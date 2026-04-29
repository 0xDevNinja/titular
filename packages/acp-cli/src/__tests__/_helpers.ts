// Shared test helpers — avoids duplicating the WriteStream + tempdir
// boilerplate across every command's spec.

import { mkdtemp } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import { Readable } from "node:stream";

/** A WriteStream stand-in that captures every chunk into a string array. */
export class CollectingStream {
  readonly chunks: string[] = [];
  write(chunk: string): boolean {
    this.chunks.push(chunk);
    return true;
  }
  text(): string {
    return this.chunks.join("");
  }
}

/** Create a fresh tempdir + return the path of a config.json that doesn't
 * exist yet. The path is returned even though the file is absent so callers
 * can pass it straight to `loadConfig` / `saveConfig`. */
export async function freshConfigPath(): Promise<{
  path: string;
  dir: string;
}> {
  const dir = await mkdtemp(join(tmpdir(), "acp-cli-test-"));
  return { path: join(dir, "config.json"), dir };
}

/** Build a fake `Response`-like object whose `body` is a `ReadableStream`
 * that yields the given chunks in order. Used by the events listener tests
 * to drive the SSE parser without a network. */
export function fakeFetchSse(chunks: readonly string[]): typeof fetch {
  return ((..._args: unknown[]) => {
    const encoder = new TextEncoder();
    const body = new ReadableStream<Uint8Array>({
      start(controller) {
        for (const chunk of chunks) controller.enqueue(encoder.encode(chunk));
        controller.close();
      },
    });
    return Promise.resolve(
      new Response(body, {
        status: 200,
        headers: { "Content-Type": "text/event-stream" },
      })
    );
  }) as typeof fetch;
}

/** Build a fake `fetch` that returns canned JSON for any URL. */
export function fakeFetchJson(value: unknown): typeof fetch {
  return ((..._args: unknown[]) => {
    return Promise.resolve(
      new Response(JSON.stringify(value), {
        status: 200,
        headers: { "Content-Type": "application/json" },
      })
    );
  }) as typeof fetch;
}

/** Build a synthetic readable stream from a list of lines (each terminated
 * with `\n`). Used by the prompt tests so we don't have to wire stdin. */
export function readableFromLines(lines: readonly string[]): Readable {
  return Readable.from([lines.map((l) => `${l}\n`).join("")]);
}
