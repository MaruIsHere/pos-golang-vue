// Universal backend starter: works on Windows, macOS, and Linux.
// - Picks the correct binary name per OS (pos-backend.exe on Windows).
// - Builds it with `go build` if missing, then runs it from backend/ dir
//   (backend expects config.json + ../frontend/dist relative to cwd).
// Usage: npm start
import { execSync, spawn } from "node:child_process";
import { existsSync } from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const rootDir = path.dirname(path.dirname(fileURLToPath(import.meta.url)));
const backendDir = path.join(rootDir, "backend");
const binName = process.platform === "win32" ? "pos-backend.exe" : "pos-backend";
const binPath = path.join(backendDir, binName);

if (!existsSync(binPath)) {
  console.log(`[start] Binary not found, building ${binName}...`);
  execSync("go build -o pos-backend .", { cwd: backendDir, stdio: "inherit" });
}

console.log(`[start] Running ${binName} on ${process.platform}/${process.arch}...`);
const child = spawn(binPath, [], { cwd: backendDir, stdio: "inherit" });
child.on("exit", (code) => process.exit(code ?? 0));
process.on("SIGINT", () => child.kill("SIGINT"));
process.on("SIGTERM", () => child.kill("SIGTERM"));
