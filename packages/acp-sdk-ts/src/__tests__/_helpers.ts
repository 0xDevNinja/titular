// Test helpers shared across agent / log-decoder unit tests. Builds an
// EVM-format event log (topics + data) from an event ABI fragment + named
// args, using the `encodeEventTopics` + `encodeAbiParameters` exports that
// viem ships in its public surface (it does not export `encodeEventLog`).

import { type AbiEvent, type AbiParameter, encodeAbiParameters, encodeEventTopics } from "viem";
import type { ParsedLog } from "../providers/types.js";
import type { EvmAddress } from "../types.js";

/**
 * Encode an event log to a {@link ParsedLog}, splitting `args` into the
 * indexed (topic) and non-indexed (data) buckets the EVM emits.
 *
 * `event` is intentionally typed loosely so callers can pass narrow `as const`
 * fragments without TypeScript widening complaints.
 */
export function encodeEventToLog(
  event: AbiEvent,
  args: Record<string, unknown>,
  address: EvmAddress,
  logIndex = 0
): ParsedLog {
  const indexedInputs = event.inputs.filter(
    (i): i is AbiEvent["inputs"][number] & { name: string } =>
      i.indexed === true && typeof i.name === "string"
  );
  const dataInputs = event.inputs.filter(
    (i): i is AbiEvent["inputs"][number] & { name: string } =>
      i.indexed !== true && typeof i.name === "string"
  );

  const indexedArgs: Record<string, unknown> = {};
  for (const i of indexedInputs) indexedArgs[i.name] = args[i.name];

  const rawTopics = encodeEventTopics({
    abi: [event],
    eventName: event.name,
    args: indexedArgs,
  });
  // encodeEventTopics returns ( `0x${string}` | `0x${string}`[] | null )[];
  // tighten to non-null `0x${string}` for our ParsedLog shape.
  const topics = rawTopics.flatMap((t): `0x${string}`[] => {
    if (t === null) return [];
    return Array.isArray(t) ? t : [t];
  });

  const data = encodeAbiParameters(
    dataInputs as readonly AbiParameter[],
    dataInputs.map((i) => args[i.name])
  );

  return { address, topics, data, logIndex };
}
