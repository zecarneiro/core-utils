#!/usr/bin/env bash

if ! ping -c 1 8.8.8.8 -q &>/dev/null || ! ping -c 1 8.8.4.4 -q &>/dev/null || ! ping -c 1 time.google.com -q &>/dev/null; then
    echo false
else
    echo true
fi