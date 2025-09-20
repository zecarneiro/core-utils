#!/usr/bin/env bash

# Se não foi passado comando: erro
if [ $# -eq 0 ]; then
    echo "Use: sudoexec command [args...]"
    exit 1
fi

cmd="$1"
shift

# Resolver o caminho completo
cmd_path="$(command -v "$cmd")"

if [ -z "$cmd_path" ]; then
    echo "Erro: comando '$cmd' não encontrado no PATH."
    exit 1
fi

# Executar com sudo -E preservando ambiente
exec sudo -E "$cmd_path" "$@"