module.exports = {
  extends: ["@commitlint/cli", "@commitlint/config-conventional"],
  rules: {
    "scope-enum": [
      2,
      "always",
      ["root", "api", "web", "ui", "zod", "openapi", "templates"],
    ],
    "type-enum": [
      2,
      "always",
      [
        "feat",
        "fix",
        "docs",
        "style",
        "refactor",
        "perf",
        "test",
        "build",
        "ci",
        "chore",
        "revert",
      ],
    ],
    "scope-empty": [2, "never"],
    "subject-case": [2, "always", "sentence-case"],
  },
};
