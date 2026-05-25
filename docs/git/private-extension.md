---
icon: lucide/lock
---

# Private overlay

The shipped git config ends with an unconditional include:

```ini title=".config/git/config (last lines)"
[include]
    path = ~/.config/private/git/config
```

Anything you don't want to commit to a public dotfiles repo lives there. Common uses:

- `user.name` / `user.email`
- Per-host emails (work vs personal)
- `commit.gpgsign = true` toggles when only some accounts have keys
- Company-specific URL rewrites (`url.<base>.insteadOf`)
- Internal CA roots for HTTPS git remotes

## Stow pattern

The recommended layout is a **separate private repository** that uses `stow` to symlink
`~/.config/private`:

```text
my-private-dotfiles/        # a separate, private repo
└── .config/
    └── private/
        ├── git/
        │   ├── config           # your name/email + host overrides
        │   └── work.gitconfig   # included from config for work hosts
        └── ssh/
            └── config           # private SSH host aliases
```

Then in that private repo:

```sh
stow --target=$HOME .
```

The public dotfiles see the symlink as if the file were there all along, but the file's
_contents_ never enter the public repository.

## Example private config

```ini title="~/.config/private/git/config"
[user]
    name = Your Name
    email = you@personal.example
    signingKey = key::ssh-ed25519-sk AAAAC3Nz... yubikey

[includeIf "gitdir:~/Repos/work/**"]
    path = ~/.config/private/git/work.gitconfig
```

```ini title="~/.config/private/git/work.gitconfig"
[user]
    email = you@company.example
    signingKey = key::ssh-ed25519-sk AAAAC3Nz... work-yubikey

[url "git@github.com:company/"]
    insteadOf = https://github.com/company/
```

!!! warning "No credentials here either"
Private repos are _not_ a credential store. Use the system keychain
(`git-credential-manager`) and resident SSH keys on a hardware token. The private overlay
is for _personal_ information (name, email, internal URLs) — not secrets.
