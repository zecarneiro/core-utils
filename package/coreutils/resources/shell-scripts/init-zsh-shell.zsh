#!/usr/bin/env zsh
# Author: José M. C. Noronha

# SHELL SCRIPTS
export PATH="$HOME/.local/coreutils/scripts/linux-shell:$PATH"

# IMPORT ALIAS
if [ -f "$HOME/.local/coreutils/alias/zsh-alias" ]; then
  source "$HOME/.local/coreutils/alias/zsh-alias"
fi
alias ..='cd ..'

# DEPENDENCY AREA
source <(fzf --zsh) # ENABLE FZF
