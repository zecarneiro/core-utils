#!/usr/bin/env bash

if [ "$1" = "install" ]; then
    shift
    pip install --break-system-packages "$@"
elif [ "$1" = "uninstall" ]; then
    shift
    pip uninstall --break-system-packages "$@"
else
    pip "$@"
fi
