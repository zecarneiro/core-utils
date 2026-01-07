#!/bin/bash

# Define a pasta de saída
ROOT_DIR="$*"
BINARY_DIR="$ROOT_DIR/bin"
BINARY_WIN_DIR="$BINARY_DIR/windows"
BINARY_LINUX_DIR="$BINARY_DIR/linux"
CMD_DIR="$ROOT_DIR/cmd"
COMMANDS=$(go list -f '{{if eq .Name "main"}}{{.Dir}}{{end}}' "$CMD_DIR/...")
(cd "$ROOT_DIR" && make clean) || (echo "❌ Error during make clean" && exit 1)

build_commands() {
    local extension="$1"
    local goos="$2"
    local dest_dir="$3"
    for dir in $COMMANDS; do
        name=$(basename "$dir")
        echo "🔨 Building $name..."
        GOOS="$goos" GOARCH=amd64 go build -o "$dest_dir/${name}${extension}" "$dir"
        if [ $? -ne 0 ]; then
            echo "❌ Error building $name"
            exit 1
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