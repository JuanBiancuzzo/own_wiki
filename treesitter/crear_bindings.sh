#!/bin/sh
set -e

BINDINGS="/gen/treesitter"
TREESITTER_FILE=treesitter.go

TEMPLATE_FILE="/treesitter.go_temp"
TEMPLATE_SO="/template_binding.go_temp"

mkdir -p "$BINDINGS"

WORKDIR="/app"
mkdir -p "$WORKDIR"; cd "$WORKDIR"

crearBindings () {
    local LENGUAJE=$(echo "$1" | awk '{$1=$1};1')
    local LENGUAJE_CAP=$(echo "$LENGUAJE" | sed 's/.*/\u&/')

    echo "Procesando '$LENGUAJE'"

    local GIT_URL="https://github.com/tree-sitter/tree-sitter-$LENGUAJE.git"
    local DIRECTORIO="parser_$LENGUAJE"

    git clone --depth 1 "$GIT_URL" "$DIRECTORIO"
    cd "$DIRECTORIO/src"

    gcc -shared -o "lib_treesitter_$LENGUAJE.so" -fPIC *.c

    cp "lib_treesitter_$LENGUAJE.so" "$BINDINGS/lib_treesitter_$LENGUAJE.so"

    sed "s/__lenguaje__/$LENGUAJE/g" "$TEMPLATE_SO" > temp
    sed "s/__Lenguaje__/$LENGUAJE_CAP/g" temp >> "$BINDINGS/$TREESITTER_FILE"

    cd "$WORKDIR"
    echo "Terminando con los bindings de $LENGUAJE"

    return 0
}

cat "$TEMPLATE_FILE" > "$BINDINGS/$TREESITTER_FILE"
crearBindings javascript

echo "Terminando de procesar todos los bindings"