#!/usr/bin/env bash
# Author: José M. C. Noronha
# shellcheck disable=SC2154
# shellcheck disable=SC1090
# shellcheck disable=SC2155

# Global Vars
declare __COREUTILS_SCRIPT_DIR__=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" &> /dev/null && pwd)
declare __COREUTILS_LIBS_SCRIPT_DIR__="${__COREUTILS_SCRIPT_DIR__}/../libs"


# Do if on Windows SO
if [[ "$OSTYPE" == "cygwin" ]]||[[ "$OSTYPE" == "msys" ]]||[[ "$OSTYPE" == "win32" ]]; then
  export LANG=C.UTF-8
fi

# Set Others Alias
alias ..='cd ..'

# DEPENDENCY AREA
eval "$(fzf --bash)" # ENABLE FZF

[ -f "$HOME/.local/share/coreutils/system-aliases/system-alias-bash" ] && source "$HOME/.local/share/coreutils/system-aliases/system-alias-bash"

build_prompt() {
  style_to_use="-1"
  if ! case "$(uname -s)" in
    CYGWIN*|MINGW*|MSYS*) true ;;
    Linux*)
      grep -qi microsoft /proc/version 2>/dev/null && true
      ;;
    *) false ;;
    esac
  then
    if command -v prompt-style >/dev/null 2>&1; then
      style_to_use="$(prompt-style -s)"
    fi
  fi
  local style_script="${__COREUTILS_LIBS_SCRIPT_DIR__}/style-prompt.sh"
  . "${style_script}" "$style_to_use"
}
PROMPT_COMMAND=build_prompt
