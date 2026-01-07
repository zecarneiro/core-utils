#!/bin/bash

# Define a pasta de saída
ROOT_DIR="$*"
BINARY_DIR="$ROOT_DIR/bin"
CMD_DIR="$ROOT_DIR/cmd"
COMMANDS=$(go list -f '{{if eq .Name "main"}}{{.Dir}}{{end}}' "$CMD_DIR/...")

echo "---------------------------------------"
echo "Init build of the projet COREUTILS"
echo "---------------------------------------"
mkdir -p "$BINARY_DIR"

for dir in $COMMANDS; do
    # Extrai o nome do binário (nome da pasta)
    name=$(basename "$dir")
    echo "🔨 Building $name..."
    # Executa o build. O comando aponta para o diretório absoluto retornado pelo go list
    go build -o "$BINARY_DIR/$name" "$dir"
    # Verifica se o último comando (go build) falhou
    if [ $? -ne 0 ]; then
        echo "❌ Error building $name"
        exit 1
    fi
done

echo "---------------------------------------"
echo "✅ All binaries generated in: $BINARY_DIR"