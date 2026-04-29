// Error taxonomy. The SDK throws structured errors so consumers can branch on
// `code` rather than scraping `.message`. Codes are stable; messages are not.

import type { GatewayErrorBody } from "./types.js";

/** Stable, machine-readable error codes surfaced by the SDK. */
export type AcpErrorCode =
  | "gateway_error"
  | "gateway_unreachable"
  | "invalid_response"
  | "invalid_param"
  | "no_signer"
  | "tx_reverted"
  | "event_not_found"
  | "provider_cannot_sign"
  | "auth_failed";

/**
 * Base error class for everything thrown by the SDK. Always carries a
 * stable `code` and, where applicable, the underlying gateway error body.
 */
export class AcpError extends Error {
  override readonly name = "AcpError";
  readonly code: AcpErrorCode;
  readonly status?: number;
  readonly gateway?: GatewayErrorBody;
  override readonly cause?: unknown;

  constructor(
    code: AcpErrorCode,
    message: string,
    options?: { status?: number; gateway?: GatewayErrorBody; cause?: unknown }
  ) {
    super(message);
    this.code = code;
    if (options?.status !== undefined) this.status = options.status;
    if (options?.gateway !== undefined) this.gateway = options.gateway;
    if (options?.cause !== undefined) this.cause = options.cause;
  }
}
