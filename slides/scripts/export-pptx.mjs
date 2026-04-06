import { existsSync, mkdirSync, renameSync, rmSync } from "node:fs"
import { resolve } from "node:path"
import { spawnSync } from "node:child_process"

const cwd = process.cwd()
const bin = process.platform === "win32"
  ? resolve(cwd, "node_modules", ".bin", "slidev.cmd")
  : resolve(cwd, "node_modules", ".bin", "slidev")

const tempOutput = resolve(cwd, "slides-export.pptx")
const finalDir = resolve(cwd, "exports")
const finalOutput = resolve(finalDir, "storm-defense.pptx")

rmSync(tempOutput, { force: true })
rmSync(finalOutput, { force: true })
mkdirSync(finalDir, { recursive: true })

const result = spawnSync(bin, ["export", "slides.md", "--format", "pptx"], {
  cwd,
  stdio: "inherit",
  shell: process.platform === "win32",
})

if (result.error) {
  console.error(result.error)
  process.exit(1)
}

if (result.status !== 0) {
  process.exit(result.status ?? 1)
}

if (!existsSync(tempOutput)) {
  console.error(`PPTX export missing: ${tempOutput}`)
  process.exit(1)
}

renameSync(tempOutput, finalOutput)
console.log(`PPTX exporte vers ${finalOutput}`)
