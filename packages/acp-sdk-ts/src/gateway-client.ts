// Typed REST client against the Titular gateway. Hand-rolled (rather than
// generated from the OpenAPI spec) to keep the npm tarball small and to keep
// auth handling under our direct control.
//
// All non-2xx responses are turned into `AcpError` with `code: "gateway_error"`
// and the upstream error body attached when the gateway returns one. Network
// failures become `gateway_unreachable`; payloads that don't match the
// expected shape become `invalid_response`.

import { AcpError } from "./errors.js";
import type {
  Agent,
  AgentsPage,
  AuthNonceResponse,
  AuthVerifyRequest,
  AuthVerifyResponse,
  BrowseParams,
  GatewayErrorEnvelope,
  JobPhase,
  JobsPage,
  Stats,
} from "./types.js";

/** Constructor parameters for {@link GatewayClient}. */
export interface GatewayClientParams {
  /** Base URL of the gateway, e.g. `https://api.titular.dev` or `http://localhost:8080`. */
  gatewayUrl: string;
  /** Optional bearer token (set after a successful SIWE/SIWS verify). */
  bearerToken?: string;
  /** Optional fetch implementation; defaults to global `fetch`. */
  fetch?: typeof fetch;
  /** Per-request timeout in ms; default 15000. */
  timeoutMs?: number;
  /** Optional `User-Agent` override. */
  userAgent?: string;
}

/** Filter accepted by {@link GatewayClient.listJobs}. */
export interface ListJobsParams {
  status?: JobPhase;
  limit?: number;
  cursor?: string;
}

/**
 * Thin typed wrapper around the gateway's REST surface. Every method maps
 * 1:1 onto a path in services/gateway-go/docs/openapi.json.
 */
export class GatewayClient {
  readonly #base: string;
  readonly #fetch: typeof fetch;
  readonly #timeoutMs: number;
  readonly #userAgent: string;
  #bearerToken: string | undefined;

  constructor(params: GatewayClientParams) {
    if (params.gatewayUrl.length === 0) {
      throw new AcpError("invalid_param", "GatewayClient: `gatewayUrl` cannot be empty");
    }
    this.#base = params.gatewayUrl.replace(/\/+$/, "");
    this.#fetch = params.fetch ?? globalThis.fetch.bind(globalThis);
    this.#timeoutMs = params.timeoutMs ?? 15_000;
    this.#userAgent = params.userAgent ?? "@titular/acp-sdk/0.1.0-alpha";
    this.#bearerToken = params.bearerToken;
  }

  /** Update the bearer token used for `Authorization` headers. */
  setBearerToken(token: string | undefined): void {
    this.#bearerToken = token;
  }

  /** Currently configured bearer token, or undefined if unauthenticated. */
  getBearerToken(): string | undefined {
    return this.#bearerToken;
  }

  /** GET /api/v1/agents — paginated registry listing. */
  async listAgents(params: BrowseParams = {}): Promise<AgentsPage> {
    const qs = new URLSearchParams();
    if (params.kind !== undefined) qs.set("kind", params.kind);
    if (params.limit !== undefined) qs.set("limit", String(params.limit));
    if (params.cursor !== undefined) qs.set("cursor", params.cursor);
    return this.#request<AgentsPage>("GET", `/api/v1/agents${qsSuffix(qs)}`);
  }

  /** GET /api/v1/agents/{id}. */
  async getAgent(id: string | number): Promise<Agent> {
    const encoded = encodeURIComponent(String(id));
    return this.#request<Agent>("GET", `/api/v1/agents/${encoded}`);
  }

  /** GET /api/v1/jobs — paginated job listing. */
  async listJobs(params: ListJobsParams = {}): Promise<JobsPage> {
    const qs = new URLSearchParams();
    if (params.status !== undefined) qs.set("status", params.status);
    if (params.limit !== undefined) qs.set("limit", String(params.limit));
    if (params.cursor !== undefined) qs.set("cursor", params.cursor);
    return this.#request<JobsPage>("GET", `/api/v1/jobs${qsSuffix(qs)}`);
  }

  /** GET /api/v1/stats — aggregate counters and last indexed block. */
  async getStats(): Promise<Stats> {
    return this.#request<Stats>("GET", "/api/v1/stats");
  }

  /** POST /auth/siwe/nonce — fetch a fresh SIWE nonce. */
  async siweNonce(): Promise<AuthNonceResponse> {
    return this.#request<AuthNonceResponse>("POST", "/auth/siwe/nonce");
  }

  /** POST /auth/siwe/verify — exchange signed SIWE message for a session. */
  async siweVerify(req: AuthVerifyRequest): Promise<AuthVerifyResponse> {
    return this.#request<AuthVerifyResponse>("POST", "/auth/siwe/verify", req);
  }

  /** POST /auth/siws/nonce — fetch a fresh SIWS nonce. */
  async siwsNonce(): Promise<AuthNonceResponse> {
    return this.#request<AuthNonceResponse>("POST", "/auth/siws/nonce");
  }

  /** POST /auth/siws/verify — exchange signed SIWS message for a session. */
  async siwsVerify(req: AuthVerifyRequest): Promise<AuthVerifyResponse> {
    return this.#request<AuthVerifyResponse>("POST", "/auth/siws/verify", req);
  }

  /** POST /auth/logout — revoke the current bearer session. */
  async logout(): Promise<void> {
    await this.#request<void>("POST", "/auth/logout", undefined, {
      expectNoContent: true,
    });
  }

  async #request<T>(
    method: "GET" | "POST",
    path: string,
    body?: unknown,
    opts?: { expectNoContent?: boolean }
  ): Promise<T> {
    const url = `${this.#base}${path}`;
    const headers: Record<string, string> = {
      Accept: "application/json",
      "User-Agent": this.#userAgent,
    };
    if (body !== undefined) headers["Content-Type"] = "application/json";
    if (this.#bearerToken !== undefined) {
      headers.Authorization = `Bearer ${this.#bearerToken}`;
    }

    const controller = new AbortController();
    const timer = setTimeout(() => controller.abort(), this.#timeoutMs);

    let res: Response;
    try {
      const init: RequestInit = {
        method,
        headers,
        signal: controller.signal,
      };
      if (body !== undefined) {
        init.body = JSON.stringify(body);
      }
      res = await this.#fetch(url, init);
    } catch (cause) {
      throw new AcpError("gateway_unreachable", `gateway request failed: ${method} ${path}`, {
        cause,
      });
    } finally {
      clearTimeout(timer);
    }

    if (!res.ok) {
      let envelope: GatewayErrorEnvelope | undefined;
      try {
        envelope = (await res.json()) as GatewayErrorEnvelope;
      } catch {
        // body was empty or non-JSON; leave envelope undefined
      }
      throw new AcpError(
        "gateway_error",
        envelope?.error?.message ?? `gateway returned ${res.status}`,
        {
          status: res.status,
          ...(envelope?.error !== undefined ? { gateway: envelope.error } : {}),
        }
      );
    }

    if (opts?.expectNoContent === true) {
      return undefined as T;
    }

    try {
      return (await res.json()) as T;
    } catch (cause) {
      throw new AcpError("invalid_response", `gateway returned non-JSON for ${method} ${path}`, {
        cause,
      });
    }
  }
}

function qsSuffix(qs: URLSearchParams): string {
  const s = qs.toString();
  return s.length === 0 ? "" : `?${s}`;
}
