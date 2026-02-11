#!/bin/bash

# Define a pasta de saída
CMD_NAME="$1"; shift
ROOT_DIR="$*"
BINARY_DIR="$ROOT_DIR/bin"
BINARY_WIN_DIR="$BINARY_DIR/windows"
BINARY_LINUX_DIR="$BINARY_DIR/linux"
CMD_DIR="$ROOT_DIR/cmd"
CONFIG_DIR="$ROOT_DIR/configs"

# Set safe dir on git. Necessary for go
if [[ -z "$(git config --global safe.directory)" ]]; then
    echo "Set safe directory for git"
    git config --global --add safe.directory '*'
    sleep 2
fi

COMMANDS=$(go list -f '{{if eq .Name "main"}}{{.Dir}}{{end}}' "$CMD_DIR/...")
if [[ -z "${CMD_NAME}" ]]; then
    (rm -rf "$BINARY_DIR") || (echo "❌ Error during clean" && exit 1)
fi

build_commands() {
    local extension="$1"
    local goos="$2"
    local dest_dir="$3"
    local SKIP_ICON="⚠️"
    for dir in $COMMANDS; do
        name=$(basename "$dir")
        if [[ -z "${CMD_NAME}" ]]||[[ "${CMD_NAME}" == "${name}" ]]; then
            skip_msg="$SKIP_ICON  Skip $name"
            if [[ "${goos}" == "windows" ]]&&[[ "$(cat "$CONFIG_DIR/windows-cmds-build-ignore" | grep -c "^$name")" -gt 0 ]]; then
                echo "${skip_msg}, Linux only"
                continue
            fi
            if [[ "${goos}" == "linux" ]]&&[[ "$(cat "$CONFIG_DIR/linux-cmds-build-ignore" | grep -c "^$name")" -gt 0 ]]; then
                echo "${skip_msg}, Windows only"
                continue
            fi
            echo "🔨 Building $name..."
            GOOS="$goos" GOARCH=amd64 go build -o "$dest_dir/${name}${extension}" "$dir"
            if [ $? -ne 0 ]; then
                echo "❌ Error building $name"
                exit 1
            fi
        fi
    done
}

echo "---------------------------------------"
echo "Init build of the projet COREUTILS for WINDOWS"
echo "---------------------------------------"
mkdir -p "$BINARY_WIN_DIR"
build_commands ".exe" "windows" "$BINARY_WIN_DIR"

echo "---------------------------------------"
echo "Init build of the projet COREUTILS for LINUX"
echo "---------------------------------------"
mkdir -p "$BINARY_LINUX_DIR"
build_commands "" "linux" "$BINARY_LINUX_DIR"

echo "---------------------------------------"
echo "✅ All binaries generated in: $BINARY_DIR"