---
icon: lucide/plug
---

# Plugins

OpenCode auto-loads every `plugin/*.{js,ts}` file under `~/.config/opencode/`, so dropping a file into
`.config/opencode/plugin/` is the whole wiring — there is no entry to add to `opencode.jsonc`. This repo ships one
plugin.

## Status indicator { #status-indicator }

`.config/opencode/plugin/agent-tmux-status.js` flags the tmux window (or terminal title) when OpenCode is waiting on
you, reusing the shared [`agent-tmux-status`](../scripts/tmux.md#agent-tmux-status) leaf script that Claude Code
drives from [its hooks](../claude/hooks-skills.md#claude-is-waiting-indicator). One indicator, every agent: whether
the pane is running OpenCode, Claude Code, or [Codex](../codex/hooks.md#the-tmux-indicator), the tmux window lights
up the same way.

The plugin subscribes to OpenCode's [event bus](https://opencode.ai/docs/plugins/) and maps each event to a state
token the script hands to [`theme.conf`](../terminal/tmux.md#agent-status):

| OpenCode event                | State       | Look                           |
| ----------------------------- | ----------- | ------------------------------ |
| `session.idle`                | `waiting`   | calm **peach** background, `●` |
| `permission.updated`          | `attention` | bold **red** background, `󰂚`   |
| `message.updated` (user role) | `clear`     | indicator removed              |

The `permission.updated` event is the louder of the two because an approval prompt actively needs you, mirroring how
the Claude Code [`Notification` hook](../claude/hooks-skills.md#claude-is-waiting-indicator) escalates over the calm
`session.idle` → `waiting` state.

It shells out via the plugin's Bun `$` helper with `.quiet().nothrow()` and wraps the call in a `try`/`catch`, so a
status blip can never disrupt a session. The pane is targeted by `$TMUX_PANE`, inherited from the shell OpenCode
launched in.

## Plugin dependency

The plugin imports [`@opencode-ai/plugin`](https://opencode.ai/docs/plugins/), pinned in
`.config/opencode/package.json`. OpenCode installs it locally; the resulting `node_modules`, lockfiles, and
generated `package.json` are gitignored (see [Files](index.md#files)), so only the hand-written plugin source is
tracked.

## Adding more

Additional plugins belong on this page. Drop the `plugin/*.js` file in, add a subsection here describing the events
it listens for, and — if it pulls a new import — confirm the install artifacts stay covered by
`.config/opencode/.gitignore`.
