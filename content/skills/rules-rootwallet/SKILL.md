---
name: rules-rootwallet
description: How to use RootWallet and how it integrates with the Orama Network — the local agent protocol for requesting SSH keys/passwords/signatures (used by the Orama CLI to deploy), the registry functions RootWallet runs on Orama, cloud sync, and the crypto/security model. Use when a task touches RootWallet, wallet signing, the agent socket, or the RootWallet↔Orama boundary.
---

# Using RootWallet

RootWallet is a self-custodial multi-chain wallet and secrets vault. It runs
registry functions on the Orama Network, and it exposes a **local agent over a
Unix socket** so tools (notably the Orama CLI) can request SSH keys, passwords,
and signatures without ever seeing the underlying secrets.

## Source of truth: fetch the docs, don't guess

Current, code-verified docs are served as raw markdown. **Fetch the index
first, then the specific doc** — don't rely on memory or this file for wire
formats, capabilities, or crypto parameters:

```
https://rootwallet.io/llms.txt
```

| Task | Doc |
|---|---|
| Have an app/CLI request a key or signature from RootWallet | `llms/agent-protocol.md` |
| Work on the functions RootWallet runs on Orama | `llms/orama-functions.md` |
| Understand vault backup / recovery | `llms/cloud-sync.md` |
| Understand the system | `llms/architecture.md` |
| Key hierarchy / derivation / keystore | `llms/crypto-architecture.md` |
| Threat model & implemented crypto | `llms/security.md` |

## How RootWallet and Orama fit together

- **Orama → RootWallet:** the Orama CLI talks to RootWallet's agent socket to
  get SSH keys / passwords / signatures for deploys. See `agent-protocol.md`.
- **RootWallet → Orama:** RootWallet deploys registry functions to its own
  Orama namespace and stores data via the namespace gateway. See
  `orama-functions.md`. For deeper Orama detail, also load the `rules-orama`
  skill.

## Gotchas worth stating up front (confirm specifics at the endpoint)

- Agent capability grants (including `wallet:sign`) **persist to disk** across
  locks and restarts — they are not re-prompted per use.
- `vault:ssh` / `vault:password` return **stored** vault entries; they do not
  derive keys.
- Cloud sync commands (`rw sync …`) are **not yet available** (stubbed).

If the endpoint is unreachable, the same docs live in the RootWallet repo under
`docs/` (root).
