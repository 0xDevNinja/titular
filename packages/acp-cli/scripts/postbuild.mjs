// Make the compiled cli entry executable so `npx acp` and bare invocations
// work after `pnpm build`. tsc doesn't preserve the exec bit and pnpm doesn't
// chmod bin files until publish-time, which means the locally-linked binary
// fails before this script runs.

import { chmod } from "node:fs/promises";
import { resolve } from "node:path";

const target = resolve(process.cwd(), "dist/cli.js");
await chmod(target, 0o755);
