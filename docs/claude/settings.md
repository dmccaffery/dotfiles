---
icon: lucide/bot
---

# Claude Code settings

`.claude/settings.json` configures [Claude Code](https://claude.com/claude-code). It is **stowed to
`~/.claude/settings.json`** — Claude Code's _user-level_ settings file — so these defaults apply in every repo on
the machine, not just this one. A project-level `.claude/settings.json` in another repo layers on top of it. (See
[Memory](memory.md) for how the user-level `.claude/` files reach `$HOME` via `stow`.) The shipped file is small
and opinionated:

```json title=".claude/settings.json"
{
    "theme": "custom:cyberdream",
    "tui": "fullscreen",
    "autoMemoryEnabled": true,
    "cleanupPeriodDays": 7,
    "editorMode": "vim",
    "effortLevel": "high",
    "attribution": { "commit": "", "pr": "" },
    "autoUpdatesChannel": "stable",
    "includeGitInstructions": false,
    "plansDirectory": "./.claude/plans",
    "respectGitignore": true,
    "feedbackSurveyRate": 0,
    "permissions": {
        "defaultMode": "auto",
        "allow": ["Read(*)", "Glob", "Grep", "WebSearch", "Edit(/tmp/**)", "..."],
        "deny": ["Read(~/.aws)", "Read(~/.config/gcloud)", "Read(~/.ssh)", "Read(~/.gnupg)", "Read(**/.env*)"]
    },
    "sandbox": {
        "enabled": true,
        "filesystem": {
            "allowRead": [
                "~/Repos",
                "~/.config",
                "~/.cache",
                "~/.local/runtime",
                "~/.local/share",
                "~/.npm",
                "/opt/homebrew",
                "/tmp"
            ],
            "allowWrite": [
                "~/Repos",
                "~/.cache/agent/worktrees",
                "~/.cache/uv",
                "~/.cache/pip",
                "~/.cache/go",
                "~/.local/share/go",
                "/tmp",
                "~/.npm"
            ],
            "denyRead": ["~/.aws", "~/.config/gcloud", "~/.ssh", "~/.gnupg", "**/.env*"]
        },
        "network": {
            "allowMachLookup": [
                "com.apple.SecurityServer",
                "com.apple.trustd",
                "com.apple.trustd.agent",
                "com.apple.mDNSResponder",
                "com.apple.dnssd",
                "com.apple.system.opendirectoryd.api",
                "com.apple.system.DirectoryService.api"
            ],
            "allowedDomains": ["github.com", "api.github.com", "..."]
        },
        "autoAllowBashIfSandboxed": false,
        "allowUnsandboxedCommands": false,
        "enableWeakerNetworkIsolation": false,
        "enableWeakerNestedSandbox": false
    },
    "statusLine": {
        "type": "command",
        "command": "oh-my-posh claude --config ~/.config/oh-my-posh/claude.yaml"
    },
    "hooks": {
        "WorktreeCreate": [{ "hooks": [{ "type": "command", "command": "~/.local/share/scripts/start-worktree" }] }],
        "WorktreeRemove": [{ "hooks": [{ "type": "command", "command": "~/.local/share/scripts/end-worktree" }] }]
    },
    "worktree": { "baseRef": "head" },
    "env": {
        "IS_DEMO": "1"
    }
}
```

> **Why literal `~/` paths instead of `${XDG_CONFIG_HOME}` / `${REPOS_DIR}` / `${HOME}`?**
> Claude Code does not perform environment-variable expansion on values in `settings.json`, so
> tokens like `${REPOS_DIR}` were being treated as literal directory names and silently failing
> to match anything. The `statusLine` and `hooks` command paths use `~/` (which Claude expands
> for command fields) or absolute roots (`/opt/homebrew`, `/tmp`). The `env` block is stricter
> still — it does **not** expand `~` or `$HOME` either, so any value there must be a literal
> absolute path. That is why Go's path overrides live in [`.zshenv`](../terminal/shell.md), not
> here: a relative-looking `~/...` or `$HOME/...` `GOPATH` makes `go` fail with
> _"GOPATH entry is relative; must be absolute path"_.

## What each block does

### Theme

```json
"theme": "custom:cyberdream"
```

Points at `.claude/themes/cyberdream.json` (relative to `~/.claude/themes/`). See [Theme](theme.md).

### TUI mode

```json
"tui": "fullscreen"
```

Renders Claude Code in alternate-screen, full-terminal mode rather than the default inline
scrollback. Pairs well with tmux: the conversation owns the pane while it's active and
restores the prior terminal contents on exit.

### Memory & cleanup

```json
"autoMemoryEnabled": true,
"cleanupPeriodDays": 7
```

The model maintains persistent memory across sessions. Anything not touched for 7 days is
garbage-collected.

### Editor mode

```json
"editorMode": "vim"
```

Vim-style modal editing in the message composer.

### Default permission mode

```json
"permissions": { "defaultMode": "auto" }
```

Nests inside the [`permissions`](#permissions) object. Sets the permission mode each session starts in. `auto`
lets Claude Code pick the mode based on context rather than always opening in the default
prompt-for-everything mode — sandbox-safe
commands run without a prompt while the [sandbox](#sandbox) and the [permission](#permissions)
`deny`/`denyRead` lists remain the real boundary. Cycle modes mid-session with ++shift+tab++.

### Effort level

```json
"effortLevel": "high"
```

Biases the model toward more thorough reasoning on each turn. `high` favours deeper analysis
over latency — the right default for the configuration and infrastructure work this repo is
mostly used for.

### Attribution

```json
"attribution": { "commit": "", "pr": "" }
```

Empty strings disable Claude Code's default attribution footers on commits and PRs. The git
config's `Signed-off-by` trailer (from the prepare-commit-msg hook) is the only attribution
that lands.

### Plans

```json
"plansDirectory": "./.claude/plans"
```

When Claude Code is in plan mode, plan files write to `<repo>/.claude/plans/`. The relative
`./.claude/plans` resolves against Claude's working directory (the repo root) — deliberately relative so this
user-level setting still scopes plans per-repo rather than dumping them all under `~/.claude/plans/`. The repo's
`.gitignore` excludes `.claude/plans/` by default, and `.stowrc` ignores `^.claude/plans` so plan artifacts are
never stowed into `$HOME`.

### Permissions

There is no `additionalDirectories` list — the working set beyond the repo root is governed
entirely by `sandbox.filesystem` below. The sandbox is the real boundary; pre-approving a
directory at the permission layer without the matching `allowRead`/`allowWrite` hole only
suppresses the prompt while the underlying access still fails with `EPERM`, so the two were
collapsed onto the sandbox as the single source of truth.

`deny` takes precedence over `allow`, so it carves secrets back out of the broad `Read(*)` grant.
It lists two kinds of entry: the credential directories `Read(~/.aws)`, `Read(~/.config/gcloud)`,
`Read(~/.ssh)`, `Read(~/.gnupg)`, and the secret-file glob `Read(**/.env*)` (which matches any
dotenv file). Per the [permissions docs](https://code.claude.com/docs/en/iam), these
`Read` deny rules apply to Claude's `Read`/`Edit` tools **and** to the file-reading built-ins
Claude Code recognises in Bash (`cat`, `head`, `tail`, `sed`) — but _not_ to an arbitrary
subprocess that opens a file itself (a `python`/`node` script, `awk`, etc.). The credential
directories matter here precisely because the sandbox does **not** govern the built-in `Read`
tool: without this list, `Read(*)` would let Claude open `~/.ssh/id_rsa` directly (the file name
trips none of the globs). The same five entries are mirrored into the sandbox's `denyRead`
(see [Sandbox](#sandbox)) to cover the subprocess path the permission layer can't reach. The two
lists are kept identical: `deny` is the tool-aware block, `denyRead` is the boundary nothing
escapes.

`allow` pre-approves common, safe tool invocations so they skip the per-call permission
prompt. The sandbox (see below) is the real safety net — `allow` only controls prompts.
The list is deliberately aligned with the sandbox's `allowWrite` for Edit entries: an Edit
permission only goes on the list when the sandbox will actually let the write succeed.
Grouped by purpose:

- **Read-only Claude tools** — `Read(*)`, `Glob`, `Grep`, `WebSearch`.
- **Path-scoped edits** — `Edit(/tmp/**)` (scratch) and `Edit(~/.cache/agent/worktrees/**)`
  (agent worktrees). Edits to `~/Repos` are allowed by the sandbox but still prompt at the
  permission layer — the prompt is the deliberate friction so you stay aware of in-place
  repo edits versus agent-isolated ones.
- **Bash inspection** — `pwd`, `hostname`, `whoami`, `id`, `uptime`, `uname`, `date`, `ls`,
  `stat`, `file`, `wc`, `tree`, `which`, `type`, `command -v`, `echo`, `printf` (all with
  optional args).
- **File inspection** — `cat`, `head`, `tail`, `grep`, `rg`, `find`, `jq`, `yq` (with args).
  These are broadly allowed, but the sandbox's `denyRead` glob `**/.env*` blocks them — and any
  other subprocess — from reading dotenv files at the OS level, so `cat .env` fails with `EPERM`
  rather than leaking into a transcript.
- **Homebrew read-only** — `brew list`, `brew search`, `brew info`, `brew bundle check`.
- **Git read-only** — `git status`, `diff`, `log`, `show`, `blame`, `ls-files`, `rev-parse`,
  `config --get`, `branch --list`, `stash list`, `worktree list`, `remote -v`,
  `remote get-url` (with optional args where applicable).
- **Git mutation, reversible** — `git add`, `restore`, `checkout`, `switch`, `commit` (with
  args). Excludes `push`, `reset --hard`, `rebase`, `branch -D` — those still prompt.

Bash patterns use the documented `cmd *` (space-star) form for "command with any args".
Some commands list both `Bash(cmd)` and `Bash(cmd *)` to cover both no-arg and with-arg
invocations, since whether `*` matches an empty trailing arg isn't explicit in the docs.

### Sandbox

```json
"sandbox": {
  "enabled": true,
  "filesystem": {
    "allowRead": [
      "~/Repos",
      "~/.config",
      "~/.cache",
      "~/.local/runtime",
      "~/.local/share",
      "~/.npm",
      "/opt/homebrew",
      "/tmp"
    ],
    "allowWrite": [
      "~/Repos",
      "~/.cache/agent/worktrees",
      "~/.cache/uv",
      "~/.cache/pip",
      "~/.cache/go",
      "~/.local/share/go",
      "/tmp",
      "~/.npm"
    ],
    "denyRead": ["~/.aws", "~/.config/gcloud", "~/.ssh", "~/.gnupg", "**/.env*"]
  },
  "network": {
    "allowMachLookup": [
      "com.apple.SecurityServer",
      "com.apple.trustd",
      "com.apple.trustd.agent",
      "com.apple.mDNSResponder",
      "com.apple.dnssd",
      "com.apple.system.opendirectoryd.api",
      "com.apple.system.DirectoryService.api"
    ],
    "allowedDomains": ["github.com", "api.github.com", "..."],
    "allowLocalBinding": true,
    "allowUnixSockets": ["/tmp", "/private/tmp"]
  }
}
```

Filesystem access is **asymmetric by design**: broad reads, narrower writes.

- `allowRead` covers the source tree (`~/Repos`), the entire XDG config tree (`~/.config` —
  required for `git` to load identity and `includeIf` overlays, for `oh-my-posh` to read its
  theme, etc.), tooling caches (`~/.cache`), `XDG_RUNTIME_DIR` (`~/.local/runtime` — ephemeral
  sockets and runtime state for `nvim`, `fnm`, etc.), per-user data (`~/.local/share`), `~/.npm`
  (npm's non-XDG cache), `/opt/homebrew` (so agents can introspect what Homebrew has
  installed), and scratch (`/tmp`). The Go toolchain needs no dedicated read entry:
  [`.zshenv`](../terminal/shell.md) relocates every Go path (`GOPATH`, `GOCACHE`,
  `GOMODCACHE`, `GOENV`) under `~/.cache/go` and `~/.local/share/go`, both of which already fall
  inside the `~/.cache` and `~/.local/share` read roots above.
  `allowWrite` covers the paths agents actually need to mutate:

- `~/Repos` — the source tree itself. Agents can edit files in checked-out repos directly. The
  permission-layer prompt on `Edit(~/Repos/**)` (see Permissions above) is what keeps in-place
  edits deliberate rather than silent.
- `~/.cache/agent/worktrees` — the dominant write target when agents use
  [worktree isolation](hooks-skills.md#worktreecreate).
- `~/.cache/uv` and `~/.cache/pip` — Python package caches, required for `make docs-build` /
  `uv sync`.
- `~/.cache/go` and `~/.local/share/go` — Go's relocated caches and workspace. `~/.cache/go`
  holds the `GOCACHE`/`GOMODCACHE`/`GOENV` targets (build cache, module cache, `go env` file);
  `~/.local/share/go` is the relocated `GOPATH`, where `go install` writes binaries under `bin`.
- `/tmp` — scratch.
- `~/.npm` — npm's non-XDG cache. Listed last so the historical "read-only" stance is obvious
  from the diff: `npm install` inside an agent session needs to populate the cache, and
  leaving this out forces a prompt (or hard fail) on every fetched tarball.

`denyRead` is the kernel-level counterpart to the permission layer's `deny`, and carries the
identical five entries: it blocks reads of the listed paths no matter which tool — or which
subprocess — reaches for them, and it accepts gitignore-style globs. The two halves of the list
do different work:

- **Credential directories** — `~/.aws`, `~/.config/gcloud`, `~/.ssh`, `~/.gnupg`.
  `~/.config/gcloud` is the load-bearing sandbox entry: it sits _inside_ the allowed `~/.config`
  root, so without an explicit carve-out it would be readable; `denyRead` subtracts it back out.
  The other three aren't under any `allowRead` root to begin with, so at the sandbox layer they
  are defense-in-depth — but they still earn their keep, since the matching `deny` entries are
  what stop the built-in `Read` tool (which the sandbox doesn't govern) from opening them.
- **Secret files anywhere** — `**/.env*`. Because `~/Repos` is both readable and writable, a
  dotenv file in a checked-out repo would otherwise be `cat`-able; this glob makes the read fail
  with `EPERM` for _any_ process, closing the gap that `deny` alone leaves open for non-built-in
  subprocesses.

> **Why only `**/.env*`?** Earlier revisions also denied `\*\*/*secret*`and`\*\*/*credentials*`,
but those globs are stowed to `~/.claude/settings.json`and so apply in **every** repo on the
machine — they blocked legitimately-named source like`internal/secret/secret.go`or a doc
named`secrets.md`. The narrower `\*\*/.env*`keeps real dotenv files out of transcripts without
shadowing ordinary source. If a repo genuinely needs a broader secret block, add it to that
repo's own`.claude/settings.json` rather than the global file.

Anything outside the read list still requires explicit permission. The narrow per-tool write
holes are the model for any future additions — open the smallest path that makes a tool work
rather than re-allowing the parent.

The `network` block keeps the sandbox strict but punches the holes macOS itself needs to be
functional. `allowMachLookup` is grouped by purpose:

- **TLS / trust** — `com.apple.SecurityServer`, `com.apple.trustd`, `com.apple.trustd.agent`.
  Required for any HTTPS-using tool (git, curl, npm, uv) to validate certificates against
  the system keychain.
- **DNS** — `com.apple.mDNSResponder`, `com.apple.dnssd`. Required for hostname resolution;
  without these, anything by name fails closed.
- **Directory services** — `com.apple.system.opendirectoryd.api`,
  `com.apple.system.DirectoryService.api`. How `whoami` / `id` resolve names from UID/GID
  via `getpwuid` / `getgrgid`. (Git identity itself comes from `~/.config/git/`, which is
  covered by `allowRead` above.)

`allowUnixSockets` allowlists AF_UNIX socket paths the sandbox may `connect()` to — here
`/tmp` and `/private/tmp` (the same scratch root under its `/private` realpath), covering
the local sockets tools drop there (test fixtures, language-server / dev-server IPC, etc.).
It does **not** rescue SSH commit signing: `ssh-agent`'s socket lives under macOS's per-user
`/var/folders/...` dir, not `/tmp`, and `connect()` to it still returns EPERM — so signing
via `ssh-agent` doesn't work from inside the sandbox; run `git commit` outside Claude Code
when a signature is required. SSH-based git **remotes** still aren't a goal either; those go
through HTTPS via `allowedDomains`.

`allowLocalBinding: true` lets sandboxed processes `bind()` and `listen()` on local
addresses, so a dev server or test harness can open a port on `localhost` and be reached
from the same session (the matching `localhost` / `127.0.0.1` / `::1` entries in
`allowedDomains` cover the outbound side).

`allowedDomains` pre-approves outbound HTTPS destinations so common tools don't trigger a
permission prompt on first contact. Anything not listed still works — Claude prompts the
first time it's hit. Grouped by purpose:

- **GitHub** — `github.com`, `api.github.com`, `objects.githubusercontent.com`,
  `raw.githubusercontent.com`, `codeload.github.com`, `ghcr.io`. Covers `git` over HTTPS,
  the `gh` CLI, large-object / LFS fetches, raw file reads, archive downloads, and
  container registry pulls.
- **Forgejo / Codeberg** — `codeberg.org`, `code.forgejo.org`. Matches the
  [`includeIf` overrides](../git/config.md#includes) for Forgejo-hosted remotes.
- **Package registries** — `registry.npmjs.org` (npm metadata + tarballs),
  `pypi.org` (Python index), `files.pythonhosted.org` (wheel storage),
  `formulae.brew.sh` (Homebrew formula API).
- **Node** — `nodejs.org`. For binary-distribution installs (e.g. `nvm install`).
- **Go** — `proxy.golang.org` (module proxy), `sum.golang.org` (checksum database),
  `index.golang.org` (module index). Covers `go mod download` / `go get` over the default
  module proxy.
- **Terraform** — `registry.terraform.io`. The provider / module registry, for
  `terraform init` provider downloads.
- **Loopback** — `localhost`, `127.0.0.1`, `::1`. Lets the session reach local servers it
  starts itself (dev servers, test fixtures); pairs with `allowLocalBinding` below.

> **Note:** `allowMachLookup` is not in the official Claude Code docs at time of writing; the
> field is treated as undocumented-but-functional based on observed behavior. If a future
> Claude Code release stops honoring it, prompts for HTTPS / DNS / directory lookups will
> reappear. `allowedDomains` / `deniedDomains` _are_ documented (see
> [Sandboxing](https://code.claude.com/docs/en/sandboxing.md)).

The three explicit `false` flags below disable common escape hatches:

- `autoAllowBashIfSandboxed: false` — Bash commands still require approval.
- `allowUnsandboxedCommands: false` — no commands can run outside the sandbox.
- `enableWeakerNetworkIsolation` / `enableWeakerNestedSandbox` — both off.

### Status line

```json
"statusLine": {
  "type": "command",
  "command": "oh-my-posh claude --config ~/.config/oh-my-posh/claude.yaml"
}
```

Claude Code runs `oh-my-posh claude …` and renders the output as a status line. The
`oh-my-posh claude` subcommand consumes Claude Code's session JSON on stdin.

See [Terminal → oh-my-posh](../terminal/oh-my-posh.md#claude-code-status-line).

### Hooks

```json
"hooks": {
  "WorktreeCreate": [
    { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/start-worktree" }] }
  ],
  "WorktreeRemove": [
    { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/end-worktree" }] }
  ],
  "Stop": [
    { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/agent-tmux-status waiting" }] }
  ],
  "Notification": [
    { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/agent-tmux-status attention" }] }
  ],
  "UserPromptSubmit": [
    { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/agent-tmux-status clear" }] }
  ],
  "SessionEnd": [
    { "hooks": [{ "type": "command", "command": "~/.local/share/scripts/agent-tmux-status clear" }] }
  ]
}
```

Two hooks bracket Claude Code's worktree lifecycle, both pointing at the same shared scripts
the [`start-tmux-session`](../scripts/tmux.md#start-tmux-session) and
[`end-tmux-session`](../scripts/tmux.md#end-tmux-session) wrappers use:

- `WorktreeCreate` → [`start-worktree`](hooks-skills.md#worktreecreate) creates the worktree
  and `agent/*` branch.
- `WorktreeRemove` → [`end-worktree`](hooks-skills.md#worktreeremove) tears the worktree back
  down once Claude is done and kills any matching tmux session.

Four more hooks drive the "Claude is waiting for you" indicator via
[`agent-tmux-status`](hooks-skills.md#claude-is-waiting-indicator) — `Stop` raises a calm
`waiting` state (peach `●`) and `Notification` a louder `attention` one (bold red `󰂚`);
`UserPromptSubmit` and `SessionEnd` clear it (you replied, or the session ended). The same
script is shared with opencode's [status-indicator plugin](../opencode/plugins.md#status-indicator).

### Worktree

```json
"worktree": { "baseRef": "head" }
```

New worktrees branch from `HEAD` of the current checkout rather than the repo's default
branch — match-what-I-see-now behavior, so a worktree carries whatever state you've staged or
committed locally.

### Environment

```json
"env": {
  "IS_DEMO": "1"
}
```

`env` injects environment variables into every Claude Code session:

| Variable  | Value | Purpose                                                                 |
| --------- | ----- | ----------------------------------------------------------------------- |
| `IS_DEMO` | `1`   | Enables Claude Code's demo mode (intentional, not a leaked credential). |

Go's path overrides (`GOPATH`, `GOCACHE`, `GOMODCACHE`, `GOENV`) deliberately do **not** live
here. The `env` block does no variable expansion, so a `~/...` or `$HOME/...` value reaches `go`
verbatim and fails as a relative path (_"GOPATH entry is relative; must be absolute path"_).
They live in [`.zshenv`](../terminal/shell.md) instead, where `${XDG_DATA_HOME}` /
`${XDG_CACHE_HOME}` expand to absolute paths. The targets still land under `~/.cache/go` and
`~/.local/share/go` — exactly the roots the [sandbox](#sandbox) grants write access to — so a Go
build inside an agent session never trips an `EPERM`.
