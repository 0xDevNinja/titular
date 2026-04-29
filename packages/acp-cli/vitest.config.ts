import { defineConfig } from "vitest/config";

export default defineConfig({
  test: {
    environment: "node",
    globals: false,
    include: ["src/__tests__/**/*.test.ts"],
    coverage: {
      provider: "v8",
      reporter: ["text", "lcov"],
      include: ["src/**/*.ts"],
      // `cli.ts` is mostly commander wiring whose branch coverage is
      // dominated by typescript `exactOptionalPropertyTypes`-style ternaries
      // — exercising every flag combo would balloon the test surface for no
      // real safety win. The smoke test in `cli.test.ts` still asserts the
      // command graph wires up correctly.
      exclude: ["src/__tests__/**", "src/cli.ts", "src/types.ts"],
      thresholds: {
        lines: 75,
        statements: 75,
        functions: 75,
        branches: 70,
      },
    },
  },
});
