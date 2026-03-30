#!/usr/bin/env fish
# Author: José M. C. Noronha

# Set Others Alias
alias ..='cd ..'

# DEPENDENCY AREA
fzf --fish | source # ENABLE FZF

test -f ~/.local/share/coreutils/system-aliases/system-alias-fish.fish; and source  ~/.local/share/coreutils/system-aliases/system-alias-fish.fish

