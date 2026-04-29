#!/usr/bin/env node
// One-shot rebuilder for `src/lib/mock-bytecode.ts`. Pipes
// `src/lib/mocks/*.sol` through the local `solc` binary and patches the
// generated runtime hex into the existing constants in mock-bytecode.ts.
//
// The script is opt-in (only runs on `pnpm -F @titular/acp-sdk-e2e run
// compile:mocks`), not on install or build, because:
//
//   * Most contributors will never touch the mocks.
//   * The mock runtime is small enough that pinning the hex keeps the
//     e2e package free of a `solc` dependency at install time.
//   * CI doesn't need solc; it imports the constants directly.
//
// The runtime hex is emitted by `solc --optimize --bin-runtime`, which
// is the form that `anvil_setCode` accepts. We deliberately do NOT use
// the deploy ("creation") bytecode — installing runtime bytes via the
// anvil cheatcode bypasses constructor execution, which is what we want
// for stateless event-emitters.

import { execSync } from "node:child_process";
import { readFileSync, writeFileSync } from "node:fs";
import { dirname, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const here = dirname(fileURLToPath(import.meta.url));
const pkgRoot = resolve(here, "..");

const mocks = [
  {
    sol: "src/lib/mocks/MockAgentRegistry.sol",
    contractName: "MockAgentRegistry",
    constName: "MOCK_AGENT_REGISTRY_RUNTIME",
  },
  {
    sol: "src/lib/mocks/MockJobFactory.sol",
    contractName: "MockJobFactory",
    constName: "MOCK_JOB_FACTORY_RUNTIME",
  },
];

// Single solc invocation gives us all artifacts at once and matches the
// exact flags the original constants were derived with.
const solcOut = execSync(`solc --optimize --bin-runtime ${mocks.map((m) => m.sol).join(" ")}`, {
  cwd: pkgRoot,
  encoding: "utf8",
});

// Parse `solc --bin-runtime` output. It groups one entry per contract:
//
//   ======= path/to/File.sol:ContractName =======
//   Binary of the runtime part:
//   <hex>
//
// We pick the entry whose ContractName matches each mock and patch the
// resulting `0x<hex>` back into the source TS constant in-place.
const sectionRe =
  /=======\s+([^\s]+):([A-Za-z_][\w]*)\s+=======[\s\S]*?Binary of the runtime part:\s*\n([0-9a-fA-F]+)/g;
const found = new Map();
for (const m of solcOut.matchAll(sectionRe)) {
  found.set(m[2], m[3]);
}

const tsPath = resolve(pkgRoot, "src/lib/mock-bytecode.ts");
let ts = readFileSync(tsPath, "utf8");

for (const { contractName, constName } of mocks) {
  const hex = found.get(contractName);
  if (!hex) {
    throw new Error(
      `compile-mocks: solc output did not include ${contractName}; got contracts: ${[...found.keys()].join(", ") || "<none>"}`
    );
  }
  const re = new RegExp(`(export const ${constName} =\\s*\\n?\\s*)"0x[0-9a-fA-F]*"`);
  if (!re.test(ts)) {
    throw new Error(`compile-mocks: could not find constant ${constName} in mock-bytecode.ts`);
  }
  ts = ts.replace(re, `$1"0x${hex}"`);
}

writeFileSync(tsPath, ts);
console.log(`compile-mocks: rewrote ${tsPath} with ${mocks.length} runtime artefact(s)`);
