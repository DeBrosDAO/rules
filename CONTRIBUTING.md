# Contributing to DeBros Rules

Thanks for considering a contribution. This repo is the source of truth for how DeBros builds software and how AI agents working in our repos must behave — but the rules are deliberately open-source because the underlying problems (supply-chain attacks, AI-agent misbehavior, multi-project drift) aren't specific to us.

## What kinds of changes are accepted

**Good candidates:**

- **New rules that catch real bugs or close real attack surfaces.** Cite a real incident or CVE class. Vague principles ("be careful with concurrency") don't make the bar; concrete verifiable rules ("Go projects must run `go test -race` in CI for any package using goroutines") do.
- **New per-language compliance baselines.** Python, React Native, Zig, Rust, Ruby — all welcome. Follow the structure of `compliance/javascript-typescript.md` and `compliance/go.md`.
- **Template improvements** — better `.npmrc`, better `renovate.json` defaults, additional CI workflows.
- **Bug fixes** in templates (e.g. a workflow that doesn't trigger on the right events, a config that conflicts with another).
- **Clarifications** to existing rules where the wording is ambiguous in practice.

**Not accepted:**

- **Removing rules** to make them more permissive without a strong argument for why the rule was wrong. Default is to keep the rule.
- **Adding rules without enforcement.** "Be careful with X" is not a rule; "Add lint check X to CI" is.
- **DeBros-internal operational details.** Project-specific deploy procedures, internal infrastructure references, customer integrations belong in each repo's own `.claude/rules/`, not here.
- **Personal style preferences.** "I think 2-space indent is better" doesn't belong here; "indent style should be enforced by formatter X" might.

## How to propose a change

1. **Open an issue first** for anything beyond a typo fix. State the problem, the proposed rule, and at least one real incident or CVE class that the rule prevents.
2. **Wait for discussion.** Maintainers will weigh in on scope and approach before you write the PR.
3. **Submit a PR** referencing the issue. Keep PRs focused — one rule change per PR.
4. **The PR description must include:**
   - What changed and why
   - The real-world incident or risk this addresses
   - How the rule is enforced (linter, CI check, sub-agent review, documentation only)
   - Migration impact: what does an existing project have to do to comply?

## Versioning

Releases follow semver:

- **Major** — breaking changes to the `debros.json` schema or to the meaning of Tier-3 blocks. Projects pinned to a previous major version keep working.
- **Minor** — new rules, new compliance baselines, new templates.
- **Patch** — editorial fixes, typo corrections, clarifications that don't change rule meaning.

A new release is tagged via `git tag v1.2.3 && git push --tags` and listed in `CHANGELOG.md`.

## Style guide for rule prose

- **Active voice.** "The agent must..." not "It is required that the agent..."
- **Concrete, not aspirational.** "Functions ≤50 lines" not "Keep functions short."
- **Cite enforcement.** Every rule names how it's enforced — linter, CI, sub-agent review, or manual review.
- **Explain why briefly.** A one-sentence rationale per rule. Not a treatise.
- **Cross-link.** Reference other rules by section number when relevant.

## Testing your changes

For rules in `DEBROS.md` or `compliance/*.md`, "testing" means:

- Read the rule out loud. Does it pass the "concrete, enforceable, cites why" bar?
- Apply the rule to an existing DeBros project (or a public OSS project of similar shape). Does the rule make sense in practice? Does it catch a real issue?
- If the rule has a Tier-3 block, walk through a scenario where the block triggers. Is the block actionable (clear message, obvious fix) or annoying?

For templates (`.npmrc`, `renovate.json`, workflows), copy them into a fresh project and confirm they work end-to-end.

## Code of Conduct

This project follows the [Contributor Covenant](CODE_OF_CONDUCT.md). Disrespectful behavior in issues or PRs results in immediate ban with no warning.

## License

By contributing, you agree your contributions are licensed under the [MIT License](LICENSE) the project uses. Don't submit code you don't have the right to relicense.
