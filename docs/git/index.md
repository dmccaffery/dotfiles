---
icon: simple/git
---

# Git

User-level git config plus the security-key signing workflow, git-town aliases for parallel
branch development, host-specific auth strategies, and a slot for private overlays.

![lazygit browsing a repository with the cyberdream theme](../assets/images/lazygit.png)

| Page                                                | Purpose                                                               |
| --------------------------------------------------- | --------------------------------------------------------------------- |
| [Git config](config.md)                             | `.config/git/config` — the canonical user-level config.               |
| [git-town](git-town.md)                             | Branch-chain workflow aliased directly under `git`.                   |
| [Authentication](auth.md)                           | Per-host auth strategy (GitHub OAuth, GCM for Codeberg/GitLab).       |
| [Signing & security keys](signing-security-keys.md) | SSH-signed commits from a YubiKey resident key.                       |
| [Private overlay](private-extension.md)             | Unconditional include at the end of `.config/git/config` for secrets. |
