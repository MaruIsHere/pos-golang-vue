// Universal backend starter with Windows/WSL bridge.
// - Picks the correct binary name per OS (pos-backend.exe on Windows).
// - Builds with `go build ./cmd/api` if missing, then runs it from backend/ dir.
// - On WSL, can run the Linux binary natively or delegate to the Windows binary.
// Usage: npm start
import { execSync, spawn, spawnSync } from "node:child_process";
import { existsSync, readFileSync } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const rootDir = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const backendDir = path.join(rootDir, "backend");
const isWindowsHost = process.platform === "win32";
const isWSL = process.platform === "linux" && (Boolean(process.env.WSL_DISTRO_NAME) || readProcVersion().toLowerCase().includes("microsoft"));

const target = (process.env.POS_BACKEND_TARGET || (isWindowsHost ? "windows" : "auto")).toLowerCase();

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

function linuxGoAvailable() {
  return has("go") || existsSync("/usr/local/go/bin/go") || existsSync(path.join(process.env.HOME || "", ".local/go/bin/go")) || existsSync(path.join(process.env.HOME || "", "go/bin/go"));
}

function wslToWindowsPath(p) {
  const r = spawnSync("wslpath", ["-w", p], { encoding: "utf8" });
  return r.status === 0 ? r.stdout.trim() : p;
}

function windowsPathToWsl(p) {
  const r = spawnSync("wslpath", ["-u", p], { encoding: "utf8" });
  return r.status === 0 ? r.stdout.trim() : p;
}

function wslSqlitePathFromConfig() {
  try {
    const cfg = JSON.parse(readFileSync(path.join(backendDir, "config.local.json"), "utf8"));
    let p = cfg.sqlite_path || "pos.db";
    if (p === ":memory:" || p.startsWith("file:")) return p;
    if (/^[A-Za-z]:[\\/]/.test(p)) {
      if (isWSL) return windowsPathToWsl(p);
      console.warn("[start] config memakai path Windows; di Linux non-WSL memakai backend/pos.db sebagai fallback.");
      return path.join(backendDir, "pos.db");
    }
    if (path.isAbsolute(p)) return p;
    return path.resolve(backendDir, p);
  } catch {
    return path.join(backendDir, "pos.db");
  }
}

if (target === "windows" || (target === "auto" && isWSL && !linuxGoAvailable())) {
  if (!isWindowsHost) {
    const winBackend = wslToWindowsPath(backendDir);
    console.log(`[start] Running Windows binary via cmd.exe (cwd: ${winBackend})...`);
    const child = spawn("cmd.exe", ["/c", `cd /d ${winBackend} && if not exist pos-backend.exe go build -o pos-backend.exe ./cmd/api && pos-backend.exe`], { stdio: "inherit" });
    child.on("exit", (code) => process.exit(code ?? 0));
    process.on("SIGINT", () => child.kill("SIGINT"));
    process.on("SIGTERM", () => child.kill("SIGTERM"));
  } else {
    const binPath = path.join(backendDir, "pos-backend.exe");
    if (!existsSync(binPath)) {
      console.log("[start] Binary not found, building pos-backend.exe...");
      execSync("go build -o pos-backend.exe ./cmd/api", { cwd: backendDir, stdio: "inherit" });
    }
    console.log("[start] Running pos-backend.exe on windows/amd64...");
    const child = spawn(binPath, [], { cwd: backendDir, stdio: "inherit" });
    child.on("exit", (code) => process.exit(code ?? 0));
    process.on("SIGINT", () => child.kill("SIGINT"));
    process.on("SIGTERM", () => child.kill("SIGTERM"));
  }
} else {
  const binName = "pos-backend";
  const binPath = path.join(backendDir, binName);
  const go = has("go") ? "go" : existsSync("/usr/local/go/bin/go") ? "/usr/local/go/bin/go" : path.join(process.env.HOME || "", ".local/go/bin/go");

  if (!existsSync(binPath)) {
    console.log(`[start] Binary not found, building ${binName}...`);
    execSync(`${go} build -o pos-backend ./cmd/api`, { cwd: backendDir, stdio: "inherit" });
  }

  const sqlitePath = wslSqlitePathFromConfig();
  console.log(`[start] Running ${binName} on ${process.platform}/${process.arch}...`);
  console.log(`[start] POS_SQLITE_PATH=${sqlitePath}`);
  const child = spawn(binPath, [], {
    cwd: backendDir,
    stdio: "inherit",
    env: { ...process.env, POS_SQLITE_PATH: sqlitePath },
  });
  child.on("exit", (code) => process.exit(code ?? 0));
  process.on("SIGINT", () => child.kill("SIGINT"));
  process.on("SIGTERM", () => child.kill("SIGTERM"));
}
