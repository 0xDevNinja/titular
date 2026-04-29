import { describe, expect, it } from "vitest";
import { AcpError } from "../errors.js";

describe("AcpError", () => {
  it("preserves code, status, gateway body, and cause", () => {
    const cause = new Error("boom");
    const err = new AcpError("gateway_error", "msg", {
      status: 503,
      gateway: { code: "upstream_down", message: "indexer is down" },
      cause,
    });
    expect(err.name).toBe("AcpError");
    expect(err.code).toBe("gateway_error");
    expect(err.status).toBe(503);
    expect(err.gateway?.code).toBe("upstream_down");
    expect(err.cause).toBe(cause);
    expect(err.message).toBe("msg");
  });

  it("omits optional fields when not supplied", () => {
    const err = new AcpError("invalid_param", "bad");
    expect(err.status).toBeUndefined();
    expect(err.gateway).toBeUndefined();
    expect(err.cause).toBeUndefined();
  });

  it("is an instance of Error", () => {
    const err = new AcpError("auth_failed", "denied");
    expect(err).toBeInstanceOf(Error);
  });
});
