---
icon: lucide/terminal
---

# The `dot` CLI

`dot` is a multi-call Go binary that backs the more complex helper commands. The logic-heavy scripts in
`.local/share/scripts/` are migrating into it so their behaviour is unit-tested instead of only verifiable by
running them against live `git`/`brew` state.

## Layout

Standard Go layout, module rooted at the repository:

| Path                                                          | What                                                                                               |
| ------------------------------------------------------------- | -------------------------------------------------------------------------------------------------- |
| `go.mod` / `go.sum`                                           | Module `github.com/dmccaffery/dotfiles`; direct deps: Cobra, charmbracelet/log, charmbracelet/huh. |
| `cmd/dot/main.go`                                             | Entry point and `argv[0]` dispatch.                                                                |
| `internal/cmd/root`                                           | Assembles the root command from the per-command packages; owns the applet registry.                |
| `internal/cmd/<command>`                                      | One package per command (`worktree`, `brewfile`, `zs`, …), each exposing `NewCmd`.                 |
| `internal/cmd/cmdutil`                                        | Shared command framework: the injectable `Deps`, `ErrSilent`, and small helpers.                   |
| `internal/cmd/cmdtest`                                        | Test-only helpers that drive a command and build fake `Deps`.                                      |
| `internal/worktree`, `internal/brewfile`, `internal/colorbar` | Pure, table-tested logic (sanitizer, hook JSON, `trust.json`, gradient).                           |
| `internal/execx`                                              | A mockable `Runner` over `os/exec` — the unit-test seam for git/tmux/brew calls.                   |
| `internal/logx`                                               | slog diagnostics via charmbracelet/log: styled levels on a TTY, JSON when piped.                   |
| `internal/ui`                                                 | A testable `Prompter` (confirmations/questions) backed by charmbracelet/huh.                       |
| `internal/envx`                                               | Environment and XDG lookups.                                                                       |

## Multi-call dispatch

One binary serves every command; it dispatches on `filepath.Base(os.Args[0])`. Each command is also exposed as a
**symlink** of the same name pointing at the binary, so existing call sites keep working unchanged:

| You run                                 | Resolves to                                |
| --------------------------------------- | ------------------------------------------ |
| `dot worktree start`                    | the worktree command                       |
| `worktree start` (symlink)              | same — `argv[0]` is `worktree`             |
| `~/.local/share/scripts/worktree start` | the Claude Code `WorktreeCreate` hook path |

`--help` works at every level (`dot --help`, `dot worktree --help`, `worktree --help`). The hidden `dot applets`
command lists the names that should be symlinked — the single source of truth the build stage reads.

## Build and install

[`setup/build.sh`](../../setup/build.sh) is the `build` install stage (it runs between `stow` and `packages`):

```sh
make build          # = ./install.sh build
```

It compiles `~/.local/bin/dot` (a per-machine artifact, never committed), self-checks that it runs, then links one
symlink per `dot applets` entry — plus a `dot` self-link — into `~/.local/share/scripts/` (already on `PATH`). `go`
is installed by the `requirements` stage via [`Brewfile.requirements`](../../setup/darwin/Brewfile.requirements).

!!! warning "Rebuild after pulling"

    The binary is compiled, not a live script — run `make build` after a `git pull` that touched `cmd/` or
    `internal/` so the installed `dot` matches the source. `dot --version` embeds `git describe` to help spot drift.

## One command, one home

A command is **either** a shell script in `.local/share/scripts/` (stowed) **or** a `dot` applet (a symlink to the
binary), never both. Porting one deletes its shell file in the same change that registers the Go command, so `stow`
and the build stage never fight over the same path.

## Logging and prompts

Diagnostics go through [`internal/logx`](../../internal/logx), which adapts `charmbracelet/log` to slog: styled,
leveled lines on a terminal and structured JSON when output is piped, with typed attributes. Command **results**
(such as the worktree path) are written to stdout only, so callers and hooks read them cleanly.

Interactive confirmations and questions go through [`internal/ui`](../../internal/ui)'s `Prompter`, backed by
`charmbracelet/huh` on `/dev/tty` (so prompts reach the terminal even when stdout is captured). With no tty the
prompt is skipped via `ErrNoTTY` and the caller falls back — e.g. brewfile leaves a tap untrusted with a warning.
`Prompter` is an interface, so commands unit-test their prompt flows against a fake.

## Testing

Three layers, all under `go test ./...`:

- **Pure logic** — table tests for the name sanitizer, hook JSON and `trust.json` parsing (no I/O).
- **Control flow** — commands run against a fake `execx.Runner` and a fake `Prompter`, asserting the right
  git/brew calls and prompt decisions without touching the system.
- **Real git** — integration tests create a temp repo and prove `worktree start`/`end` actually create and remove
  the worktree directory and its `agent/*` branch (they skip when `git` is absent; CI has it).

## What's ported

Each command's external interface (name, arguments, output) is unchanged — only the implementation moved into
`dot`.

- **Wave 1** — [`worktree`](../scripts/tmux.md), [`agent-tmux-status`](../scripts/tmux.md#agent-tmux-status),
  [`brewfile`](../scripts/misc.md#brewfile)
- **Wave 2** — [`print-colors`](../scripts/misc.md#print-colors),
  [`profile-shell`](../scripts/misc.md#profile-shell), [`git-resign`](../scripts/security-keys.md#git-resign),
  [`gh-switch-user`](../scripts/security-keys.md#gh-switch-user),
  [`fzf-image-preview`](../scripts/misc.md#fzf-image-preview),
  [`reset-background-items`](../scripts/misc.md#reset-background-items), [`zs`](../scripts/zscaler.md)

Still shell, ported later: `git-github-auth` and `tmux-session`, then the security-key flows `ssh-sk` and
`ssh-askpass` **last** — so commit signing and ssh-agent login never come to depend on a successful build.

## CI

[`pull-request-go.yaml`](../../.github/workflows/pull-request-go.yaml) runs `go build`, `go vet`, `go test` and a
`gofmt` check on any change under `cmd/`, `internal/`, or `go.mod`. Dependabot tracks the `gomod` ecosystem
alongside npm, uv and GitHub Actions.
