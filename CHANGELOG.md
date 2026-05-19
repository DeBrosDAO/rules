# Changelog

All notable changes to DeBros Rules are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Pin your project to a specific version via `debros.json.rules.version`. The AI agent surfaces newer versions on session start but never auto-upgrades.

---

## [Unreleased]

### Added

- `DEBROS.md` §10 — new top-level "Application Security" section. Eight
  language-agnostic rules covering parameterized queries, route-level
  auth, password storage, boundary validation, secrets management,
  security headers, file uploads, and rate limiting. Earned its place
  via an external-repo audit in which §10.5 caught a missing CI
  secret-scanning step a baseline reviewer (without §10) missed.
- `DEBROS.md` §3.9 — new subsection "Verify external references
  against upstream." When the agent writes prose referencing an
  external tool's command, config syntax, API endpoint, or version
  pin, it fetches the tool's current docs before committing. Earned
  its place via an external-repo audit in which §3.9 caught two
  deprecated GitHub Action pins via live WebFetch. Includes trace-
  before-report verification methodology baked into the rule.
- `DEBROS.md` §3.8 — new subsection "No environment disclosure in
  commit messages or PRs." Agent-side rules for text destined for
  public repositories: no absolute filesystem paths, no OS / shell
  disclosure unless user-facing, no shell prompts in pasted output,
  no first-person voice. Conceptually paired with §3.9.
- `compliance/web-service.md` — new compliance file: HTTP-service
  implementation patterns for §10 (auth middleware, security
  headers, rate limiting, file uploads, boundary validation),
  language-agnostic.
- `compliance/javascript-typescript.md` §8 — JS/TS-specific library
  guidance for §10 (zod, @node-rs/argon2 or bcrypt, helmet,
  express-rate-limit, multer + file-type, gitleaks).
- `DEBROS.md` §10.0 — new "Trust model for security-critical
  configuration" preamble. Names what §3.2's untrusted-agent-output
  stance means for §10: every edit to `public_routes[]`,
  `tier3_overrides[]`, `exceptions[]`, the §3.9 verification log,
  and `.env.example` placeholders requires independent reviewer
  verification, not agent self-reporting. Added in response to an
  independent security audit that surfaced the gap.
- `DEBROS.md` §10.2 — explicit forbidden-pattern subrule for
  `public_routes[]`: root path, naked wildcards, admin/internal
  prefixes, and over-broad wildcards are rejected without inline
  justification + tracked `exceptions[]` entry.
- `DEBROS.md` §4.1 — extended sub-agent-review trigger list to
  include any edit to `public_routes[]`, `tier3_overrides[]`, or
  `exceptions[]` regardless of line count. Closes the "one-line
  JSON edit bypasses §4" loophole.
- `DEBROS.md` §10.5 — added one-time history-scan requirement at
  adoption (`gitleaks git .` or equivalent). Closes the gap between
  the rule mandating rotation upon historic exposure and the rule
  not mandating detection of historic exposure.
- `DEBROS.md` §10.8 — added mandatory trust-proxy specification.
  Rate-limit keying must pin trusted proxy hops; rotating
  `X-Forwarded-For` must not yield additional buckets. Closes the
  `req.ip` bypass surfaced by the audit.

### Changed

- `DEBROS.md` §3.9 — widened scope from "agent writes prose" to
  "agent writes or reviews prose." Codifies the behavior that
  earned the rule its place in the experiment. Verification-log
  format now requires an evidence anchor per row (commit SHA,
  version tag, or quoted line) so the log is spot-checkable
  without re-fetching every reference.
- `DEBROS.md` §10.0 — added explicit human-approval requirement
  for `tier3_overrides[]` entries (which REMOVE structural
  guarantees). Sub-agent review may accompany but does not
  substitute. Closes the recursive-trust loophole where a fully
  agent-internal review chain satisfied "human OR sub-agent"
  without involving a human.
- `DEBROS.md` §10.6 — rewritten in invariant form ("every
  response carries the configured headers"). Regression test
  matrix broadened from `401` alone to `{401, 404, 429, 413, 400}`
  to cover every middleware that can short-circuit before the
  handler.
- `DEBROS.md` §4.1, §10.2, §10.5, §10.8 — round-1 fixes rewritten
  in invariant form with one worked example each, rather than
  enumerating cases. Closes the literal-compliance loophole where
  an adopter satisfies the worked example without satisfying the
  underlying principle.

### Rationale

Rule selection driven by an external-repo audit experiment: only
rules with demonstrated catches on real third-party code earn
first-class DEBROS.md slots. Refinements driven by a follow-up
independent security audit that probed the gap between stated
intent and rule text. Operational protocols and feature-bound
rules (dismissal protocol, schema fields for bootstrap scripts,
machine-readable Tier-3 tooling) stay with the feature patches
that introduced them.

## [v0.2.0] — 2026-05-12

### Added

- `templates/agent-pointers/` — drop-in pointer files for every major AI
  coding tool (Claude Code, Cursor, GitHub Copilot Chat, Aider/Cline/Goose
  via `AGENTS.md`). Each is ~3 lines that point the tool at the canonical
  `DEBROS.md` and remind it of §3.7 (no AI co-author on commits).
- Bootstrap prompt in README now includes step 7: copy the pointer files
  into the adopted repo so AI tools find `DEBROS.md` without per-session
  setup.
- README "Wiring your AI" section rewritten to describe the pointer-file
  pattern instead of asking adopters to hand-wire each tool.
- README "Repository structure" updated to show `templates/agent-pointers/`.

### Rationale

There's no universal "AI agents read this file" convention. Each tool
looks in its own place by default. v0.1.0 left adoption-of-the-rules
to the user (paste the bootstrap prompt every session). v0.2.0 closes
that gap by establishing a pointer-file convention that works with the
defaults every major tool already uses.

## [v0.1.0] — 2026-05-12

Initial public release. Expect breaking changes before v1.0.

### Added

- `DEBROS.md` — canonical engineering and AI-agent rules. Covers supply-chain hygiene, code quality, AI-agent behavior, sub-agent review, compliance drift, and exceptions.
- `compliance/javascript-typescript.md` — JS/TS tooling baseline (`.npmrc`, `renovate.json`, lockfile, Node version pinning, vuln scanning, TypeScript strict mode).
- `compliance/go.md` — Go tooling baseline (`go.mod` toolchain directive, `go.sum`, `govulncheck`, `staticcheck`).
- `compliance/python.md` — Python tooling baseline (uv/Poetry lockfiles with hashes, `--only-binary :all:`, `pip-audit`, mypy strict, Ruff).
- `compliance/zig.md` — Zig tooling baseline (hashed `build.zig.zon`, pinned compiler with signature verification, `build.zig` review discipline).
- `DEBROS.md §3.7` — no AI co-authorship on git commits, ever. Removes `Co-Authored-By: Claude/Cursor/...` trailers and `--author` overrides.
- `DEBROS.md §8` — Agent Identity: AnBuddy. DeBros default persona for the AI agent — Spartan voice, direct, honest, light wit. Explicitly separable: other orgs adopting these rules can fork or replace this section without touching the technical rules.
- Bootstrap prompt updated to trigger the AnBuddy introduction on first activation and reinforce the no-co-author rule.
- `templates/debros.json` — per-project metadata + rules-version tracking, plus a JSON Schema (`templates/debros.schema.json`) for validation.
- `templates/.npmrc` — canonical npm config with `ignore-scripts=true` and other supply-chain hardening.
- `templates/renovate.json` — canonical Renovate config with 30-day cooldown and security-CVE override.
- `templates/github-workflows/security.yml` — auto-detecting CI workflow for JS/TS, Go, and Python.
- One-prompt adoption flow in the README — copy-paste prompt for any AI coding assistant.
- `README.md`, `CONTRIBUTING.md`, `CODE_OF_CONDUCT.md`, MIT `LICENSE`.

### Known gaps (roadmap)

- `compliance/react-native.md` — needs to cover the native (iOS Podfile, Android Gradle) side beyond what JS/TS already covers.
- Sync tool — currently manual; a CLI or GitHub Action that updates `debros.json.rules.sha` is on the roadmap.
- Per-rule enforcement tests — a way to verify a project actually meets each rule programmatically.
- Optional `compliance/rust.md` and `compliance/ruby.md` if there's demand.

---

[Unreleased]: https://github.com/DeBrosDAO/rules/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/DeBrosDAO/rules/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/DeBrosDAO/rules/releases/tag/v0.1.0
