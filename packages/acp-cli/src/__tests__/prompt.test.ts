import { Writable } from "node:stream";
import { describe, expect, it } from "vitest";

import { ask, closePromptSession, select } from "../prompt.js";

import { readableFromLines } from "./_helpers.js";

class SinkWritable extends Writable {
  override _write(_chunk: unknown, _enc: string, cb: (e?: Error | null) => void): void {
    cb();
  }
}

function ioFromLines(lines: readonly string[]) {
  return { input: readableFromLines(lines), output: new SinkWritable() };
}

describe("ask", () => {
  it("returns the trimmed reply", async () => {
    const io = ioFromLines(["  hello  "]);
    expect(await ask("q", undefined, io)).toBe("hello");
    closePromptSession(io);
  });

  it("returns the fallback when the user submits an empty line", async () => {
    const io = ioFromLines([""]);
    expect(await ask("q", "default", io)).toBe("default");
    closePromptSession(io);
  });

  it("throws when the stream closes without a reply and no fallback", async () => {
    const io = ioFromLines([]);
    await expect(ask("q", undefined, io)).rejects.toThrow(/prompt closed/);
    closePromptSession(io);
  });
});

describe("select", () => {
  it("accepts a numeric reply", async () => {
    const io = ioFromLines(["2"]);
    expect(await select("pick", ["red", "green", "blue"], undefined, io)).toBe("green");
    closePromptSession(io);
  });

  it("accepts a slug reply", async () => {
    const io = ioFromLines(["blue"]);
    expect(await select("pick", ["red", "green", "blue"], undefined, io)).toBe("blue");
    closePromptSession(io);
  });

  it("falls back to the default when blank", async () => {
    const io = ioFromLines([""]);
    expect(await select("pick", ["a", "b"], "a", io)).toBe("a");
    closePromptSession(io);
  });

  it("rejects an unknown reply", async () => {
    const io = ioFromLines(["purple"]);
    await expect(select("pick", ["red", "blue"], undefined, io)).rejects.toThrow(/not one of/);
    closePromptSession(io);
  });

  it("rejects an empty choice list", async () => {
    const io = ioFromLines([]);
    await expect(select("pick", [], undefined, io)).rejects.toThrow(/must not be empty/);
    closePromptSession(io);
  });
});
