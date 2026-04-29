import { describe, expect, it, vi } from "vitest";
import { AcpError } from "../errors.js";
import { GatewayClient } from "../gateway-client.js";

function makeFetch(handler: (url: string, init: RequestInit) => Response): typeof fetch {
  return vi.fn(async (input: RequestInfo | URL, init?: RequestInit) => {
    return handler(String(input), init ?? {});
  }) as unknown as typeof fetch;
}

function jsonResponse(body: unknown, init?: ResponseInit): Response {
  return new Response(JSON.stringify(body), {
    status: 200,
    headers: { "Content-Type": "application/json" },
    ...init,
  });
}

describe("GatewayClient.constructor", () => {
  it("rejects an empty gatewayUrl", () => {
    expect(() => new GatewayClient({ gatewayUrl: "" })).toThrow(AcpError);
  });

  it("strips trailing slashes from the base URL", async () => {
    let observed = "";
    const fetchImpl = makeFetch((url) => {
      observed = url;
      return jsonResponse({ data: [], next_cursor: "" });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://localhost:8080///",
      fetch: fetchImpl,
    });
    await client.listAgents();
    expect(observed).toBe("http://localhost:8080/api/v1/agents");
  });
});

describe("GatewayClient bearer token handling", () => {
  it("set/get cycle round-trips", () => {
    const client = new GatewayClient({ gatewayUrl: "http://x" });
    expect(client.getBearerToken()).toBeUndefined();
    client.setBearerToken("abc");
    expect(client.getBearerToken()).toBe("abc");
    client.setBearerToken(undefined);
    expect(client.getBearerToken()).toBeUndefined();
  });

  it("attaches Authorization header when a token is set", async () => {
    let auth: string | null = null;
    const fetchImpl = makeFetch((_url, init) => {
      auth = (init.headers as Record<string, string>).Authorization ?? null;
      return jsonResponse({
        total_agents: 0,
        total_jobs: 0,
        total_trades: 0,
        last_block_indexed: 0,
      });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
      bearerToken: "tok",
    });
    await client.getStats();
    expect(auth).toBe("Bearer tok");
  });

  it("does not attach Authorization when no token is set", async () => {
    let auth: string | undefined;
    const fetchImpl = makeFetch((_url, init) => {
      auth = (init.headers as Record<string, string>).Authorization;
      return jsonResponse({
        total_agents: 0,
        total_jobs: 0,
        total_trades: 0,
        last_block_indexed: 0,
      });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await client.getStats();
    expect(auth).toBeUndefined();
  });
});

describe("GatewayClient.listAgents", () => {
  it("hits /api/v1/agents with no query params by default", async () => {
    let observed = "";
    const fetchImpl = makeFetch((url) => {
      observed = url;
      return jsonResponse({ data: [], next_cursor: "" });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    const page = await client.listAgents();
    expect(observed).toBe("http://x/api/v1/agents");
    expect(page).toEqual({ data: [], next_cursor: "" });
  });

  it("encodes kind, limit, and cursor", async () => {
    let observed = "";
    const fetchImpl = makeFetch((url) => {
      observed = url;
      return jsonResponse({ data: [], next_cursor: "" });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await client.listAgents({ kind: "acp", limit: 10, cursor: "c1" });
    expect(observed).toBe("http://x/api/v1/agents?kind=acp&limit=10&cursor=c1");
  });
});

describe("GatewayClient.getAgent", () => {
  it("URL-encodes the id", async () => {
    let observed = "";
    const fetchImpl = makeFetch((url) => {
      observed = url;
      return jsonResponse({});
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await client.getAgent("123");
    expect(observed).toBe("http://x/api/v1/agents/123");
  });
});

describe("GatewayClient.listJobs", () => {
  it("encodes status, limit, cursor", async () => {
    let observed = "";
    const fetchImpl = makeFetch((url) => {
      observed = url;
      return jsonResponse({ data: [], next_cursor: "" });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await client.listJobs({ status: "active", limit: 5, cursor: "c2" });
    expect(observed).toBe("http://x/api/v1/jobs?status=active&limit=5&cursor=c2");
  });
});

describe("GatewayClient auth methods", () => {
  it("siweNonce hits the right endpoint with POST", async () => {
    let observed = "";
    let method = "";
    const fetchImpl = makeFetch((url, init) => {
      observed = url;
      method = init.method ?? "GET";
      return jsonResponse({ nonce: "n", expires_at: 1 });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    const r = await client.siweNonce();
    expect(observed).toBe("http://x/auth/siwe/nonce");
    expect(method).toBe("POST");
    expect(r.nonce).toBe("n");
  });

  it("siweVerify posts the message + signature", async () => {
    let body = "";
    const fetchImpl = makeFetch((_url, init) => {
      body = String(init.body ?? "");
      return jsonResponse({ address: "0x", expires_at: 1, token: "t" });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await client.siweVerify({ message: "msg", signature: "sig" });
    expect(JSON.parse(body)).toEqual({ message: "msg", signature: "sig" });
  });

  it("siwsNonce + siwsVerify hit the SIWS routes", async () => {
    const seen: string[] = [];
    const fetchImpl = makeFetch((url) => {
      seen.push(url);
      return jsonResponse({
        nonce: "n",
        expires_at: 1,
        address: "0x",
        token: "t",
      });
    });
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await client.siwsNonce();
    await client.siwsVerify({ message: "m", signature: "s" });
    expect(seen).toEqual(["http://x/auth/siws/nonce", "http://x/auth/siws/verify"]);
  });

  it("logout returns void on no-content responses", async () => {
    const fetchImpl = makeFetch(() => new Response(null, { status: 204 }));
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await expect(client.logout()).resolves.toBeUndefined();
  });
});

describe("GatewayClient error handling", () => {
  it("wraps a structured 4xx body into AcpError(gateway_error)", async () => {
    const fetchImpl = makeFetch(
      () =>
        new Response(
          JSON.stringify({
            error: { code: "not_found", message: "no such agent" },
          }),
          {
            status: 404,
            headers: { "Content-Type": "application/json" },
          }
        )
    );
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await expect(client.getAgent(1)).rejects.toMatchObject({
      code: "gateway_error",
      status: 404,
      gateway: { code: "not_found", message: "no such agent" },
    });
  });

  it("falls back to a generic message when the body is not the expected shape", async () => {
    const fetchImpl = makeFetch(
      () =>
        new Response("oops", {
          status: 500,
          headers: { "Content-Type": "text/plain" },
        })
    );
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await expect(client.getStats()).rejects.toMatchObject({
      code: "gateway_error",
      status: 500,
    });
  });

  it("turns network failures into AcpError(gateway_unreachable)", async () => {
    const fetchImpl = vi
      .fn()
      .mockRejectedValue(new TypeError("network down")) as unknown as typeof fetch;
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await expect(client.getStats()).rejects.toMatchObject({
      code: "gateway_unreachable",
    });
  });

  it("turns non-JSON 200 responses into AcpError(invalid_response)", async () => {
    const fetchImpl = makeFetch(
      () =>
        new Response("not json", {
          status: 200,
          headers: { "Content-Type": "text/plain" },
        })
    );
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
    });
    await expect(client.getStats()).rejects.toMatchObject({
      code: "invalid_response",
    });
  });

  it("aborts the request when timeoutMs elapses", async () => {
    // Fetch implementation that never resolves but honours AbortSignal.
    const fetchImpl = vi.fn((_url: string, init?: { signal?: AbortSignal }) => {
      return new Promise((_resolve, reject) => {
        init?.signal?.addEventListener("abort", () => reject(new Error("aborted")));
      });
    }) as unknown as typeof fetch;
    const client = new GatewayClient({
      gatewayUrl: "http://x",
      fetch: fetchImpl,
      timeoutMs: 5,
    });
    await expect(client.getStats()).rejects.toMatchObject({
      code: "gateway_unreachable",
    });
  });
});
