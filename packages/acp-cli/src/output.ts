// Output helpers shared across commands. Two motivations:
//   1. Always serialise `bigint` and `Uint8Array`-shaped fields as decimal
//      strings before writing JSON; `JSON.stringify` throws on bigints
//      otherwise, and the gateway already encodes uint256 as decimal strings
//      so the round-trip is lossless.
//   2. Centralise NO_COLOR / TTY detection so every command honours the env
//      identically (`NO_COLOR=1`, non-tty stdout, or `--no-color` all yield
//      uncoloured output).

import pc from "picocolors";

/** Stream sink, narrowed for test injection. `process.stdout` and
 * `process.stderr` both satisfy this. */
export interface WriteStream {
  write(chunk: string): boolean | undefined;
}

/** Bag of streams + flags every command receives. Keeping all I/O behind
 * this object makes the CLI fully testable without spawning a child
 * process. */
export interface OutputContext {
  stdout: WriteStream;
  stderr: WriteStream;
  json: boolean;
  color: boolean;
}

/** Build an `OutputContext` from the current process. The `json` flag is the
 * caller's responsibility (it usually comes from `--json`). */
export function buildOutputContext(opts: {
  json: boolean;
  stdout?: WriteStream;
  stderr?: WriteStream;
  env?: NodeJS.ProcessEnv;
  isTTY?: boolean;
}): OutputContext {
  const env = opts.env ?? process.env;
  const stdout = opts.stdout ?? process.stdout;
  const stderr = opts.stderr ?? process.stderr;

  // Honour the `NO_COLOR` standard (https://no-color.org), the FORCE_COLOR
  // override, and a non-TTY stdout. Tests pass `isTTY: false` explicitly so
  // assertions are not contaminated by ANSI escapes from a developer's
  // terminal.
  const isTTY = opts.isTTY ?? (process.stdout as { isTTY?: boolean }).isTTY === true;
  const forceColor = env.FORCE_COLOR;
  const noColor = env.NO_COLOR;
  let color = isTTY;
  if (noColor !== undefined && noColor.length > 0) color = false;
  if (forceColor !== undefined && forceColor.length > 0 && forceColor !== "0") {
    color = true;
  }
  if (opts.json) color = false;

  return { stdout, stderr, json: opts.json, color };
}

/** Replace `bigint` values with their decimal string form recursively so the
 * payload is `JSON.stringify`-safe. Cycles are not expected in CLI output;
 * if one is ever introduced, JSON.stringify will surface it. */
export function stringifyBigInts(value: unknown): unknown {
  if (typeof value === "bigint") return value.toString(10);
  if (Array.isArray(value)) return value.map((v) => stringifyBigInts(v));
  if (value !== null && typeof value === "object") {
    const out: Record<string, unknown> = {};
    for (const [k, v] of Object.entries(value as Record<string, unknown>)) {
      out[k] = stringifyBigInts(v);
    }
    return out;
  }
  return value;
}

/** Write `value` to stdout as a JSON document followed by a single newline. */
export function writeJson(ctx: OutputContext, value: unknown): void {
  ctx.stdout.write(`${JSON.stringify(stringifyBigInts(value), null, 2)}\n`);
}

/** Write a single text line to stdout. */
export function writeLine(ctx: OutputContext, line: string): void {
  ctx.stdout.write(`${line}\n`);
}

/** Write a key/value pair, dimming the key when colour is enabled. */
export function writeKv(ctx: OutputContext, key: string, value: string): void {
  const k = ctx.color ? pc.dim(`${key}:`) : `${key}:`;
  ctx.stdout.write(`${k.padEnd(ctx.color ? 28 : 16)} ${value}\n`);
}

/** Render a list-of-records as a fixed-width table. Falls back to plain
 * spacing rather than ASCII borders so output stays readable when piped. */
export function writeTable(
  ctx: OutputContext,
  columns: readonly string[],
  rows: readonly (readonly string[])[]
): void {
  if (rows.length === 0) {
    writeLine(ctx, ctx.color ? pc.dim("(no rows)") : "(no rows)");
    return;
  }
  const widths = columns.map((c, i) => Math.max(c.length, ...rows.map((r) => (r[i] ?? "").length)));
  const header = columns.map((c, i) => c.padEnd(widths[i] ?? c.length)).join("  ");
  const separator = widths.map((w) => "-".repeat(w)).join("  ");
  const headerLine = ctx.color ? pc.bold(header) : header;
  const sepLine = ctx.color ? pc.dim(separator) : separator;
  ctx.stdout.write(`${headerLine}\n${sepLine}\n`);
  for (const row of rows) {
    const line = row.map((cell, i) => (cell ?? "").padEnd(widths[i] ?? 0)).join("  ");
    ctx.stdout.write(`${line}\n`);
  }
}

/** Stylise a message as an error, dimming when colour is disabled. */
export function errorString(ctx: OutputContext, msg: string): string {
  return ctx.color ? pc.red(msg) : msg;
}
