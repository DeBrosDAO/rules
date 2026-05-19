# Compliance — Web Service

> Concrete patterns to satisfy [DEBROS.md §10](../DEBROS.md) for HTTP services
> in any language. Applies to projects with `debros.json.project.type` in
> `{service, web, sdk}`. Language-specific implementation details live in
> the matching `compliance/<language>.md`.

---

## 1. Auth middleware (DEBROS.md §10.2)

**Pattern.** A single middleware reads `debros.json.project.public_routes[]` at startup, matches each incoming request path against the list, lets public routes through, and rejects everything else with `401` unless a valid auth token is present.

**Recommended starting set for `public_routes[]`** (copy into your `debros.json` and prune):

```json
"public_routes": [
  "/health",
  "/ready",
  "/api/auth/login",
  "/api/auth/register",
  "/api/auth/forgot-password",
  "/api/auth/reset-password",
  "/api/auth/verify-email",
  "/api/auth/oauth/callback"
]
```

**Per-framework sketches** (the real code goes through your auth library; these show the wiring):

- **Express / Hono / Fastify**: middleware registered globally, calls `next()` if the path matches the allowlist, otherwise verifies the auth token before delegating.
- **FastAPI**: a `Depends()` guard on the router, with explicit exemptions on public routes via `dependencies=[]`.
- **Gin**: a middleware on the global router that skips for allowlisted paths, runs auth otherwise.
- **Spring Boot**: a `SecurityFilterChain` bean with the allowlist in `authorizeHttpRequests`.

Adding a route to `public_routes[]` is a security-sensitive change. Sub-agent review (§4) applies.

---

## 2. Security headers (DEBROS.md §10.6)

| Language / framework | Library |
|---|---|
| JavaScript / TypeScript | `helmet` (≥7.x) — also `@fastify/helmet`, `next-safe` |
| Python / Django | `django-csp`, `django-permissions-policy` |
| Python / FastAPI | `secure` |
| Ruby / Rails | `secure_headers` |
| Go | `unrolled/secure` |

**Baseline CSP**, copy-pasteable starting point — tighten per-app:

```
default-src 'self';
frame-ancestors 'none';
script-src 'self';
style-src 'self' 'unsafe-inline';
img-src 'self' data:;
object-src 'none';
base-uri 'self';
form-action 'self';
```

**Other required headers**:

```
Strict-Transport-Security: max-age=31536000; includeSubDomains; preload
X-Content-Type-Options: nosniff
Referrer-Policy: strict-origin-when-cross-origin
```

Do **not** set `X-Frame-Options` if you have CSP `frame-ancestors` — pick one. CSP wins.

---

## 3. Rate limiting (DEBROS.md §10.8)

| Stack | Library |
|---|---|
| JavaScript / TypeScript | `express-rate-limit` (in-memory), `@upstash/ratelimit` (Redis) |
| Python / FastAPI | `slowapi` |
| Python / Django | `django-ratelimit` |
| Go | `golang.org/x/time/rate`, `didip/tollbooth` |
| Ruby / Rails | `rack-attack` |

**Storage.** In-memory is fine for single-node deployments. Multi-node deployments need Redis (or equivalent) — otherwise per-IP limits are per-node, and an attacker fans out across nodes.

**Required behavior.** Reject excess requests with `429 Too Many Requests`. Logging-only mode is not compliant with §10.8 — the limit must actually reject.

**Apply to**: `/api/auth/login`, `/api/auth/register`, `/api/auth/forgot-password`, `/api/auth/reset-password`, `/api/auth/verify-email` (token-guessing target), `/api/auth/oauth/callback` (public exchange route), and any endpoint that triggers an email or SMS send. Any route in `public_routes[]` that accepts a token, code, or secret as input MUST also be rate-limited — being publicly reachable does not exempt it.

---

## 4. File uploads (DEBROS.md §10.7)

| Stack | Upload + MIME-sniff |
|---|---|
| JavaScript / TypeScript | `multer` + `file-type` |
| Python | `python-multipart` + `python-magic` (or `filetype`) |
| Go | `http.DetectContentType` (stdlib) + `mime/multipart` |
| Ruby / Rails | `marcel` + `active_storage` |

**Workflow**:

1. Enforce the size limit at the parser level (`multer({ limits: { fileSize: 10 * 1024 * 1024 } })` etc.). Don't read 100 MB into memory just to reject it.
2. Sniff the actual content for MIME — **do not** trust the client's `Content-Type` header.
3. Compare against an allowlist (e.g. `['image/png', 'image/jpeg', 'application/pdf']`). Allowlist beats denylist.
4. Generate the storage filename server-side (`<uuid>.<sniffed-extension>`). Discard the client name.
5. Store outside the web root, or behind a controller that serves with `Content-Disposition: attachment` and the right MIME.

---

## 5. Input validation at boundaries (DEBROS.md §10.4)

Language-specific tooling lives in the matching `compliance/<language>.md`. This file documents **where** validation belongs:

- HTTP request body, query string, path parameters
- Selected headers — Authorization, Content-Type (not every header)
- Webhook payloads (validate even when the signature is verified — they're separate concerns)
- File uploads (alongside §10.7's MIME check)
- External API responses your code parses

Validate **once at the boundary**. Internal code trusts the validated types and does not re-validate.

---

## 6. Common patterns to enforce / forbid

These follow from §10 but are common-enough mistakes to call out explicitly.

- **HttpOnly + Secure + SameSite cookies for auth.** Never store auth tokens in `localStorage` (XSS-readable). Cookies with `HttpOnly` + `Secure` + `SameSite=Lax` (or `Strict` where the UX allows) are the standard pattern.
- **No stack traces in production responses.** Error responses include a short message and a correlation ID. Stack traces and internal error details go to logs only.
- **Generic auth-failure messages.** `401` for "invalid credentials" — not "user not found" vs "wrong password." The latter is a user-enumeration oracle.
- **CSRF protection on state-changing routes** that authenticate via cookies. Routes that authenticate via `Authorization: Bearer` and don't accept cookies are not CSRF-vulnerable; routes that accept cookies are.
- **Constant-time string comparison for tokens and signatures.** `===` / `==` short-circuits on the first byte and leaks timing.

---

## 7. Tier-3 blocks that affect web services

The agent refuses to proceed past these unless overridden in `debros.json.compliance.tier3_overrides[]`:

| Block | Triggered by | Rule |
|---|---|---|
| `hardcoded-secret` | A literal-looking secret in source | §10.5 |
| `insecure-password-hash` | MD5/SHA-* used for password hashing | §10.3 |
| `sql-string-concat` | String concat/interpolation building SQL | §10.1 |
| `missing-auth-on-route` | Route handler, no auth, not in `public_routes[]` | §10.2 |

These are the rule statements; a machine-readable equivalent for CI integration is provided by the `templates/tier3.json` tooling patch (not required for the prose rules to apply).

---

## 8. File-by-file checklist

| Check | Required? |
|---|---|
| Auth middleware exists and is wired into the app entry point | ✅ |
| `debros.json.project.public_routes[]` populated | ✅ |
| Security-headers middleware present (helmet / secure / etc.) | ✅ |
| Rate-limit middleware on auth routes, rejecting (not logging) | ✅ |
| File-upload handler validates MIME via content-sniffing + 10 MB cap | ✅ if uploads exist |
| Per-language input-validator chosen and used at every boundary | ✅ |
| Auth cookies use HttpOnly + Secure + SameSite | ✅ if cookie auth |
| Generic 401 message on auth failure (no user enumeration) | ✅ |

---

## Migration from a stock web service

1. **Populate `public_routes[]`** in your `debros.json` from the recommended starting set, then prune.
2. **Add auth middleware** if not already present, with the public-route allowlist wired in.
3. **Install a security-headers library** with the baseline CSP from §2.
4. **Add rate limiting** to `/api/auth/*` routes. In-memory is fine to start; switch to Redis on multi-node.
5. **Audit file-upload routes** for MIME sniffing + size caps + server-generated filenames.
6. **Audit password storage** — migrate any MD5/SHA-hashed passwords to bcrypt/argon2id on next login (rehash transparently when the user authenticates).
7. **Run a secret scanner** (gitleaks / trufflehog) against git history. Rotate anything found.

Expect a few days for an existing service that hasn't followed these patterns. Most of the cost is in the password-storage migration and the secret-rotation sweep.
