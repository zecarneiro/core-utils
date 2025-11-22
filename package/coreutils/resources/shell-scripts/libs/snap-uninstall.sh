#!/usr/bin/env bash
# Author: José M. C. Noronha

declare config_dir="$HOME/snap/$1"
declare config_system_dir="/snap/$1"
eval-advanced "sudo snap remove --purge '$1'"
eval-advanced "snap saved"
read -p "Insert the number on the line of App(ENTER TO SKIP): " userInput
if [[ -n "${userInput}" ]]; then
    eval-advanced "sudo snap forget ${userInput}"
fi
if [ -d "$config_dir" ]; then
    eval-advanced "rm -rf \"$config_dir\""
fi
if [ -d "$config_system_dir" ]; then
    eval-advanced "sudo rm -rf \"$config_system_dir\""
fi
