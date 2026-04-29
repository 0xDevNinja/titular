import { describe, expect, it } from "vitest";
import { TOOL_NAMES, availableTools } from "../llm/tools.js";

describe("availableTools", () => {
  it("returns one tool per AcpAgent action", () => {
    const tools = availableTools();
    const names = tools.map((t) => t.function.name).sort();
    expect(names).toEqual([TOOL_NAMES.browse, TOOL_NAMES.createJob, TOOL_NAMES.register].sort());
  });

  it("emits OpenAI-shaped envelopes", () => {
    for (const tool of availableTools()) {
      expect(tool.type).toBe("function");
      expect(typeof tool.function.name).toBe("string");
      expect(typeof tool.function.description).toBe("string");
      expect(tool.function.description.length).toBeGreaterThan(20);
      expect(tool.function.parameters.type).toBe("object");
    }
  });

  it("returns a fresh array each call (mutation safety)", () => {
    const a = availableTools();
    const b = availableTools();
    expect(a).not.toBe(b);
    a.pop();
    expect(b.length).toBe(3);
  });

  it("filters with `include`", () => {
    const tools = availableTools({ include: [TOOL_NAMES.browse] });
    expect(tools.map((t) => t.function.name)).toEqual([TOOL_NAMES.browse]);
  });

  it("filters with `include` set to multiple", () => {
    const tools = availableTools({
      include: [TOOL_NAMES.browse, TOOL_NAMES.register],
    });
    const names = tools.map((t) => t.function.name).sort();
    expect(names).toEqual([TOOL_NAMES.browse, TOOL_NAMES.register].sort());
  });

  it("returns an empty array when `include` matches nothing", () => {
    // @ts-expect-error: deliberately passing an unknown name to verify behavior
    const tools = availableTools({ include: ["does_not_exist"] });
    expect(tools).toEqual([]);
  });

  describe("acp_register schema", () => {
    const tool = availableTools().find((t) => t.function.name === TOOL_NAMES.register);

    it("requires metadataUri + capabilities", () => {
      expect(tool?.function.parameters.required).toEqual(["metadataUri", "capabilities"]);
    });

    it("describes the controller field", () => {
      const props = tool?.function.parameters.properties;
      expect(props?.controller?.pattern).toBe("^0x[0-9a-fA-F]{40}$");
    });

    it("constrains capabilities to decimal uint", () => {
      const props = tool?.function.parameters.properties;
      expect(props?.capabilities?.pattern).toBe("^(0|[1-9][0-9]*)$");
    });

    it("disallows additional properties", () => {
      expect(tool?.function.parameters.additionalProperties).toBe(false);
    });
  });

  describe("acp_browse schema", () => {
    const tool = availableTools().find((t) => t.function.name === TOOL_NAMES.browse);

    it("has no required fields", () => {
      expect(tool?.function.parameters.required).toBeUndefined();
    });

    it("constrains kind to enum", () => {
      const props = tool?.function.parameters.properties;
      expect(props?.kind?.enum).toEqual(["launchpad", "acp"]);
    });

    it("constrains limit to [1, 200]", () => {
      const props = tool?.function.parameters.properties;
      expect(props?.limit?.type).toBe("integer");
      expect(props?.limit?.minimum).toBe(1);
      expect(props?.limit?.maximum).toBe(200);
    });
  });

  describe("acp_create_job schema", () => {
    const tool = availableTools().find((t) => t.function.name === TOOL_NAMES.createJob);

    it("requires token, budget, deadline", () => {
      expect(tool?.function.parameters.required).toEqual(["token", "budget", "deadline"]);
    });

    it("constrains jobType enum", () => {
      const props = tool?.function.parameters.properties;
      expect(props?.jobType?.enum).toEqual(["direct", "evaluated"]);
    });

    it("typed hooks as array of addresses", () => {
      const props = tool?.function.parameters.properties;
      expect(props?.hooks?.type).toBe("array");
      expect(props?.hooks?.items?.pattern).toBe("^0x[0-9a-fA-F]{40}$");
    });
  });
});
