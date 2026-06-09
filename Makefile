.DEFAULT_GOAL := help

# Run Node CLIs straight from node_modules so package-lock.json is the source of
# truth. Do not use npx or global developer-installed tools.
NPMBIN := ./node_modules/.bin

.PHONY: backup
backup: ## Move conflicting configs out of $HOME into ./backups
	@ ./backup.sh

.PHONY: restore
restore: ## Restore a chosen backup from ./backups back into $HOME (uses fzf)
	@ ./restore.sh

.PHONY: install
install: ## Run every install stage in order
	@ ./install.sh

.PHONY: xdg
xdg: ## Create the XDG base directories under $HOME
	@ ./install.sh xdg

.PHONY: requirements
requirements: ## Install all prerequisites for the configuration to function
	@ ./install.sh requirements

.PHONY: config
config: ## Apply macOS system defaults
	@ ./install.sh config

.PHONY: stow
stow: ## Symlink configs into ~/ via stow
	@ ./install.sh stow

.PHONY: build
build: ## Build the dot CLI and link its applets into ~/.local/share/scripts
	@ ./install.sh build

.PHONY: packages
packages: ## Install all packages, including those in a selected profile
	@ ./install.sh packages

.PHONY: shell
shell: ## Set Zsh from Homebrew as the default login shell
	@ ./install.sh shell

node_modules: package.json package-lock.json
	@ npm ci
	@ touch node_modules

.PHONY: fmt
fmt: node_modules ## Install locked npm deps and auto-format the repo with prettier
	@ $(NPMBIN)/prettier --write .

.PHONY: scrub
scrub: ## Remove generated private state from tracked config files
	@ ./hack/scrub-codex-config-state.sh

.PHONY: lint
lint: SHELL := bash
lint: .SHELLFLAGS := -eu -o pipefail -c
lint: fmt ## Format the repo, then run shellcheck on shell scripts and markdownlint on every .md file
	@ shopt -s globstar nullglob; shellcheck --severity=warning --external-sources \
		install.sh restore.sh backup.sh \
		hack/*.sh \
		setup/**/*.sh \
		stow/.local/share/scripts/* \
		stow/.config/git/template/hooks/* \
		stow/.ssh/rc
	@ $(NPMBIN)/markdownlint-cli2 '**/*.md'

.PHONY: docs-serve
docs-serve: ## Sync deps and serve the docs site at http://localhost:8000
	@ uv sync
	@ uv run zensical serve

.PHONY: docs-build
docs-build: lint ## Lint, then sync deps and build the docs site into ./site
	@ uv sync
	@ uv run zensical build --clean

.PHONY: commit
commit: ## Run ./commit.sh when the current batch has one
	@ if [ -f ./commit.sh ]; then ./commit.sh; else printf 'No commit.sh found; skipping.\n'; fi

.PHONY: pr
pr: scrub fmt lint commit ## Scrub generated state, format, lint, then run commit.sh if present

.PHONY: upgrade
upgrade: ## Upgrade npm + uv deps (bypasses dependabot cooldown — confirms first)
	@ printf '\n'
	@ printf 'WARNING: `make upgrade` bypasses the 7-day dependabot cooldown configured\n'
	@ printf 'in .github/dependabot.yaml. Pulling fresh releases from npm and PyPI\n'
	@ printf 'registries before that window elapses can expose this repo to supply-\n'
	@ printf 'chain attacks (typosquat releases, hijacked packages, malicious patch\n'
	@ printf 'versions). Dependabot batches minor/patch updates into one PR after the\n'
	@ printf 'cooldown so the wider ecosystem has time to flag malicious releases.\n\n'
	@ printf 'Prefer merging the dependabot PR. Continue only if you have a reason.\n\n'
	@ printf 'Continue with upgrade? [y/N] '; read ans </dev/tty; \
		case "$$ans" in \
			y|Y|yes|Yes|YES) ;; \
			*) echo "Aborted."; exit 1 ;; \
		esac
	@ npm update
	@ uv sync --upgrade

.PHONY: help
help: ## Print this help
	@ awk 'BEGIN { FS = ":.*?## "; printf "\nUsage: make \033[36m<target>\033[0m\n\nTargets:\n" } /^[a-zA-Z_-]+:.*?## / { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
