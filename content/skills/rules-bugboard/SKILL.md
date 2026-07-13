---
name: rules-bugboard
description: How to use bugboard.ai — the task/bug tracker with an MCP server where the coding agent is a teammate. Covers connecting a project to the bugboard MCP, the task tools (create/list/update/comment/link tasks and initiatives), and where to file work. Use when filing bugs/features/tasks, querying a project's tasks, or wiring the bugboard MCP into a repo.
---

# Using bugboard.ai

bugboard is a task and bug tracker with a first-class **MCP server** — you
interact with it through MCP tools, not a REST API you hand-roll. Use it to
file bugs/features, list and update tasks, comment, and link tasks into
initiatives.

## Connecting a project to the bugboard MCP

Add the server to the repo's MCP config (`.mcp.json` for Claude Code). The
token is per-user, per-org and starts with `bb_`; keep it out of committed
files — reference an env var:

```json
{
  "mcpServers": {
    "bugboard": {
      "type": "http",
      "url": "https://api.bugboard.ai/mcp/sse",
      "headers": { "Authorization": "Bearer ${BUGBOARD_MCP_TOKEN}" }
    }
  }
}
```

Create the token in the bugboard web app (Settings → API tokens). Every tool
call is authorized to that user's org and audit-logged.

## The tools you'll actually use

Prefer the unified **task** surface (a task's `type` is bug / feature / chore):

- `list_projects`, `list_tasks`, `search_tasks`, `get_task`
- `create_task`, `create_tasks` (batch), `update_task`, `update_tasks`,
  `comment_task`, `set_task_parent`, `change_task_type`
- `check_duplicates` before filing, so you don't create a dup
- Initiatives & links: `create_initiative`, `link_task_to_initiative`,
  `link_tasks` (blocks / relates_to / duplicate_of / chain)

(Legacy `*_bug` / `*_feature` tools exist but are hidden — use the task tools.)

## When to file

When you spot out-of-scope work — a bug, a follow-up, tech debt, a security
issue noticed in passing — file it as a bugboard task rather than dropping it.
Check `check_duplicates` first.

## Full docs

Product overview, architecture, and MCP details are served as raw markdown at:

```
https://bugboard.ai/llms.txt
```

Fetch that index for the current tool list and connection details before
relying on anything above.
