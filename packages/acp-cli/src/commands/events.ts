// `acp events listen` — subscribe to the gateway's SSE event stream.
//
// We hand-parse `text/event-stream` rather than pulling in `eventsource` so
// the CLI's npm tarball stays small and so we can unit-test the parser
// against synthetic byte streams. The parser implements only the minimum
// SSE features the gateway emits (id, event, data, comment-as-heartbeat).

import { defaultConfigPath, loadConfig } from "../config.js";
import {
  type OutputContext,
  buildOutputContext,
  errorString,
  writeJson,
  writeLine,
} from "../output.js";
import type { CliConfig, GlobalFlags } from "../types.js";

export interface EventsListenFlags extends GlobalFlags {
  subjects?: string;
  /** Optional last-event-id for replay. */
  lastEventId?: string;
  /** Maximum number of events to print before exiting (test hook). */
  maxEvents?: number;
}

export interface EventsDeps {
  /** Optional fetch override. Tests pass an implementation that returns a
   * `ReadableStream` over canned bytes. */
  fetch?: typeof fetch;
  stdout?: NodeJS.WriteStream;
  stderr?: NodeJS.WriteStream;
  env?: NodeJS.ProcessEnv;
  isTTY?: boolean;
}

/** A single decoded SSE event. */
export interface SseEvent {
  id?: string;
  event?: string;
  data: string;
}

export async function runEventsListen(
  flags: EventsListenFlags,
  deps: EventsDeps = {}
): Promise<SseEvent[]> {
  const env = deps.env ?? process.env;
  const path = flags.config ?? defaultConfigPath(env);
  const cfg = await loadConfig(path);
  if (cfg === undefined) {
    throw new EventsError("no_config", `no config at ${path}; run \`acp configure\` first`);
  }
  const ctx = buildOutputContext({
    json: flags.json === true,
    ...(deps.stdout !== undefined ? { stdout: deps.stdout } : {}),
    ...(deps.stderr !== undefined ? { stderr: deps.stderr } : {}),
    env,
    ...(deps.isTTY !== undefined ? { isTTY: deps.isTTY } : {}),
  });

  const url = buildEventsUrl(cfg, flags.subjects);
  const headers: Record<string, string> = {
    Accept: "text/event-stream",
    "Cache-Control": "no-cache",
  };
  if (cfg.apiKey !== undefined) headers.Authorization = `Bearer ${cfg.apiKey}`;
  if (flags.lastEventId !== undefined) headers["Last-Event-ID"] = flags.lastEventId;

  const fetchImpl = deps.fetch ?? globalThis.fetch.bind(globalThis);
  const res = await fetchImpl(url, { headers });
  if (!res.ok) {
    throw new EventsError(
      "gateway_error",
      `events listen: gateway returned ${res.status} ${res.statusText}`
    );
  }
  if (res.body === null) {
    throw new EventsError("no_body", "events listen: gateway returned an empty response body");
  }

  const events: SseEvent[] = [];
  const max = flags.maxEvents ?? Number.POSITIVE_INFINITY;
  const decoder = new TextDecoder();
  let buffer = "";

  const reader = res.body.getReader();
  while (events.length < max) {
    const { value, done } = await reader.read();
    if (done) break;
    buffer += decoder.decode(value, { stream: true });
    while (true) {
      const nl = findEventBoundary(buffer);
      if (nl === -1) break;
      const raw = buffer.slice(0, nl);
      buffer = buffer.slice(nl).replace(/^(\r?\n)+/, "");
      const evt = parseEventBlock(raw);
      if (evt === null) continue; // pure comment/heartbeat
      events.push(evt);
      emitEvent(ctx, evt);
      if (events.length >= max) break;
    }
  }
  // Cooperatively close the stream so an injected mock body (or a real
  // server) can flush its own goroutines/promises promptly.
  try {
    reader.releaseLock();
  } catch {
    // releaseLock throws when the reader is already released; ignore.
  }
  return events;
}

/** Extract a single complete event block from `buffer`, returning the index
 * of the (consumable) trailing blank-line separator. Returns -1 when no
 * complete event is yet buffered. */
function findEventBoundary(buffer: string): number {
  // SSE separates events with a blank line. Per spec the line terminator
  // can be \n, \r, or \r\n; the gateway emits \n but be liberal in what
  // we accept.
  const matches = [
    buffer.indexOf("\n\n"),
    buffer.indexOf("\r\n\r\n"),
    buffer.indexOf("\r\r"),
  ].filter((i) => i !== -1);
  if (matches.length === 0) return -1;
  return Math.min(...matches);
}

/** Parse one event block (without its trailing blank-line separator) into
 * an `SseEvent`. Returns null when the block contained only comments and
 * therefore yields no observable event. */
export function parseEventBlock(block: string): SseEvent | null {
  const lines = block.split(/\r?\n|\r/);
  let id: string | undefined;
  let event: string | undefined;
  const data: string[] = [];
  let saw = false;
  for (const line of lines) {
    if (line.length === 0) continue;
    if (line.startsWith(":")) continue; // comment / heartbeat
    const colon = line.indexOf(":");
    const field = colon === -1 ? line : line.slice(0, colon);
    const valueRaw = colon === -1 ? "" : line.slice(colon + 1);
    // Per spec a single leading space after the colon is consumed.
    const value = valueRaw.startsWith(" ") ? valueRaw.slice(1) : valueRaw;
    switch (field) {
      case "id":
        id = value;
        saw = true;
        break;
      case "event":
        event = value;
        saw = true;
        break;
      case "data":
        data.push(value);
        saw = true;
        break;
      default:
        // unknown field — ignore per SSE spec
        break;
    }
  }
  if (!saw) return null;
  const evt: SseEvent = { data: data.join("\n") };
  if (id !== undefined) evt.id = id;
  if (event !== undefined) evt.event = event;
  return evt;
}

function emitEvent(ctx: OutputContext, evt: SseEvent): void {
  if (ctx.json) {
    // For JSON mode, attempt to parse the data payload — the gateway
    // emits JSON event envelopes — and fall back to the raw string when
    // it isn't valid JSON.
    let parsed: unknown = evt.data;
    try {
      parsed = JSON.parse(evt.data);
    } catch {
      // leave parsed as the raw string
    }
    const out: Record<string, unknown> = { data: parsed };
    if (evt.id !== undefined) out.id = evt.id;
    if (evt.event !== undefined) out.event = evt.event;
    writeJson(ctx, out);
    return;
  }
  const prefix = evt.event !== undefined ? `[${evt.event}]` : "[event]";
  const idSuffix = evt.id !== undefined ? ` (id=${evt.id})` : "";
  writeLine(ctx, `${prefix}${idSuffix} ${evt.data}`);
}

function buildEventsUrl(cfg: CliConfig, subjects?: string): string {
  const base = cfg.gatewayUrl.replace(/\/+$/, "");
  const url = new URL(`${base}/events`);
  if (subjects !== undefined && subjects.length > 0) {
    url.searchParams.set("subjects", subjects);
  }
  return url.toString();
}

export class EventsError extends Error {
  override readonly name = "EventsError";
  readonly code: string;
  constructor(code: string, message: string) {
    super(message);
    this.code = code;
  }
}

// Surface a few helpers to tests that want to assert parser behaviour
// without spinning up a fake server.
export const __testing__ = {
  findEventBoundary,
  parseEventBlock,
  buildEventsUrl,
};

// Silence the unused import warning when the file is referenced by errorString
// from tests that exercise the JSON branch.
void errorString;
