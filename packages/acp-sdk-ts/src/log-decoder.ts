// Internal log-decoder helper. AcpAgent uses this to pull the right event
// out of a tx receipt without making the public EvmProviderAdapter interface
// depend on viem.
//
// The adapter contract only requires that `waitForReceipt` returns logs in
// emission order; this module handles the (relatively expensive) ABI decode
// here so adapters built on different libraries (ethers, web3, raw EIP-1193)
// don't have to reimplement decoding.

import { type AbiEvent, decodeEventLog, keccak256, toBytes } from "viem";
import type { ParsedLog } from "./providers/types.js";

/** Compute the topic hash for an event ABI fragment. */
export function eventTopic(event: AbiEvent): `0x${string}` {
  const types = event.inputs.map((i) => i.type).join(",");
  return keccak256(toBytes(`${event.name}(${types})`));
}

/**
 * Find the first log matching `eventName` in `abi` and decode it.
 * Returns `null` if no matching log was emitted.
 *
 * `abi` is intentionally typed as `readonly object[]` so the caller (AcpAgent)
 * can pass the same hand-rolled ABI it ships in {@link ./abi.ts} without
 * fighting viem's deeply-typed const generics.
 */
export function decodeFirstEvent<TArgs extends Record<string, unknown> = Record<string, unknown>>(
  abi: readonly object[],
  eventName: string,
  logs: readonly ParsedLog[]
): TArgs | null {
  const eventFrag = abi.find(
    (f): f is AbiEvent =>
      typeof f === "object" &&
      f !== null &&
      (f as { type?: unknown }).type === "event" &&
      (f as { name?: unknown }).name === eventName
  );
  if (!eventFrag) return null;

  const wantTopic = eventTopic(eventFrag);
  for (const log of logs) {
    if (log.topics.length === 0 || log.topics[0] !== wantTopic) continue;
    try {
      const decoded = decodeEventLog({
        // viem accepts the same fragment shape we use; the cast is a
        // narrowing convenience, not a type override.
        abi: [eventFrag] as const,
        data: log.data,
        topics: log.topics as [`0x${string}`, ...`0x${string}`[]],
      });
      // viem returns `{ args, eventName }`; we only care about args.
      return decoded.args as unknown as TArgs;
    } catch {
      // Topic matched but payload didn't decode — keep scanning. This
      // happens when the same event signature is emitted by an unrelated
      // contract in the same tx.
    }
  }
  return null;
}
