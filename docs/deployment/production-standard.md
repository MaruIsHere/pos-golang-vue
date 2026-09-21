# Standar Deploy Produksi — Referensi Industri DEV

Status: ACUAN / REFERENSI
Tanggal: 21 Sep 2026
Tujuan: Standar deploy produksi yang baik untuk tim DEV (DevOps / SRE) — bisa jadi pelajaran industri.

---

## 1. PRINSIP: "Tidak Ada Deploy Tanpa Pipeline"

```text
Kode (Main) → CI (Build + Test) → Artifact → Staging → Approval → Production
                                          ↑
                                    Rollback (Otomatis / Manual)
```

---

## 2. ARSITEKTUR DEPLOY (SaaS POS — Multi-Tenant)

```
┌─────────────────────────────────────────────────────────────┐
│  CDN / Edge (Cloudflare) — SSL, Cache, WAF                    │
│  ├── Custom Domain Routing (per tenant / white-label)        │
└─────────────────────────────┬─────────────────────────────────┘
                              │
              ┌───────────────┼───────────────┐
              ▼               ▼               ▼
      ┌──────────────┐ ┌──────────────┐ ┌──────────────┐
      │ Load Balancer │ │ Load Balancer │ │ Load Balancer │
      │ (Nginx / Traefik) │ │ (API) │ │ (Dashboard) │
      └──────┬───────┘ └──────┬───────┘ └──────┬───────┘
             │                │                │
      ┌──────▼──────┐  ┌──────▼──────┐  ┌──────▼──────┐
      │ API Pods    │  │ PWA Pods    │  │ Dashboard   │
      │ (Fastify/Go)│  │ (Vue 3)     │  │ (Vue 3)     │
      │ Replicas: 3+│  │ Replicas: 2+│  │ Replicas: 1 │
      └──────┬──────┘  └─────────────┘  └─────────────┘
             │
      ┌──────▼──────┐
      │ PostgreSQL   │  (Primary + Read Replica + Backup Harian)
      │ Redis Cluster│  (Cache + Session + Queue)
      │ S3/MinIO     │  (Receipt Image + Backup)
      └─────────────┘
```

---

## 3. CI/CD PIPELINE (GitHub Actions — Contoh Lengkap)

```yaml
# .github/workflows/deploy.yml
name: Production Deploy

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

env:
  DOCKER_IMAGE: ghcr.io/maruisHere/pos-api:latest
  ENV: production

jobs:
  test-and-build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Setup Node
        uses: actions/setup-node@v4
        with: { node-version: '20', cache: 'pnpm' }
      - run: pnpm install --frozen-lockfile
      - run: pnpm run lint
      - run: pnpm run test -- --coverage --watchAll=false
      - run: pnpm run build
      - name: Security Scan (Snyk / Trivy)
        run: |
          npm audit --audit-level=high
          docker run --rm -v $(pwd):/repo aquasec/trivy fs /repo
      - name: Build Docker Image
        run: docker build -t $DOCKER_IMAGE .
      - name: Push to Registry
        if: github.ref == 'refs/heads/main'
        run: docker push $DOCKER_IMAGE

  deploy-staging:
    needs: test-and-build
    if: github.ref == 'refs/heads/main'
    environment: staging
    runs-on: ubuntu-latest
    steps:
      - name: Deploy to Staging
        run: |
          kubectl set image deployment/api api=$DOCKER_IMAGE -n staging
          kubectl rollout status deployment/api -n staging
      - name: Smoke Test
        run: curl -f https://staging.pos.com/health || exit 1

  deploy-production:
    needs: deploy-staging
    if: github.ref == 'refs/heads/main'
    environment:
      name: production
      url: https://api.pos.com/health
    runs-on: ubuntu-latest
    steps:
      - name: Wait for Approval
        # Manual approval gate — wajib sebelum production
        uses: trstringer/manual-approval@v1
        with:
          secret: ${{ secrets.APPROVAL_TOKEN }}
          approvers: MaruIsHere
          minimum-approvals: 1
          timeout-days: 7
      - name: Deploy Production
        run: |
          kubectl set image deployment/api api=$DOCKER_IMAGE -n production
          kubectl rollout status deployment/api -n production --timeout=300s
      - name: Health Check Post-Deploy
        run: |
          sleep 10
          curl -f https://api.pos.com/health || exit 1
          curl -f https://api.pos.com/metrics || exit 1
      - name: Rollback (Jika Gagal)
        if: failure()
        run: |
          kubectl rollout undo deployment/api -n production
          kubectl rollout status deployment/api -n production
```

---

## 4. DOCKER MULTI-STAGE (Contoh — Go + Vue)

```dockerfile
# Dockerfile — Multi-stage build
# Stage 1: Build Backend (Go)
FROM golang:1.22-alpine AS backend-builder
WORKDIR /app
COPY apps/api/ .
RUN go mod download && go build -o pos-api ./main.go

# Stage 2: Build Frontend (Vue)
FROM node:20-alpine AS frontend-builder
WORKDIR /app
COPY apps/pos-pwa/ .
COPY apps/dashboard/ ./dashboard/
RUN pnpm install --frozen-lockfile && pnpm run build:frontend && pnpm run build:dashboard

# Stage 3: Runtime (Minimal)
FROM alpine:3.19
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=backend-builder /app/pos-api ./
COPY --from=frontend-builder /app/dist ./frontend/dist
COPY --from=frontend-builder /app/dashboard/dist ./dashboard/dist
COPY apps/api/.env.example ./.env
COPY scripts/entrypoint.sh ./
RUN chmod +x entrypoint.sh
EXPOSE 3000
ENTRYPOINT ["./entrypoint.sh"]
```

---

## 5. ENVIRONMENT MANAGEMENT (12-Factor App)

```yaml
# .env.production (tidak di-commit — simpan di Vault / Secrets Manager)
NODE_ENV=production
DATABASE_URL=postgresql://user:pass@postgres:5432/pos_prod
REDIS_URL=redis://redis:6379
JWT_SECRET=<dari Vault>
MIDTRANS_SERVER_KEY=<dari Vault>
MIDTRANS_CLIENT_KEY=<dari Vault>
S3_BUCKET=pos-marketplace-prod
S3_ACCESS_KEY=<dari Vault>
FEATURE_FLAGS_DEFAULT='{"qris":true,"inventory":true,"loyalty":false}'
```

**Prinsip:**
- ✅ Konfigurasi via env vars (bukan file `.env` di repo)
- ✅ Secret dari Vault / Kubernetes Secret (bukan plain text)
- ✅ 1 image Docker → bisa deploy ke staging/prod hanya dengan env berbeda

---

## 6. DATABASE MIGRATION (Zero-Downtime)

```bash
# Pipeline — migrasi otomatis sebelum deploy
# 1. Buat snapshot (backup) — wajib
pg_dump -Fc pos_prod > backup_$(date +%Y%m%d_%H%M).sql

# 2. Migrasi (Prisma)
pnpm prisma migrate deploy --schema=prisma/schema.prisma

# 3. Cek rollback script siap
pnpm prisma migrate rollback --to-schema-datamodel=prisma/migrations/...
```

**Aturan:**
- ✅ Migration **backward-compatible** (tidak hapus kolom yang dipakai versi lama)
- ✅ Migration **tidak mengubah data** (hanya struktur) — data migration terpisah
- ❌ Tidak pernah edit `.sql` manual tanpa version control

---

## 7. HEALTH CHECK & MONITORING (Wajib)

### Health Endpoint
```typescript
// apps/api/src/health/controller.ts
GET /health
{
  status: 'ok',
  uptime: 3600,
  services: {
    database: 'connected',
    redis: 'connected',
    storage: 'connected'
  },
  version: '1.2.0',
  timestamp: new Date().toISOString()
}
```

### Monitoring Stack (Self-Hosted / Cloud)
```yaml
# docker-compose.monitoring.yml
services:
  prometheus:     # Metrik (CPU, memory, request rate)
  grafana:        # Dashboard visual
  loki:           # Log aggregation
  alertmanager:   # Alert (Slack/Email) jika error rate > 5% atau latency > 2s
```

**Alert Wajib:**
- ✅ Error rate > 5% selama 5 menit → Alert
- ✅ Response time (p95) > 2 detik → Warning
- ✅ Disk usage > 80% → Warning
- ✅ Database connection pool > 80% → Warning
- ✅ Pod restart > 3x dalam 5 menit → Critical

---

## 8. ROLLBACK STRATEGY (Wajib Ada Sebelum Deploy)

```bash
# 1. Blue-Green Deploy (Ideal)
# Versi baru (Green) deploy → health check → traffic switch → Blue tetap standby 10 menit → hapus Blue

# 2. Manual Rollback (Jika Blue-Green tidak mungkin)
kubectl rollout undo deployment/api -n production --to-revision=prev
# Atau via Docker:
docker-compose -f production.yml up --force-recreate -d
# Atau restore DB:
pg_restore -Fc backup.sql > pos_prod_rollback
```

**Aturan:** Rollback harus bisa dieksekusi dalam **< 3 menit** tanpa akses ke server manual (via script / CI). Setiap deploy harus menyimpan `revision` atau `docker image tag`.

---

## 9. SECURITY (Industry Standard)

```typescript
// Middleware wajib (Fastify / Express)
app.register(helmet)        // Security headers
app.register(cors, { origin: ['https://*.pos.com'] })
app.register(rateLimit, { max: 100, timeWindow: '1 minute' })
```

**Headers Wajib:**
- `Content-Security-Policy`
- `X-Content-Type-Options: nosniff`
- `X-Frame-Options: DENY`
- `Strict-Transport-Security`
- `Referrer-Policy`

**API Security:**
- ✅ Rate limiting per IP + per API key
- ✅ Input validation (`zod` / `joi`) — tidak pernah percaya input user
- ✅ Parameterized query (Prisma) — tidak pernah raw SQL tanpa parameter
- ✅ JWT refresh rotation (deteksi reuse → revoke semua token user)
- ✅ Audit log semua operasi sensitif (create/delete/update tenant, billing)

---

## 10. ZERO-DOWNTIME DEPLOY (Contoh Lengkap)

```bash
#!/bin/bash
# scripts/deploy.sh
set -euo pipefail

VERSION=$(git rev-parse --short HEAD)
IMAGE="ghcr.io/maruisHere/pos-api:${VERSION}"

# 1. Build
make build IMAGE=${IMAGE}

# 2. Database migration (sebelum deploy image baru)
pnpm prisma migrate deploy --schema=prisma/schema.prisma

# 3. Push image
make push IMAGE=${IMAGE}

# 4. Deploy (rolling update)
kubectl set image deployment/api api=${IMAGE} -n production
kubectl rollout status deployment/api -n production --timeout=300s

# 5. Health check post-deploy
sleep 5
curl -sf https://api.pos.com/health || { echo "DEPLOY GAGAL"; kubectl rollout undo deployment/api -n production; exit 1; }

# 6. Notifikasi
curl -X POST https://hooks.slack.com/services/... -d '{"text":"✅ POS Deployed ${VERSION}"}'
```

---

## 11. CHECKLIST DEPLOY PRODUKSI (Wajib Dibaca Sebelum Push)

```markdown
[ ] 1. .env tidak berisi secret (gunakan Vault / Kubernetes Secret)
[ ] 2. .gitignore melindungi .env, node_modules, .db, *.log
[ ] 3. Image Docker multi-stage (tidak menyimpan build artifact)
[ ] 4. Health check (`/health`) berjalan di semua service
[ ] 5. Database migration backward-compatible dan sudah diuji
[ ] 6. Rollback script sudah diuji (tidak hanya teori)
[ ] 7. Monitoring (Prometheus / Grafana) aktif dan alert terkonfigurasi
[ ] 8. Rate limit dan security headers aktif
[ ] 9. Audit log aktif untuk operasi sensitif
[ ] 10. Backup database otomatis (harian) dan diuji restore
[ ] 11. SSL aktif (Let's Encrypt / Cloudflare) — tidak HTTP
[ ] 12. Log rotasi (tidak memenuhi disk)
```

---

## 12. REFERENSI UNTUK INDUSTRI

- **12-Factor App**: https://12factor.net/ (Standar aplikasi modern — wajib baca)
- **Docker Multi-Stage**: https://docs.docker.com/build/building/multi-stage/
- **Kubernetes Zero-Downtime**: https://kubernetes.io/docs/concepts/workloads/controllers/deployment/
- **PostgreSQL Migration Best Practice**: https://www.postgresql.org/docs/current/ (Backup/Restore)
- **Security Headers (OWASP)**: https://owasp.org/www-project-secure-headers/

---

## 13. CATATAN UNTUK TIM (`MaruIsHere` / Kontributor)

File ini (`docs/deployment/production-standard.md`) adalah **referensi industri** — bukan sekadar catatan internal. Jika kamu (`MaruIsHere`) atau kontributor menambahkan fitur:

1. **Baca checklist (Bagian 11)** sebelum deploy
2. **Jika fitur besar** → buat ADR (`docs/decisions/`) dulu
3. **Jika fitur baru** → cek apakah perlu update dokumen ini (misal: plugin baru → update `COMPONENTS.md`)
4. **Tidak pernah** commit `.env`, `node_modules`, `.db`, `*.exe`, `dist/` — `.gitignore` sudah melindungi (cek commit `9990ee05`)

---

*Disusun oleh Hermes Agent — 21 Sep 2026*
*Untuk tim POS (`pos-golang-vue`) sebagai referensi deploy produksi. Dokumen `pulau-segara` (game) hidup terpisah di `~/Projects/gamedev/pulau-segara/` — jangan dicampur dengan standar POS ini.*
