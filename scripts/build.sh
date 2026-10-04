#!/bin/bash

CMD_NAME="$1"
ROOT_DIR="$PWD"
SCRIPTS_DIR="$ROOT_DIR/scripts"

echo "🚀 Building..."
. "$SCRIPTS_DIR/installer-and-others/builder.sh" "$CMD_NAME" "$ROOT_DIR"
echo "✅ Building finished!"
