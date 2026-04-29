import { describe, expect, it } from "vitest";

import { buildProgram, main } from "../cli.js";

describe("buildProgram", () => {
  it("registers every documented command", () => {
    const program = buildProgram();
    const names = program.commands.map((c) => c.name());
    for (const expected of ["configure", "agent", "browse", "job", "events"]) {
      expect(names).toContain(expected);
    }
  });

  it("declares the global --json flag", () => {
    const program = buildProgram();
    const opts = program.options.map((o) => o.long);
    expect(opts).toContain("--json");
    expect(opts).toContain("--config");
  });

  it("nests agent subcommands", () => {
    const program = buildProgram();
    const agent = program.commands.find((c) => c.name() === "agent");
    const sub = agent?.commands.map((c) => c.name()) ?? [];
    expect(sub).toEqual(expect.arrayContaining(["create", "list", "status"]));
  });
});

describe("main", () => {
  it("returns 0 for --help", async () => {
    const code = await main(["--help"]);
    expect(code).toBe(0);
  });

  it("returns 0 for --version", async () => {
    const code = await main(["--version"]);
    expect(code).toBe(0);
  });
});
