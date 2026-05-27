---
icon: lucide/chevron-right
---

# oh-my-posh

The prompt is [oh-my-posh](https://ohmyposh.dev). Two configs ship in this repo:

| File                             | Used by                                           |
| -------------------------------- | ------------------------------------------------- |
| `.config/oh-my-posh/prompt.yaml` | Main shell prompt (initialized in `.zshrc`).      |
| `.config/oh-my-posh/claude.yaml` | [Claude Code](../claude/settings.md) status line. |

## Main prompt

The prompt is a two-line, multi-segment design with a cyberdream palette baked in. Top line
shows OS, path, and Git status (left) plus language version segments (right). Bottom line is a
minimal `╰─` continuation with an exit-code-colored arrow.

### Layout

```text
╭─  user   ~/Repos/dotfiles    main ✔                              Lua 5.4
╰─ ❯
```

### Palette

The `palette:` block at the top of `prompt.yaml` defines every color used in the prompt. The
cyberdream-anchored values:

```yaml title=".config/oh-my-posh/prompt.yaml"
palette:
    base: "#16181A" # background
    text: "#FFFFFF"
    blue: "#5EA1FF"
    mauve: "#FF5EA0"
    green: "#5EFF6C"
    red: "#FF6E5E"
    yellow: "#F1FF5E"
    peach: "#FFBD5E"
    teal: "#5EF1FF"
    lavender: "#BD5EFF"
    flamingo: "#FFACF6"
    pink: "#FF5EF1"
    sapphire: "#007CDE"
    # ...
```

These names are referenced from segments as `p:blue`, `p:mauve`, etc.

### Segments shipped

| Side   | Segment  | Notes                                                       |
| ------ | -------- | ----------------------------------------------------------- |
| Left   | `os`     | OS icon in a leading diamond.                               |
| Left   | `path`   | Folder-style path; home shown as a house icon.              |
| Left   | `git`    | Branch + upstream + working/staging counts; cached `none`.  |
| Right  | `go`     | Triggered on `go.mod` / `go.work`.                          |
| Right  | `lua`    | Always shown (matches the dotfiles' nvim/tmux Lua use).     |
| Right  | `python` | Includes venv name when active.                             |
| Right  | `node`   | Triggered on `package.json`; shows package-manager icon.    |
| Right  | `gcp`    | Active `gcloud` project.                                    |
| Bottom | `text`   | `╰─` continuation.                                          |
| Bottom | `status` | The arrow, recoloured red when last exit code was non-zero. |

### Transient and secondary prompts

```yaml
secondary_prompt:
    template: "❯❯ "
    foreground: p:red

transient_prompt:
    template: "❯ "
    foreground_templates:
        - "{{if gt .Code 0}}p:red{{end}}"
        - "{{if eq .Code 0}}p:subtext_1{{end}}"
```

The transient prompt collapses past prompts down to a single `❯` after each command, keeping
scrollback readable. Color reflects the exit code of the _previous_ command.

## Claude Code status line

`.config/oh-my-posh/claude.yaml` powers the status line at the bottom of the Claude Code TUI.
It's wired up via `.claude/settings.json`:

```json
"statusLine": {
  "type": "command",
  "command": "oh-my-posh claude --config ${XDG_CONFIG_HOME}/oh-my-posh/claude.yaml"
}
```

`oh-my-posh claude` is a built-in subcommand that knows how to read Claude Code's session JSON
from stdin.

### What it shows

Left to right, the status line renders the OS icon, the current path, a `claude` segment, and the
git segment. The `claude` segment packs four pieces of information:

| Field                                | Source field                        | Shown as                         |
| ------------------------------------ | ----------------------------------- | -------------------------------- |
| Current model                        | `.Model.DisplayName`                | e.g. `Claude Opus 4.7`           |
| Reasoning effort (when supported)    | `.Effort.Level`                     | e.g. `· high` (omitted if empty) |
| Weekly (7-day rolling) usage         | `.SevenDayGauge` + `.SevenDayUsage` | gauge `▰▰▱▱▱` + `42%`            |
| Current session context-window usage | `.TokenUsagePercent`                | `18%`                            |

The 7-day window aggregates usage across **all models** and **all sessions** on the account, so
it's the limit most likely to surprise you during heavy use — the gauge gives an at-a-glance read
without needing to remember the exact percentage. The effort segment is conditional: models that
don't expose a reasoning-effort level (`.Effort.Level` is empty) skip the `· level` suffix entirely.
