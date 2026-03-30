#!/usr/bin/env zsh
# Author: José M. C. Noronha

# Set Others Alias
alias ..='cd ..'

# DEPENDENCY AREA
source <(fzf --zsh) # ENABLE FZF

[ -f "$HOME/.local/share/coreutils/system-aliases/system-alias-zsh" ] && source "$HOME/.local/share/coreutils/system-aliases/system-alias-zsh"
