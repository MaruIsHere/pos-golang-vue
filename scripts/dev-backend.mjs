// Bridge backend runner for Windows / WSL / Linux.
//
// Why:
// - Windows dev expects paths like C:\... and a Windows Go toolchain.
// - WSL dev expects paths like /mnt/c/... and a Linux Go toolchain.
// Both can still use the same physical SQLite file.
//
// Usage:
//   node scripts/dev-backend.mjs              # auto choose
//   node scripts/dev-backend.mjs --target=windows
//   node scripts/dev-backend.mjs --target=wsl
//   POS_BACKEND_TARGET=windows node scripts/dev-backend.mjs
import { existsSync, readFileSync } from "node:fs";
import { spawn, spawnSync } from "node:child_process";
import path from "node:path";
import { fileURLToPath } from "node:url";

const rootDir = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const backendDir = path.join(rootDir, "backend");
const isWindowsHost = process.platform === "win32";
const isLinuxHost = process.platform === "linux";
const isWSL = isLinuxHost && (Boolean(process.env.WSL_DISTRO_NAME) || readProcVersion().toLowerCase().includes("microsoft"));

const argTarget = process.argv.find((a) => a.startsWith("--target="))?.split("=")[1];
const requestedTarget = (argTarget || process.env.POS_BACKEND_TARGET || "auto").toLowerCase();

function readProcVersion() {
  try {
    return readFileSync("/proc/version", "utf8");
  } catch {
    return "";
  }
}

function has(cmd, args = ["--version"]) {
  const r = spawnSync(cmd, args, { stdio: "ignore" });
  return r.status === 0;
}

function wslToWindowsPath(p) {
  const r = spawnSync("wslpath", ["-w", p], { encoding: "utf8" });
  return r.status === 0 ? r.stdout.trim() : p;
}

function windowsPathToWsl(p) {
  const r = spawnSync("wslpath", ["-u", p], { encoding: "utf8" });
  return r.status === 0 ? r.stdout.trim() : p;
}

function linuxGoBin() {
  if (process.env.GO && existsSync(process.env.GO)) return process.env.GO;
  if (has("go", ["version"])) return "go";
  const candidates = ["/usr/bin/go", "/usr/local/go/bin/go", path.join(process.env.HOME || "", ".local/go/bin/go"), path.join(process.env.HOME || "", "go/bin/go")];
  for (const p of candidates) {
    if (p && existsSync(p)) return p;
  }
  return null;
}

function resolveTarget() {
  if (requestedTarget === "windows" || requestedTarget === "win") return "windows";
  if (requestedTarget === "wsl" || requestedTarget === "linux") return "wsl";

  if (isWindowsHost) return "windows";
  if (isWSL) return linuxGoBin() ? "wsl" : "windows";
  return "wsl";
}

function localConfig() {
  try {
    return JSON.parse(readFileSync(path.join(backendDir, "config.local.json"), "utf8"));
  } catch {
    try {
      return JSON.parse(readFileSync(path.join(backendDir, "config.json"), "utf8"));
    } catch {
      return {};
    }
  }
}

function wslSqlitePathFromConfig() {
  const cfg = localConfig();
  let p = cfg.sqlite_path || "pos.db";
  if (p === ":memory:" || p.startsWith("file:")) return p;
  if (/^[A-Za-z]:[\\/]/.test(p)) {
    if (isWSL) return windowsPathToWsl(p);
    console.warn("[bridge] config memakai path Windows; di Linux non-WSL memakai backend/pos.db sebagai fallback.");
    return path.join(backendDir, "pos.db");
  }
  if (path.isAbsolute(p)) return p;
  return path.resolve(backendDir, p);
}

function runWindowsFromWindows() {
  const go = process.env.GO && existsSync(process.env.GO) ? process.env.GO : "go.exe";
  console.log("[bridge] Running backend with Windows Go toolchain...");
  const child = spawn(go, ["run", "./cmd/api"], { cwd: backendDir, stdio: "inherit" });
  wireChild(child);
}

function runWindowsFromWsl() {
  const winBackend = wslToWindowsPath(backendDir);
  console.log(`[bridge] Running backend on Windows side via cmd.exe (cwd: ${winBackend})...`);
  const child = spawn("cmd.exe", ["/c", `cd /d ${winBackend} && go run ./cmd/api`], { stdio: "inherit" });
  wireChild(child);
}

function runNativeLinux() {
  const go = linuxGoBin();
  if (!go) {
    console.error("[bridge] Linux Go toolchain not found. Install Go in WSL or run with --target=windows.");
    process.exit(1);
  }

  const sqlitePath = wslSqlitePathFromConfig();
  console.log("[bridge] Running backend with Linux/WSL Go toolchain...");
  console.log(`[bridge] POS_SQLITE_PATH=${sqlitePath}`);
  const child = spawn(go, ["run", "./cmd/api"], {
    cwd: backendDir,
    stdio: "inherit",
    env: {
      ...process.env,
      POS_SQLITE_PATH: sqlitePath,
    },
  });
  wireChild(child);
}

function wireChild(child) {
  child.on("exit", (code) => process.exit(code ?? 0));
  process.on("SIGINT", () => child.kill("SIGINT"));
  process.on("SIGTERM", () => child.kill("SIGTERM"));
}

const target = resolveTarget();
if (isWSL && target === "windows") {
  console.warn("[bridge] WARNING: backend akan jalan di Windows. Jika frontend juga jalan di WSL, proxy /api ke 127.0.0.1 biasanya tidak nyambung.");
}
console.log(`[bridge] platform=${process.platform}, wsl=${isWSL}, target=${target}`);

if (target === "windows") {
  if (isWindowsHost) runWindowsFromWindows();
  else if (isWSL) runWindowsFromWsl();
  else {
    console.error("[bridge] target=windows is only available on Windows or WSL.");
    process.exit(1);
  }
} else {
  runNativeLinux();
}
