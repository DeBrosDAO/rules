# Compliance — Python

> The concrete files every Python project must have to satisfy [DEBROS.md](../DEBROS.md).

Python's supply-chain story is meaningfully weaker than Go's: PyPI is rife with typosquatting, source distributions (sdists) can run arbitrary code at install time, and historic tooling (pip alone, plain `requirements.txt`) doesn't verify integrity hashes. The defenses below close those gaps.

---

## Recommended toolchain

We recommend **[uv](https://github.com/astral-sh/uv)** for new projects. Reasons:

- Native lockfile with embedded hashes (cryptographic integrity, like Go's `go.sum`)
- Built-in dependency-resolution + virtualenv management
- ~10× faster than pip/Poetry; running it in CI is essentially free
- Backed by an active well-funded team (Astral)
- Supports PEP 621 (`pyproject.toml`) — no proprietary config format

**Poetry** is acceptable for existing projects already using it. Plain `pip + requirements.txt` is acceptable ONLY if used with `pip-tools` + `--require-hashes`.

This document covers all three setups; commands assume uv unless noted.

---

## Required files

### 1. `pyproject.toml` + lockfile

**Tier 3 block.** Commits to `pyproject.toml` without a corresponding lockfile change are rejected.

| Toolchain | Lockfile | CI install command |
|---|---|---|
| uv | `uv.lock` | `uv sync --frozen` |
| Poetry | `poetry.lock` | `poetry install --no-root --sync` |
| pip-tools | `requirements.txt` with `--hash` entries | `pip install --require-hashes -r requirements.txt` |
| Pipenv (legacy) | `Pipfile.lock` | `pipenv install --deploy` |

A `pyproject.toml` MUST declare:

```toml
[project]
name = "your-project"
requires-python = ">=3.12"     # pin a major version
version = "0.1.0"

[project.optional-dependencies]
dev = [
  "pytest>=7",
  "ruff>=0.4",
  "mypy>=1.10",
  "pip-audit>=2.7",
]
```

`requires-python` enforces that the project won't install on a Python version it wasn't tested against.

### 2. `.python-version` — pin the runtime

```
3.12.4
```

Used by `uv`, `pyenv`, and `mise`. CI must use the pinned version.

### 3. Block source-distribution installs by default

**Tier 3 block** for production deploys. Without this, packages can run arbitrary code in `setup.py` at install time.

For **uv**, add to `pyproject.toml`:

```toml
[tool.uv]
# Refuse source distributions. Forces wheels (precompiled binaries),
# which cannot execute code at install. Packages that genuinely need
# to build from source must be explicitly allowlisted below.
no-build = true

# Allowlist for packages that have no wheel available. Each entry needs
# a one-line justification and the entry should be reviewed when changed.
no-build-package = []
```

For **pip**, in CI:

```bash
pip install --only-binary :all: -r requirements.txt
```

When a dependency genuinely requires building from source (rare for popular packages), add an explicit `--only-binary :all: --no-binary <package>` override AND review the package's `setup.py` for any code that runs at install time.

### 4. `renovate.json` with 30-day cooldown for Python

Renovate supports Python via `pep621`, `poetry`, and `pip_requirements` managers. Copy [`templates/renovate.json`](../templates/renovate.json) — the same file works across ecosystems.

Key config:

```jsonc
{
  "pep621": { "enabled": true },
  "poetry": { "enabled": true },
  "pip_requirements": { "enabled": true },
  "minimumReleaseAge": "30 days",
  "automerge": false
}
```

`minimumReleaseAge` is the rule §1.1 enforcement and covers typosquatting + compromised-account attacks.

### 5. `pip-audit` in CI

`pip-audit` is maintained by PyPA and queries the [PyPI Advisory Database](https://github.com/pypa/advisory-database).

Add to `.github/workflows/security.yml`:

```yaml
- name: Install pip-audit
  run: pip install pip-audit

- name: Run pip-audit
  run: |
    if [ -f uv.lock ]; then
      uv export --no-hashes | pip-audit --requirement /dev/stdin --strict
    elif [ -f poetry.lock ]; then
      pip-audit --strict
    elif [ -f requirements.txt ]; then
      pip-audit --requirement requirements.txt --strict
    fi
```

Findings at severity HIGH or higher fail the build.

### 6. Type checking

mypy (or pyright) in **strict mode**, run on every PR:

```toml
[tool.mypy]
strict = true
warn_return_any = true
warn_unused_ignores = true
disallow_any_explicit = true
disallow_untyped_decorators = true
no_implicit_optional = true
```

Untyped code is forbidden in new files. Existing untyped files may be exempted via a tracked exception in `debros.json.compliance.exceptions[]` with a migration plan.

### 7. Linter + formatter

**[Ruff](https://github.com/astral-sh/ruff)** for both. Single binary, ~100× faster than flake8/pylint, includes formatter (replaces black + isort).

```toml
[tool.ruff]
target-version = "py312"
line-length = 100

[tool.ruff.lint]
select = ["E", "F", "W", "I", "B", "C4", "UP", "N", "S", "ARG", "SIM"]
ignore = ["E501"]  # line length handled by formatter

[tool.ruff.format]
quote-style = "double"
```

The `S` (bandit) rules catch common security antipatterns. Run in CI; `git commit --no-verify` is forbidden.

---

## File-by-file checklist

| File | Path | Required? | Tier-3 block? |
|---|---|---|---|
| `pyproject.toml` | repo root | ✅ | — |
| Lockfile (`uv.lock` / `poetry.lock` / `requirements.txt` with hashes) | repo root | ✅ | ✅ |
| `.python-version` | repo root | ✅ | — |
| `renovate.json` | repo root | ✅ | — |
| `.github/workflows/security.yml` running `pip-audit` | `.github/workflows/` | ✅ | — |
| `[tool.uv] no-build = true` (or equivalent) | in `pyproject.toml` | ✅ | ✅ (production only) |
| `[tool.mypy] strict = true` | in `pyproject.toml` | ✅ | — |
| `[tool.ruff]` config | in `pyproject.toml` | ✅ | — |

---

## Code patterns to enforce

### Type annotations everywhere

Per the strict-mypy setting: every function signature is typed, every variable that crosses a module boundary is typed, `Any` is forbidden except behind a comment explaining why.

```python
# Good
def get_user(user_id: UserId) -> User | None: ...

# Bad
def get_user(user_id): ...
```

### Use `NewType` to make illegal states unrepresentable

Per DEBROS.md §2.2 principle 3:

```python
from typing import NewType

UserId = NewType("UserId", str)
SessionId = NewType("SessionId", str)

# Now this is a type error — caught at lint time:
def lookup(sid: SessionId) -> User: ...
lookup(user_id)  # mypy: error
```

### Errors carry context

Per DEBROS.md §2.2 principle 6:

```python
# Good
try:
    conn = connect(host, port)
except OSError as e:
    raise RuntimeError(f"connect to olric at {host}:{port}: {e}") from e

# Bad — loses context
except OSError:
    raise RuntimeError("connection failed")

# Forbidden — swallows the error
except OSError:
    pass
```

### Function and file sizes

Per DEBROS.md §2.1:
- Functions ≤50 lines
- Files ≤300 lines

Ruff's `PLR0915` rule + `complexity-too-much` enforces complexity. Add to `[tool.ruff.lint.per-file-ignores]` only for tests where setup naturally lengthens.

### Logging, not printing

`print()` in production code is forbidden. Use `logging` (or a structured logger like `structlog`).

```python
# Good
logger.info("user registered", extra={"user_id": user_id, "namespace": ns})

# Bad
print(f"user {user_id} registered")
```

### Async vs sync

Per DEBROS.md §2.2 principle 8: no premature concurrency. Default to sync. Move to async only when the project is genuinely I/O-bound and benchmarks justify it. Mixed sync/async codebases are a known footgun — pick a model per project and stick to it.

---

## Dependency additions

When adding a Python package, the agent MUST verify:

1. **Cooldown.** The package version was published ≥30 days ago (rule §1.1) OR there's a Renovate `vulnerabilityAlerts` waiver.
2. **Typosquat check.** The package name is spelled correctly. Common typosquats include misspelled popular packages (`reqests`, `urlib3`). Compare against the canonical name on PyPI; check download counts (low downloads on a "popular-named" package is a red flag).
3. **Maintainer count.** Single-maintainer packages with broad reach are higher risk. Check `Maintainers` on the PyPI project page.
4. **Wheels available.** Prefer packages that ship wheels. Source-only packages run `setup.py` at install time and are a higher attack surface.
5. **License compatible.** Check `License` on PyPI. The project's license MUST be compatible.

The agent reports its findings on each of these BEFORE adding the dependency.

---

## Common antipatterns

Forbidden in `pyproject.toml` / `requirements.txt`:

- Unversioned dependencies (`pandas` without a version constraint)
- Wildcard versions (`pandas = "*"`)
- VCS dependencies in production (`pandas @ git+https://...`) — these bypass PyPI and the lockfile's hash check; use only in dev with an explicit exception
- `pip install -e .` in production — editable installs bypass the lockfile
- Lifecycle hooks in `setup.py` (`cmdclass = {"install": MyCustomInstall}`) — extremely suspicious; review carefully

---

## Migration from a stock Python project

1. **Decide on a toolchain.** If starting fresh, use `uv`. If you have `requirements.txt`, run `uv init` and import; if Poetry, keep Poetry.
2. **Lock dependencies with hashes.** `uv sync` or `poetry lock` or `pip-compile --generate-hashes`.
3. **Pin Python version.** Add `.python-version`.
4. **Add Ruff + mypy strict config to `pyproject.toml`.** Fix what they catch — usually a few hours for a mid-size project.
5. **Add the CI workflow** with `pip-audit`.
6. **Add `renovate.json`.**
7. **Update `debros.json`** to record Python compliance is satisfied.

Expect first migration to take a day for projects with significant untyped code; subsequent maintenance is minimal.

---

## Notes on Python's supply-chain story

What Python protects against (when configured correctly):
- Hash-pinned lockfiles catch swapped packages (`uv.lock`, `poetry.lock`, `--require-hashes` mode)
- Wheel-only installs (`--only-binary :all:`) prevent install-time code execution
- `pip-audit` checks against the official PyPI advisory database

What Python does NOT protect against (handle these out of band):
- **Typosquatting.** Massive ongoing problem on PyPI. Defense is human review of every new package name.
- **Compromised maintainer accounts.** Defense is 30-day cooldown + watching for unusual recent releases.
- **`setup.py` arbitrary code.** Defense is `--only-binary :all:` and reviewing source-required packages.
- **Dependency confusion** (private package name shadowed by a public PyPI package). Defense is private index priority configuration + name reservation on PyPI.

Python supply-chain security is improving (PEP 458 / PEP 480 for repository signing, sigstore for provenance) but most of those defenses are not yet default. Until they are, the rules above are the floor.
