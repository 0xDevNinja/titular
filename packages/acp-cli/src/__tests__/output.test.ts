import { describe, expect, it } from "vitest";

import {
  buildOutputContext,
  errorString,
  stringifyBigInts,
  writeJson,
  writeKv,
  writeLine,
  writeTable,
} from "../output.js";

import { CollectingStream } from "./_helpers.js";

describe("stringifyBigInts", () => {
  it("converts top-level bigint", () => {
    expect(stringifyBigInts(42n)).toBe("42");
  });

  it("recurses into objects", () => {
    expect(stringifyBigInts({ a: 1n, b: { c: 2n } })).toEqual({
      a: "1",
      b: { c: "2" },
    });
  });

  it("recurses into arrays", () => {
    expect(stringifyBigInts([1n, 2n])).toEqual(["1", "2"]);
  });

  it("leaves non-bigints unchanged", () => {
    expect(stringifyBigInts("x")).toBe("x");
    expect(stringifyBigInts(7)).toBe(7);
    expect(stringifyBigInts(null)).toBeNull();
  });
});

describe("buildOutputContext", () => {
  it("disables colour when NO_COLOR is set", () => {
    const ctx = buildOutputContext({
      json: false,
      env: { NO_COLOR: "1" },
      isTTY: true,
    });
    expect(ctx.color).toBe(false);
  });

  it("disables colour for non-TTY", () => {
    const ctx = buildOutputContext({ json: false, env: {}, isTTY: false });
    expect(ctx.color).toBe(false);
  });

  it("respects FORCE_COLOR", () => {
    const ctx = buildOutputContext({
      json: false,
      env: { FORCE_COLOR: "1" },
      isTTY: false,
    });
    expect(ctx.color).toBe(true);
  });

  it("forces colour off in JSON mode regardless of TTY", () => {
    const ctx = buildOutputContext({ json: true, env: {}, isTTY: true });
    expect(ctx.color).toBe(false);
  });
});

describe("writeJson / writeLine / writeKv", () => {
  function ctxWith(stream: CollectingStream, json = false) {
    return buildOutputContext({
      json,
      stdout: stream,
      stderr: stream,
      env: { NO_COLOR: "1" },
      isTTY: false,
    });
  }

  it("writeJson stringifies bigints", () => {
    const s = new CollectingStream();
    writeJson(ctxWith(s, true), { id: 7n });
    expect(s.text()).toContain('"id": "7"');
  });

  it("writeLine appends a newline", () => {
    const s = new CollectingStream();
    writeLine(ctxWith(s), "hello");
    expect(s.text()).toBe("hello\n");
  });

  it("writeKv prints key: value", () => {
    const s = new CollectingStream();
    writeKv(ctxWith(s), "key", "value");
    expect(s.text()).toContain("key:");
    expect(s.text()).toContain("value");
  });
});

describe("writeTable", () => {
  it("renders an empty marker for no rows", () => {
    const s = new CollectingStream();
    const ctx = buildOutputContext({
      json: false,
      stdout: s,
      stderr: s,
      env: { NO_COLOR: "1" },
      isTTY: false,
    });
    writeTable(ctx, ["a", "b"], []);
    expect(s.text()).toContain("(no rows)");
  });

  it("renders header + rows", () => {
    const s = new CollectingStream();
    const ctx = buildOutputContext({
      json: false,
      stdout: s,
      stderr: s,
      env: { NO_COLOR: "1" },
      isTTY: false,
    });
    writeTable(
      ctx,
      ["id", "name"],
      [
        ["1", "alice"],
        ["2", "bob"],
      ]
    );
    const text = s.text();
    expect(text).toContain("id");
    expect(text).toContain("alice");
    expect(text).toContain("bob");
  });
});

describe("errorString", () => {
  it("returns the raw string when colour is off", () => {
    const ctx = buildOutputContext({
      json: false,
      env: { NO_COLOR: "1" },
      isTTY: false,
    });
    expect(errorString(ctx, "boom")).toBe("boom");
  });
});
