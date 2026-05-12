# Changelog

All notable changes to DeBros Rules are documented here. The format follows [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

Pin your project to a specific version via `debros.json.rules.version`. The AI agent surfaces newer versions on session start but never auto-upgrades.

---

## [Unreleased]

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

[Unreleased]: https://github.com/DeBrosDAO/rules/compare/v0.1.0...HEAD
[v0.1.0]: https://github.com/DeBrosDAO/rules/releases/tag/v0.1.0
