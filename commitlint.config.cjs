module.exports = {
  extends: ["@commitlint/config-conventional"],
  rules: {
    "subject-case": [2, "always", "lower-case"],
    "header-max-length": [2, "always", 72],
    "type-enum": [
      2,
      "always",
      [
        "feat",
        "fix",
        "refactor",
        "test",
        "docs",
        "chore",
        "perf",
        "style",
        "ci",
        "build",
        "revert",
      ],
    ],
  },
};
