.PHONY: requirements
requirements:
	./install.sh requirements

.PHONY: packages
packages:
	./install.sh requirements packages

.PHONY: stow
stow:
	./install.sh stow
