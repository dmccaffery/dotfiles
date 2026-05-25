---
icon: lucide/git-branch
---

# git-town

[git-town](https://www.git-town.com) is a workflow tool for parallel branch development —
syncing chains of feature branches against `main`, shipping in order, and keeping branch
parents tracked. This repo aliases every git-town command directly under `git`, so the muscle
memory stays `git <verb>`.

## Aliases

```ini title=".config/git/config (aliases section)"
[alias]
    append      = town append
    browse      = "!f() { git-town repo upstream 1>/dev/null 2>&1 || git-town repo 1>/dev/null 2>&1; }; f"
    compress    = town compress
    contribute  = town contribute
    delete      = town delete
    diff-parent = town diff-parent
    done        = "!f() { git checkout main; git sync; }; f"
    hack        = town hack
    observe     = town observe
    park        = town park
    pr          = town prepare
    prepend     = town prepend
    propose     = town propose
    prototype   = town hack --prototype
    rename      = town rename
    repo        = town repo
    set-parent  = town set-parent
    sync        = town sync --all
    start       = town hack
```

| Alias                     | What it does                                                                                                 |
| ------------------------- | ------------------------------------------------------------------------------------------------------------ |
| `git hack <branch>`       | Create a new feature branch off `main`, with `main` as its parent.                                           |
| `git start <branch>`      | Synonym for `hack` — more natural opening verb.                                                              |
| `git append <branch>`     | Stack a new branch on top of the current one (current becomes parent).                                       |
| `git prepend <branch>`    | Insert a branch between the current branch and its parent.                                                   |
| `git sync`                | Pull `main`, rebase every local feature branch (`--all`) on top of `main` and on top of any parent branches. |
| `git propose`             | Open a PR for the current branch.                                                                            |
| `git ship`/`done`         | `git checkout main && git sync` (after merge).                                                               |
| `git compress`            | Squash all commits on the current branch into one.                                                           |
| `git contribute <branch>` | Mark a branch as "I'm just contributing" — sync stays minimal.                                               |
| `git park`                | Stop syncing this branch (e.g., on hiatus).                                                                  |
| `git observe`             | Make a branch read-only locally.                                                                             |
| `git delete <branch>`     | Delete a branch and its lineage cleanly.                                                                     |
| `git rename`              | Rename a branch _and_ update upstream tracking.                                                              |
| `git set-parent`          | Re-point the current branch's parent.                                                                        |
| `git diff-parent`         | Diff against the parent branch (not `main`).                                                                 |
| `git repo`                | Open the repo's web URL in a browser.                                                                        |
| `git browse`              | Like `repo`, but prefers the upstream remote.                                                                |

## Git-town config

```ini title=".config/git/config ([git-town] section)"
[git-town]
    perennial-regex         = ^v([[:digit:]]\\.?)+$
    push-hook               = false
    ship-delete-tracking-branch = false
    sync-feature-strategy   = rebase
    sync-tags               = false
    sync-upstream           = true
    share-new-branches      = push
    github-connector        = gh
    dev-remote              = origin
    main-branch             = main
```

| Setting                          | Effect                                                                    |
| -------------------------------- | ------------------------------------------------------------------------- |
| `perennial-regex`                | Branches matching `^v1`, `^v2.3`, … are perennial (never synced to main). |
| `push-hook = false`              | Skip pre-push hooks during sync — common pattern for speed.               |
| `sync-feature-strategy = rebase` | Rebase, don't merge, when syncing feature branches.                       |
| `share-new-branches = push`      | New branches are pushed to `origin` immediately.                          |
| `github-connector = gh`          | Use the `gh` CLI for PR operations.                                       |

Provider-specific `forge-type` is set per-host via the includeIf'd `github.gitconfig` and
`forgejo.gitconfig`.
