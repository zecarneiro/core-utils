#!/usr/bin/env fish
# Author: José M. C. Noronha

# SHELL SCRIPTS
set -gx PATH $HOME/.local/coreutils/scripts/linux-shell $PATH

# IMPORT ALIAS
if [ -f "$HOME/.local/coreutils/alias/fish-alias" ]
  source "$HOME/.local/coreutils/alias/fish-alias"
end
alias ..='cd ..'

# DEPENDENCY AREA
fzf --fish | source # ENABLE FZF

