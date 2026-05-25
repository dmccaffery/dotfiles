.DEFAULT_GOAL := help

.PHONY: backup
backup: ## Move conflicting configs out of $HOME into ./backups
	./backup.sh

.PHONY: restore
restore: ## Restore a chosen backup from ./backups back into $HOME (uses fzf)
	./restore.sh

.PHONY: install
install: ## Run every install stage in order
	./install.sh

.PHONY: xdg
xdg: ## Create the XDG base directories under $HOME
	./install.sh xdg

.PHONY: requirements
requirements: ## Install all prerequisites for the configuration to function
	./install.sh requirements

.PHONY: config
config: ## Apply macOS system defaults
	./install.sh config

.PHONY: stow
stow: ## Symlink configs into ~/ via stow
	./install.sh stow

.PHONY: packages
packages: ## Install all packages, including those in a selected profile
	./install.sh packages

.PHONY: shell
shell: ## Set Zsh from Homebrew as the default login shell
	./install.sh shell

.PHONY: help
help: ## Print this help
	@awk 'BEGIN { FS = ":.*?## "; printf "\nUsage: make \033[36m<target>\033[0m\n\nTargets:\n" } /^[a-zA-Z_-]+:.*?## / { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)
