import { describe, expect, it } from "vitest";
import { decodeFirstEvent, eventTopic } from "../log-decoder.js";
import type { ParsedLog } from "../providers/types.js";
import type { EvmAddress } from "../types.js";
import { encodeEventToLog } from "./_helpers.js";

const ADDR: EvmAddress = "0x1111111111111111111111111111111111111111";

const FOO_EVENT = {
  type: "event",
  name: "Foo",
  inputs: [{ name: "x", type: "uint256", indexed: false }],
  anonymous: false,
} as const;

const BAR_EVENT = {
  type: "event",
  name: "Bar",
  inputs: [{ name: "y", type: "address", indexed: true }],
  anonymous: false,
} as const;

const ABI = [FOO_EVENT, BAR_EVENT, { type: "function", name: "noop", inputs: [], outputs: [] }];

function fooLog(x: bigint): ParsedLog {
  return encodeEventToLog(FOO_EVENT, { x }, ADDR);
}

describe("eventTopic", () => {
  it("computes a stable 32-byte hash", () => {
    const t = eventTopic(FOO_EVENT);
    expect(t).toMatch(/^0x[0-9a-f]{64}$/);
  });
});

describe("decodeFirstEvent", () => {
  it("decodes the first matching log", () => {
    const args = decodeFirstEvent<{ x: bigint }>(ABI, "Foo", [fooLog(42n)]);
    expect(args?.x).toBe(42n);
  });

  it("returns null when no log matches", () => {
    const args = decodeFirstEvent(ABI, "Foo", []);
    expect(args).toBeNull();
  });

  it("returns null when the event name is not in the ABI", () => {
    const args = decodeFirstEvent(ABI, "DoesNotExist", [fooLog(1n)]);
    expect(args).toBeNull();
  });

  it("skips logs with mismatched topic[0]", () => {
    const stranger: ParsedLog = {
      address: ADDR,
      topics: ["0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef"],
      data: "0x",
      logIndex: 0,
    };
    const args = decodeFirstEvent<{ x: bigint }>(ABI, "Foo", [stranger, fooLog(7n)]);
    expect(args?.x).toBe(7n);
  });

  it("skips logs with empty topics", () => {
    const empty: ParsedLog = {
      address: ADDR,
      topics: [],
      data: "0x",
      logIndex: 0,
    };
    const args = decodeFirstEvent<{ x: bigint }>(ABI, "Foo", [empty, fooLog(9n)]);
    expect(args?.x).toBe(9n);
  });

  it("skips logs whose data fails to decode and continues scanning", () => {
    const corrupt: ParsedLog = {
      address: ADDR,
      topics: [eventTopic(FOO_EVENT)],
      data: "0x00",
      logIndex: 0,
    };
    const args = decodeFirstEvent<{ x: bigint }>(ABI, "Foo", [corrupt, fooLog(11n)]);
    expect(args?.x).toBe(11n);
  });
});
