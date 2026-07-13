---
name: rules-orama
description: How to use the Orama Network — deploy apps, SQLite databases, custom domains, and serverless WASM functions, run releases and rolling upgrades, and use the client SDK. Use whenever a task involves Orama: deploying to it, writing Orama functions, calling an Orama gateway, or operating Orama nodes.
---

# Using Orama Network

Orama is a decentralized platform for deploying web apps, SQLite databases,
custom domains, and serverless WASM functions across a peer-to-peer node
network, reached through one API gateway per namespace.

## Source of truth: fetch the docs, don't guess

The current, code-verified docs are served as raw markdown. **Fetch the index
first, then the specific doc you need** — do not rely on memory or on this file
for commands, flags, or endpoints:

```
https://iofo4ifs.orama.network/llms.txt
```

That index lists each doc with a one-line description and a direct URL. Fetch
the one matching your task:

| Task | Doc |
|---|---|
| Deploy an app (static / Next.js / Go / Node) | `llms/deploying-apps.md` |
| Write / deploy / invoke a function | `llms/functions.md` |
| Cut a release or roll out a cluster upgrade | `llms/release-and-rollout.md` |
| Understand the system | `llms/architecture.md` |
| Call a gateway from Go | `llms/client-sdk.md` |
| Check cluster / node health | `llms/monitoring.md` |
| Diagnose a failure | `llms/troubleshooting.md` |

## The one gotcha worth stating up front

`orama deploy <type> <dir>` takes a **source directory** and builds + tarballs
it for you (with `CGO_ENABLED=0`). Do **not** pre-build a tarball or a binary
and pass that — it fails. Everything else: confirm against the endpoint above.

## Operating rules (when working ON an Orama cluster)

- Rolling upgrades only — never restart multiple RQLite voters at once (quorum).
- Drive nodes through the `orama` CLI (`orama node …`), never raw `systemctl`.
- Inter-node traffic is on the WireGuard overlay (`10.0.0.x`), not public IPs.

If the endpoint is unreachable, the same docs live in the Orama repo under
`docs/` (root).
