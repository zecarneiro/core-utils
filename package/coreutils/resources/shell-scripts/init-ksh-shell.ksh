#!/usr/bin/env ksh
# Author: José M. C. Noronha

# SHELL SCRIPTS
export PATH="$HOME/.local/coreutils/scripts/linux-shell:$PATH"

# IMPORT ALIAS
if [ -f "$HOME/.local/coreutils/alias/ksh-alias" ]; then
  source "$HOME/.local/coreutils/alias/ksh-alias"
fi
alias ..='cd ..'
