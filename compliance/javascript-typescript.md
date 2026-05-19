# Compliance — JavaScript / TypeScript

> The concrete files every JS/TS project must have to satisfy [DEBROS.md](../DEBROS.md). Applies to Node, Bun, Deno, and React Native (RN has its own [supplementary file](react-native.md) for the native side).

---

## Required files

### 1. `.npmrc` — block install-time scripts

**Tier 3 block.** Without this file, the agent refuses to run `pnpm install` or `npm install`.

Copy [`templates/.npmrc`](../templates/.npmrc) to the repo root.

Minimum contents:

```ini
# Block postinstall / preinstall / install scripts by default.
# Packages that genuinely need them (esbuild, sharp, sqlite) must be
# allowlisted in package.json `pnpm.onlyBuiltDependencies`.
ignore-scripts=true

# Fail audits at moderate severity or higher.
audit-level=moderate

# Don't install peer dependencies automatically — explicit is better.
auto-install-peers=false

# Prefer offline cache when available (reproducibility).
prefer-offline=true

# Block packages from manipulating the lockfile shape.
strict-peer-dependencies=true
```

For repos that need a few packages with install scripts, allowlist them in `package.json`:

```json
{
  "pnpm": {
    "onlyBuiltDependencies": [
      "esbuild",
      "sharp"
    ]
  }
}
```

Reviewing this allowlist counts as a security-sensitive code change (sub-agent review required per DEBROS.md §4).

### 2. `renovate.json` — enforce 30-day cooldown

Copy [`templates/renovate.json`](../templates/renovate.json) to the repo root.

Key configuration:

```jsonc
{
  "$schema": "https://docs.renovatebot.com/renovate-schema.json",
  "extends": ["config:recommended"],
  "minimumReleaseAge": "30 days",
  "automerge": false,
  "vulnerabilityAlerts": {
    "minimumReleaseAge": "0 days",
    "labels": ["security"]
  },
  "lockFileMaintenance": {
    "enabled": true,
    "schedule": ["before 4am on monday"]
  }
}
```

`minimumReleaseAge: "30 days"` is the rule §1.1 enforcement. The `vulnerabilityAlerts` override allows immediate upgrades when Renovate detects a published CVE.

If your project doesn't use Renovate, use Dependabot's `cooldown` option in `.github/dependabot.yml`:

```yaml
version: 2
updates:
  - package-ecosystem: "npm"
    directory: "/"
    schedule:
      interval: "weekly"
    cooldown:
      semver-major-days: 30
      semver-minor-days: 30
      semver-patch-days: 30
    open-pull-requests-limit: 10
```

### 3. Lockfile committed

**Tier 3 block.** Commits to `package.json` without a corresponding lockfile change are rejected.

| Package manager | Lockfile |
|---|---|
| pnpm | `pnpm-lock.yaml` |
| npm | `package-lock.json` |
| yarn | `yarn.lock` |
| bun | `bun.lockb` |

CI MUST install with frozen-lockfile:
- pnpm: `pnpm install --frozen-lockfile`
- npm: `npm ci`
- yarn: `yarn install --frozen-lockfile`
- bun: `bun install --frozen-lockfile`

A CI run that mutates the lockfile fails.

### 4. Node version pinned

Add `.nvmrc` or `.tool-versions` at the repo root:

```
# .nvmrc
20.11.1
```

or

```
# .tool-versions
nodejs 20.11.1
```

CI MUST use the pinned version. Reference it in workflow files:

```yaml
- uses: actions/setup-node@v4
  with:
    node-version-file: '.nvmrc'
```

### 5. CI vulnerability scanning

Copy [`templates/github-workflows/security.yml`](../templates/github-workflows/security.yml) into `.github/workflows/`.

It runs on every PR and:
- Verifies the lockfile is committed and frozen
- Runs `pnpm audit --prod` (or equivalent for the package manager in use)
- Fails the build on findings at severity HIGH or CRITICAL
- Logs MEDIUM/LOW findings for review

### 6. TypeScript: strict mode

For TypeScript projects, `tsconfig.json` MUST include:

```jsonc
{
  "compilerOptions": {
    "strict": true,
    "noUncheckedIndexedAccess": true,
    "noImplicitOverride": true,
    "noFallthroughCasesInSwitch": true,
    "noPropertyAccessFromIndexSignature": true,
    "exactOptionalPropertyTypes": true
  }
}
```

The full `strict: true` is the floor. Individual strictness flags above it are added per-project but never removed below `strict: true`.

### 7. Linter + formatter

- ESLint (or Biome) configured and run in CI
- Prettier (or Biome) configured and run in CI
- A pre-commit hook (husky / lefthook / git hooks) that runs the linter and formatter before commit
- `git commit --no-verify` is forbidden (per DEBROS.md §3.4)

### 8. Application-security libraries (DEBROS.md §10)

For projects with HTTP / web surface, the following are the canonical
library choices. Web-service patterns (middleware wiring, header policy,
file-upload workflow) live in [`web-service.md`](web-service.md); this
section covers JS/TS-specific package selection.

#### Schema validation — §10.4

- **`zod` ≥3.22.** TypeScript-first schemas, type inference end-to-end. Use `schema.parse()` in handlers that should throw on bad input; `schema.safeParse()` at boundaries where you want to shape the failure response yourself.
- Validate at the boundary (the HTTP handler), not deep in business logic. Internal functions accept the inferred type (`z.infer<typeof Schema>`) and trust it.

#### Password hashing — §10.3

- **`@node-rs/argon2`** (preferred — native binding, fast, modern argon2id) OR
- **`bcrypt`** (the npm package, NOT `bcryptjs` — the latter is pure-JS and ~30× slower, which hurts the user with no security benefit). Cost factor **≥12**.
- **Forbidden for passwords:** `crypto.createHash('md5'|'sha1'|'sha256'|'sha512')`, `node:crypto`'s `pbkdf2` with fewer than 600_000 iterations, any in-house hash function.
- bcrypt's 72-byte input limit: if you allow passwords longer than that, pre-hash with `crypto.createHash('sha256')` and bcrypt the hex digest. Note this in `debros.json.ai_agent_notes` so future contributors don't think the truncation is fine.

#### Security headers — §10.6

- **Express:** `helmet` ≥7.x with `contentSecurityPolicy` configured explicitly (the default is restrictive in ways that break inline tooling — set it yourself).
- **Fastify:** `@fastify/helmet`.
- **Next.js:** `next-safe` or a `next.config.js` `headers()` function.
- Always set `frame-ancestors` in CSP. Do not also set `X-Frame-Options` (see §10.6 — they conflict).

#### Rate limiting — §10.8

- **`express-rate-limit`** for single-node deployments (in-memory store).
- **`@upstash/ratelimit`** or **`rate-limiter-flexible`** with Redis for multi-node.
- Apply specifically to `/api/auth/*` and any endpoint that sends email/SMS.

#### File uploads — §10.7

- **`multer`** with `limits.fileSize = 10 * 1024 * 1024` (10 MB).
- **`file-type`** for content-sniffed MIME validation. Never trust `req.file.mimetype` — that's the client's claim, not the actual content.
- Generate the storage path server-side: `${uuid()}.${extFromSniffedMime}`.

#### Secret scanning — §10.5 + Tier-3 `hardcoded-secret`

- Wire **`gitleaks`** or **`trufflehog`** into CI on every PR.
- Husky/lefthook pre-commit hook with `gitleaks git --staged` catches secrets before they leave the dev's machine.

---

## File-by-file checklist

| File | Path | Required? | Tier-3 block? |
|---|---|---|---|
| `.npmrc` | repo root | ✅ | ✅ |
| `renovate.json` or `.github/dependabot.yml` | repo root or `.github/` | ✅ | — |
| Lockfile (`pnpm-lock.yaml` etc.) | repo root | ✅ | ✅ |
| `.nvmrc` or `.tool-versions` | repo root | ✅ | — |
| `.github/workflows/security.yml` | `.github/workflows/` | ✅ | — |
| `tsconfig.json` with `strict: true` | repo root (TS only) | ✅ | — |
| ESLint / Biome config | repo root | ✅ | — |
| Pre-commit hook config | repo root | ✅ | — |
| Zod schema for every HTTP handler boundary | source | ✅ if HTTP service | — |
| bcrypt / argon2 in use; no MD5/SHA-* on password values | source | ✅ if user accounts | ✅ (`insecure-password-hash`) |
| `helmet` / equivalent + explicit CSP | source | ✅ if HTTP service | — |
| Rate-limit middleware on `/api/auth/*` | source | ✅ if HTTP service | — |
| `multer` `limits` + `file-type` MIME sniff | source | ✅ if uploads | — |
| `gitleaks` / `trufflehog` in CI | `.github/workflows/` | ✅ | ✅ (`hardcoded-secret`) |

---

## Common patterns to enforce

### Package additions

When the agent or a human adds a new dependency, the agent MUST verify:

1. The package's most recent version was published ≥30 days ago (per rule §1.1) OR there's a Renovate `securityVulnerabilityAlerts` waiver
2. The package does not have install scripts, OR if it does, those scripts are reviewed and the package is explicitly allowlisted in `pnpm.onlyBuiltDependencies`
3. The package has more than one maintainer (single-maintainer packages with broad reach are a known supply-chain risk)
4. The package's `package.json` does not show signs of recent ownership transfer (check on npm registry — recent maintainer email change is a red flag)

The agent reports its findings on each of these before adding the dependency.

### `package.json` curation

Forbidden in `package.json`:
- `"dependencies": { ..., "*": "..." }` — never depend on `*` versions
- `"scripts": { "postinstall": "curl ... | sh" }` — never run remote shell scripts in lifecycle hooks
- `"resolutions"` / `"overrides"` without a tracked ticket explaining why

### Test framework

Use Vitest, Jest, or the platform's native test runner. The unit suite MUST run in <30 seconds (DEBROS.md §2.4). Tests with real network calls or `setTimeout`-based waits are forbidden — use fake timers and mock servers.

---

## Migration from a stock project

If you're adopting these rules in an existing project:

1. **Add `.npmrc` first.** This is the highest-value change. Expect some packages to fail to install — their install scripts were doing real work. Add those packages to `pnpm.onlyBuiltDependencies`.
2. **Audit existing dependencies.** Run `pnpm audit --prod` and resolve HIGH/CRITICAL findings. Run `npm ls --all` and look for single-maintainer packages with broad reach. Consider removing or replacing.
3. **Add `renovate.json`.** Renovate will start opening upgrade PRs respecting the 30-day cooldown. Review them; don't auto-merge.
4. **Add the CI security workflow.** Fix anything it catches.
5. **Update `debros.json`** to record that JS/TS compliance is satisfied.

Expect the first migration to take half a day. Subsequent maintenance is minimal.
