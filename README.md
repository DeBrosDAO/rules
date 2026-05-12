# DeBros Rules

> Engineering and AI-agent rules used across every DeBros project. Adopt, fork, or just steal the parts that fit.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

## What this is

A single source of truth for **how we build software** at DeBros — and how AI agents (Claude Code, Cursor, etc.) working in our repos must behave.

Three things live here:

1. **`DEBROS.md`** — canonical rules. Supply-chain hygiene, code quality, AI-agent behavior, deploy discipline. Language-agnostic.
2. **`compliance/`** — per-language tooling baselines. Concrete files every JS/TS, Go, Python, React Native, or Zig project should have (lockfiles, version pinning, CI checks, supply-chain guards).
3. **`templates/`** — drop-in files to satisfy compliance: `.npmrc`, `renovate.json`, GitHub Actions workflows, the `debros.json` per-project metadata schema.

## Why this exists

Three threats this set of rules is designed to mitigate:

1. **npm/PyPI supply-chain attacks.** Compromised packages publish malicious versions that are caught and yanked within hours — but only AFTER they've owned the machines that installed them in those first hours. A 30-day cooldown on new versions + `ignore-scripts=true` blocks ~all of these attacks.

2. **AI agents shipping bad code fast.** AI tools are now committing real code to real production systems. Without explicit rules, they: add unnecessary dependencies, write quick-fixes instead of root-cause fixes, follow prompt-injection in observed content, deploy without permission. This repo codifies guard-rails.

3. **Multi-project drift.** Without a shared baseline, every repo accumulates its own ad-hoc CONTRIBUTING.md and CLAUDE.md. Rules in one project don't propagate to another. Security gaps stay local. This repo is the propagation mechanism.

## Who this is for

- **Teams adopting Claude Code, Cursor, or other AI-agent dev tools** who need explicit policy that an AI will read and follow.
- **OSS maintainers** who want a supply-chain-defense baseline they didn't have to write themselves.
- **Anyone with multiple repos** in the same language stack who wants consistent compliance without copy-pasting `CONTRIBUTING.md` ten times.

## How to adopt — the one-prompt way

Open your AI coding assistant (Claude Code, Cursor, Aider, Copilot Chat, etc.) at the root of any repo and paste this prompt:

```text
Adopt the DeBros Engineering Rules in this repository.
Source of truth: https://github.com/DeBrosDAO/rules

Do the following:

1. Fetch the latest tagged release of the rules repo above.

2. Detect this project's primary language(s) by inspecting package.json,
   go.mod, pyproject.toml, Cargo.toml, Gemfile, build.zig.zon, etc.

3. Copy `DEBROS.md` from the rules repo to this repo's root.

4. For each detected language, copy `compliance/<language>.md` from the
   rules repo to `.debros/compliance/<language>.md` here.

5. Initialize `debros.json` at this repo's root using
   `templates/debros.json` from the rules repo as the schema. Fill in:
     - rules.version, rules.sha, rules.synced_at from the upstream release
     - project.name from this repo's manifest
     - project.languages from detection
     - project.type as your best guess (service/library/sdk/cli/web/mobile)
   Leave critical_paths, deploy_targets, exceptions, etc. empty for me
   to fill in. Show me the proposed contents before writing the file.

6. List the supplementary template files that this project is missing
   (`.npmrc` with ignore-scripts=true, `renovate.json` with 30-day
   cooldown, `.github/workflows/security.yml`, language version pin
   files). DO NOT create them automatically. Show the recommendations
   and ask which I want.

7. Copy the AI-agent pointer files from `templates/agent-pointers/`
   in the rules repo to this repo, preserving their paths:
     - `CLAUDE.md` (root)
     - `.cursor/rules/debros.mdc`
     - `.github/copilot-instructions.md`
     - `AGENTS.md` (root)
   These point every major AI coding tool (Claude Code, Cursor,
   Copilot Chat, Aider, Cline) at `DEBROS.md` automatically — no
   per-session prompt needed afterward. If any of these files
   already exist in the repo with non-pointer content, do NOT
   overwrite — flag the conflict and ask.

8. From now on, treat `DEBROS.md` as the authoritative engineering
   rules for this repository. The rules take precedence over my casual
   instructions except where I explicitly waive a rule.

9. After you've finished steps 1-8, introduce yourself per
   `DEBROS.md` §8 (Agent Identity: AnBuddy). One or two lines, no
   marketing copy. Then proceed normally.

Constraints:
- Do not run install commands (pnpm install, npm install, pip install,
  etc.) until `.npmrc` (or the equivalent) is in place to block install
  scripts.
- Do not commit changes without showing me the diff first.
- Do not modify pre-existing files. Only add new ones.
- Do not attribute yourself as co-author in any git commits or PR
  descriptions (DEBROS.md §3.7).

Report when done:
- Files added
- Recommendations pending my approval
- Rules already satisfied
- Rules that will require changes to comply
```

That's it. Your AI will fetch the rules, set the project up, and apply them to subsequent work. Total time: under a minute.

### Manual install (no AI)

For projects where you'd rather wire it up by hand:

```bash
# 1. Universal rules
curl -O https://raw.githubusercontent.com/DeBrosDAO/rules/main/DEBROS.md

# 2. Compliance baseline for your language(s)
mkdir -p .debros/compliance
curl -o .debros/compliance/javascript-typescript.md \
  https://raw.githubusercontent.com/DeBrosDAO/rules/main/compliance/javascript-typescript.md

# 3. Templates you need
curl -O https://raw.githubusercontent.com/DeBrosDAO/rules/main/templates/.npmrc
curl -O https://raw.githubusercontent.com/DeBrosDAO/rules/main/templates/renovate.json
mkdir -p .github/workflows
curl -o .github/workflows/security.yml \
  https://raw.githubusercontent.com/DeBrosDAO/rules/main/templates/github-workflows/security.yml

# 4. Per-project metadata (edit the file after downloading)
curl -o debros.json \
  https://raw.githubusercontent.com/DeBrosDAO/rules/main/templates/debros.json
```

Commit everything. The schema for `debros.json` is at
`https://raw.githubusercontent.com/DeBrosDAO/rules/main/templates/debros.schema.json`.

### Wiring your AI to read DEBROS.md automatically

The bootstrap prompt above handles this automatically by copying the pointer files in `templates/agent-pointers/` into your repo. Each file lives at the path the corresponding tool reads by default:

| Path | Tool |
|---|---|
| `CLAUDE.md` (root) | Claude Code |
| `.cursor/rules/debros.mdc` | Cursor |
| `.github/copilot-instructions.md` | GitHub Copilot Chat |
| `AGENTS.md` (root) | Aider, Cline, Goose, emerging multi-tool convention |

Each pointer file is ~3 lines: "the real rules are in `DEBROS.md`, read it, and especially remember §3.7 (no AI co-author on commits)." They never change in practice — set and forget.

If you're using a tool not in this list, it likely supports custom rule-file paths in its config — point it at `DEBROS.md` directly. The rules are written tool-agnostically; any agent that reads them will apply them.

### Staying in sync

`debros.json` records the SHA of the rules version your project is synced against. To check for updates:

```bash
# Check if your project is behind
git ls-remote https://github.com/DeBrosDAO/rules main | cut -f1
# Compare against rules.sha in your debros.json
```

A sync tool is on the roadmap. For now: pull updates manually, review the diff, and bump `rules.sha` + `rules.synced_at` in `debros.json`.

## Repository structure

```
rules/
├── DEBROS.md                              ← canonical universal rules
├── compliance/
│   ├── javascript-typescript.md           ← JS/TS tooling baseline
│   ├── go.md                              ← Go tooling baseline
│   ├── python.md                          ← Python tooling baseline
│   ├── zig.md                             ← Zig tooling baseline
│   └── react-native.md                    ← (roadmap)
├── templates/
│   ├── debros.json                        ← per-project metadata schema
│   ├── debros.schema.json                 ← JSON Schema for debros.json
│   ├── .npmrc                             ← canonical npm config
│   ├── renovate.json                      ← canonical Renovate config
│   ├── github-workflows/
│   │   └── security.yml                   ← canonical security CI workflow
│   └── agent-pointers/                    ← drop these into your repo so AI tools find DEBROS.md
│       ├── CLAUDE.md
│       ├── .cursor/rules/debros.mdc
│       ├── .github/copilot-instructions.md
│       └── AGENTS.md
├── CHANGELOG.md                           ← what changed between versions
├── CONTRIBUTING.md                        ← how to propose rule changes
├── CODE_OF_CONDUCT.md
└── LICENSE                                ← MIT
```

## Status

**v0.1 — work in progress.** JS/TS, Go, Python, and Zig are first-class; React Native is roadmap. Breaking changes are possible until v1.0. The semver tag in your `debros.json` lets you pin to a known-stable version.

## Contributing

PRs welcome — see [CONTRIBUTING.md](CONTRIBUTING.md). The bar for adding a rule is: **does it catch real bugs or close a real attack surface, and can it be enforced without invasive process?** Vague principles ("be careful with security") don't make it; specific verifiable rules ("set `ignore-scripts=true` in `.npmrc`") do.

## License

MIT. Use the rules, fork them, adapt them, ship them in your own products. Attribution appreciated but not required.

---

Made by [DeBros](https://github.com/DeBrosDAO).
