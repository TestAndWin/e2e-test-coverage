// Checks that every package in package-lock.json was published at least
// DAYS days ago (supply chain cooldown). Exits with 1 if a younger one is found.
//
// Run it from the npm project directory:
//
//	node ../scripts/npmcooldown.mjs [days]   (default 7)
import { readFileSync } from 'node:fs'

const days = Number(process.argv[2] ?? 7)
const cutoff = Date.now() - days * 24 * 60 * 60 * 1000
const registry = 'https://registry.npmjs.org/'

const lock = JSON.parse(readFileSync('package-lock.json', 'utf8'))
const wanted = new Map() // name -> Set(version)
for (const [path, pkg] of Object.entries(lock.packages ?? {})) {
  if (!path || pkg.link || !pkg.version || !pkg.resolved?.startsWith(registry)) continue
  const name = pkg.name ?? path.slice(path.lastIndexOf('node_modules/') + 'node_modules/'.length)
  if (!wanted.has(name)) wanted.set(name, new Set())
  wanted.get(name).add(pkg.version)
}

const young = []
const failed = []
const names = [...wanted.keys()]
async function worker() {
  while (names.length) {
    const name = names.pop()
    try {
      const res = await fetch(registry + name.replace('/', '%2f'))
      if (!res.ok) throw new Error(`HTTP ${res.status}`)
      const { time = {} } = await res.json()
      for (const v of wanted.get(name)) {
        const t = Date.parse(time[v])
        if (Number.isNaN(t)) failed.push(`${name}@${v}: no publish time`)
        else if (t > cutoff) young.push(`${name}@${v} published ${time[v].slice(0, 10)}`)
      }
    } catch (e) {
      failed.push(`${name}: ${e.message}`)
    }
  }
}
await Promise.all(Array.from({ length: 16 }, worker))

for (const f of failed.sort()) console.log(`ERROR ${f}`)
for (const y of young.sort()) console.log(`YOUNG ${y} (< ${days} days)`)
if (!young.length && !failed.length) {
  console.log(`all ${wanted.size} locked npm packages are at least ${days} days old`)
}
process.exit(young.length || failed.length ? 1 : 0)
