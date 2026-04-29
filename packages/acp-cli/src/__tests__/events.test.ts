import { describe, expect, it } from "vitest";

import { EventsError, __testing__, parseEventBlock, runEventsListen } from "../commands/events.js";
import { saveConfig } from "../config.js";

import { CollectingStream, fakeFetchSse, freshConfigPath } from "./_helpers.js";

const BASE_CFG = {
  gatewayUrl: "https://api.titular.dev",
  chain: "base-sepolia" as const,
  apiKey: "k",
};

describe("parseEventBlock", () => {
  it("parses an event with id, event and data", () => {
    const evt = parseEventBlock('id: 7\nevent: agent_registered\ndata: {"x":1}');
    expect(evt).not.toBeNull();
    expect(evt?.id).toBe("7");
    expect(evt?.event).toBe("agent_registered");
    expect(evt?.data).toBe('{"x":1}');
  });

  it("returns null on a comment-only block", () => {
    expect(parseEventBlock(":\n")).toBeNull();
    expect(parseEventBlock(": heartbeat")).toBeNull();
  });

  it("joins multi-line data fields with newlines", () => {
    const evt = parseEventBlock("data: line1\ndata: line2");
    expect(evt?.data).toBe("line1\nline2");
  });

  it("ignores unknown fields", () => {
    const evt = parseEventBlock("retry: 1000\ndata: x");
    expect(evt?.data).toBe("x");
  });
});

describe("findEventBoundary", () => {
  const { findEventBoundary } = __testing__;
  it("returns -1 when no full event is buffered", () => {
    expect(findEventBoundary("data: incomplete")).toBe(-1);
  });

  it("finds \\n\\n boundaries", () => {
    expect(findEventBoundary("data: a\n\nrest")).toBe(7);
  });

  it("finds \\r\\n\\r\\n boundaries", () => {
    expect(findEventBoundary("data: a\r\n\r\nrest")).toBe(7);
  });
});

describe("buildEventsUrl", () => {
  const { buildEventsUrl } = __testing__;
  it("strips trailing slashes from base", () => {
    expect(buildEventsUrl({ gatewayUrl: "https://api/", chain: "base-sepolia" })).toBe(
      "https://api/events"
    );
  });

  it("encodes subjects as a query param", () => {
    const url = buildEventsUrl(
      { gatewayUrl: "https://api", chain: "base-sepolia" },
      "agent.registered,job.created"
    );
    expect(url).toContain("subjects=agent.registered%2Cjob.created");
  });
});

describe("runEventsListen", () => {
  it("rejects when no config is present", async () => {
    const { path } = await freshConfigPath();
    await expect(
      runEventsListen({ config: path }, { env: {}, isTTY: false })
    ).rejects.toBeInstanceOf(EventsError);
  });

  it("decodes a small SSE stream and prints text by default", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const stdout = new CollectingStream();
    const events = await runEventsListen(
      { config: path, maxEvents: 2 },
      {
        fetch: fakeFetchSse([
          "id: 1\nevent: agent_registered\ndata: hello\n\n",
          "id: 2\nevent: job_created\ndata: world\n\n",
        ]),
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(events).toHaveLength(2);
    expect(stdout.text()).toContain("[agent_registered]");
    expect(stdout.text()).toContain("hello");
    expect(stdout.text()).toContain("[job_created]");
  });

  it("emits JSON envelopes when --json is set, parsing data when valid JSON", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const stdout = new CollectingStream();
    await runEventsListen(
      { config: path, maxEvents: 1, json: true },
      {
        fetch: fakeFetchSse(['id: 9\nevent: x\ndata: {"a":1}\n\n']),
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    const text = stdout.text();
    expect(text).toContain('"id": "9"');
    expect(text).toContain('"event": "x"');
    expect(text).toContain('"a": 1');
  });

  it("falls back to raw string when --json data is not valid JSON", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const stdout = new CollectingStream();
    await runEventsListen(
      { config: path, maxEvents: 1, json: true },
      {
        fetch: fakeFetchSse(["data: hello not-json\n\n"]),
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: {},
        isTTY: false,
      }
    );
    expect(stdout.text()).toContain('"data": "hello not-json"');
  });

  it("treats heartbeat-only chunks as no-ops", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const stdout = new CollectingStream();
    const events = await runEventsListen(
      { config: path, maxEvents: 1 },
      {
        fetch: fakeFetchSse([": heartbeat\n\n", "id: 1\ndata: real\n\n"]),
        stdout: stdout as unknown as NodeJS.WriteStream,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(events).toHaveLength(1);
    expect(stdout.text()).toContain("real");
  });

  it("surfaces a typed error on non-2xx responses", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    const fetchImpl = ((..._args: unknown[]) =>
      Promise.resolve(
        new Response("nope", { status: 503, statusText: "unavailable" })
      )) as typeof fetch;
    await expect(
      runEventsListen(
        { config: path, maxEvents: 1 },
        {
          fetch: fetchImpl,
          env: { NO_COLOR: "1" },
          isTTY: false,
        }
      )
    ).rejects.toBeInstanceOf(EventsError);
  });

  it("threads --subjects and --last-event-id into the request", async () => {
    const { path } = await freshConfigPath();
    await saveConfig(path, BASE_CFG);
    let capturedUrl = "";
    let capturedHeaders: Record<string, string> = {};
    const fetchImpl = ((url: string, init: RequestInit) => {
      capturedUrl = url;
      capturedHeaders = init.headers as Record<string, string>;
      const body = new ReadableStream<Uint8Array>({
        start(controller) {
          controller.close();
        },
      });
      return Promise.resolve(new Response(body, { status: 200 }));
    }) as unknown as typeof fetch;
    await runEventsListen(
      {
        config: path,
        subjects: "agent.registered",
        lastEventId: "abc",
      },
      {
        fetch: fetchImpl,
        env: { NO_COLOR: "1" },
        isTTY: false,
      }
    );
    expect(capturedUrl).toContain("subjects=agent.registered");
    expect(capturedHeaders["Last-Event-ID"]).toBe("abc");
    expect(capturedHeaders.Authorization).toBe("Bearer k");
  });
});
