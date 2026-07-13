# rules

`rules` installs the DeBros engineering rules and skills into any project, per
project, for Claude Code. The rule block and the skills are embedded into the
binary (`go:embed`), so it needs no network at runtime.

Running `rules init` in a project:

- installs three namespaced Claude Code skills into `.claude/skills/`:
  - `rules-bugboard` — using bugboard.ai (task/bug tracker MCP)
  - `rules-orama` — deploying to and operating the Orama Network
  - `rules-rootwallet` — RootWallet and its Orama integration
- injects the DeBros hard-rules block into the project's `CLAUDE.md`, between
  `<!-- rules:start -->` / `<!-- rules:end -->` markers
- writes a `.rules.lock` manifest so future updates only touch files you have
  not edited yourself

## Install

```sh
brew install debros/tap/rules
```

Or build from source:

```sh
go build -o rules .
```

## Commands

### `rules init [--force]`

Installs the skills and injects the rules block into `CLAUDE.md` in the current
directory.

- If `CLAUDE.md` does not exist, it is created containing just the wrapped block.
- If it exists with the markers, the text between them is replaced.
- If it exists without the markers, the block is appended after a blank line and
  all existing content is preserved.
- A skill file that already exists, differs from what `rules` would write, and
  is not tracked in `.rules.lock` is treated as **your** file: it is reported as
  a conflict and skipped. Pass `--force` to overwrite conflicts. Files already
  managed by `rules` are updated normally.

### `rules update`

Refreshes managed skill files from the embedded content and re-injects the
`CLAUDE.md` block. A skill file you have edited (its on-disk hash no longer
matches the lock) is left untouched and reported. Requires an existing
`.rules.lock` — run `rules init` first. The block injection only ever touches
text between the markers.

### `rules doctor`

Reports install health: CLI version vs. installed version, whether each managed
file is present / missing / modified, and whether `CLAUDE.md` carries the rules
block. Exits non-zero on any drift or missing file.

### `rules version`

Prints the version. It is stamped at build time via
`-ldflags "-X main.version=<v>"` and defaults to `dev`.

## What `rules` manages

| Path | Managed how |
|---|---|
| `.claude/skills/rules-*/SKILL.md` | Hash-tracked in `.rules.lock`; refreshed by `update` unless you edited them |
| `CLAUDE.md` (between markers) | Marker-managed, not hash-tracked; always safe to re-inject |
| `.rules.lock` | Written by `init` / `update`; records version, timestamp, and per-file hashes |

## Releases

Releases are cut with [GoReleaser](https://goreleaser.com) (`.goreleaser.yaml`),
which builds darwin/linux amd64+arm64 binaries and publishes a `rules` formula
to the `DeBros/homebrew-tap` repository. `Formula/rules.rb` is a reference copy
of that formula.
