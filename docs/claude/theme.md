---
icon: lucide/palette
---

# Claude Code theme

`.claude/themes/cyberdream.json` is a complete `base: dark` theme override that maps every
visual role Claude Code exposes to a cyberdream color.

## Activation

In `.claude/settings.json`:

```json
"theme": "custom:cyberdream"
```

The `custom:` prefix tells Claude Code to look in `.claude/themes/<name>.json` rather than
loading a built-in theme.

## Structure

```json title=".claude/themes/cyberdream.json"
{
    "name": "Cyberdream",
    "base": "dark",
    "overrides": {
        "claude": "#FF5EA0",
        "claudeShimmer": "#FF5EF1",
        "text": "#FFFFFF",
        "inverseText": "#16181A",
        "inactive": "#3C4048",
        "subtle": "#3C4048",
        "suggestion": "#5EA1FF",
        "permission": "#BD5EFF",
        "remember": "#BD5EFF",

        "success": "#5EFF6C",
        "error": "#FF6E5E",
        "warning": "#F1FF5E",
        "merged": "#BD5EFF",

        "promptBorder": "#5EA1FF",
        "planMode": "#5EF1FF",
        "autoAccept": "#5EFF6C",
        "bashBorder": "#FF5EA0",
        "ide": "#5EA1FF",
        "fastMode": "#FFBD5E"
        // … plus diff colors, rate-limit bars, subagent rainbow, etc.
    }
}
```

## Notable mappings

| Role                            | Color                 | Meaning                                            |
| ------------------------------- | --------------------- | -------------------------------------------------- |
| `claude`                        | `#FF5EA0` pink        | Claude's own messages and avatar.                  |
| `permission`                    | `#BD5EFF` purple      | Permission prompts (e.g., "Allow Bash to run X?"). |
| `planMode`                      | `#5EF1FF` cyan        | Plan mode border / accents.                        |
| `autoAccept`                    | `#5EFF6C` green       | Auto-accept indicator.                             |
| `bashBorder`                    | `#FF5EA0` pink        | Bash invocation border.                            |
| `fastMode`                      | `#FFBD5E` orange      | Fast-mode badge.                                   |
| `success` / `error` / `warning` | green / red / yellow  | Standard status colors.                            |
| `diffAdded` / `diffRemoved`     | dark green / dark red | Diff backgrounds in code blocks.                   |

## Shimmer pairs

Several roles have a "shimmer" companion — a lighter variant used during animations
(thinking, streaming, etc.). Examples:

- `claude` / `claudeShimmer` — pink → magenta
- `permission` / `permissionShimmer` — purple → lavender
- `fastMode` / `fastModeShimmer` — orange → light orange

These keep the animated effects on-brand even when they pulse out to lighter shades.

## Subagent rainbow

`*_FOR_SUBAGENTS_ONLY` keys define the colors used to distinguish parallel subagents in the
UI. They map directly to the cyberdream palette's primary hues:

| Key                         | Color     |
| --------------------------- | --------- |
| `red_FOR_SUBAGENTS_ONLY`    | `#FF6E5E` |
| `blue_FOR_SUBAGENTS_ONLY`   | `#5EA1FF` |
| `green_FOR_SUBAGENTS_ONLY`  | `#5EFF6C` |
| `yellow_FOR_SUBAGENTS_ONLY` | `#F1FF5E` |
| `purple_FOR_SUBAGENTS_ONLY` | `#BD5EFF` |
| `orange_FOR_SUBAGENTS_ONLY` | `#FFBD5E` |
| `pink_FOR_SUBAGENTS_ONLY`   | `#FF5EA0` |
| `cyan_FOR_SUBAGENTS_ONLY`   | `#5EF1FF` |
