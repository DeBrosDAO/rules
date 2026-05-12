# Agent Pointer Files

Drop these into any DeBros-adopted repo so AI coding tools find `DEBROS.md`
without per-session setup.

Each file lives at the path the corresponding tool looks for by default:

| File | Tool |
|---|---|
| `CLAUDE.md` | Claude Code |
| `.cursor/rules/debros.mdc` | Cursor |
| `.github/copilot-instructions.md` | GitHub Copilot Chat |
| `AGENTS.md` | Aider, Cline, Goose, and the emerging multi-tool convention |

All four files have essentially the same 3-line content: "the real rules
are in `DEBROS.md`, go read it, and especially remember §3.7 (no AI
co-author on commits)."

## Why all four

There's no universal "AI agents read this file" convention today. Each
tool looks in its own place. Dropping all four files means whichever tool
a contributor runs, they'll be pointed at `DEBROS.md` automatically.

If you only use one or two of those tools, you can omit the others — but
the pointer files are 3 lines each, never change, so the cost of having
all four is negligible.

## What you do NOT do

- Don't put the actual rules content in these files. The rules live in
  `DEBROS.md`. These are pointers, not copies.
- Don't symlink. Plain text files work everywhere; symlinks break on
  Windows and behave inconsistently across tools.
- Don't add `Co-Authored-By: Claude/Cursor/...` to commits. That's the
  one rule we restate in every pointer file because tool defaults
  frequently violate it.
