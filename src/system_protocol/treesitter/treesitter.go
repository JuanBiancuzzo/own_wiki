package treesitter

import (
	"fmt"

	"github.com/ebitengine/purego"
)


func LanguageJavascript() (uintptr, error) {
	path := "./lib_treesitter_javascript.so"
	lib, err := purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
	if err != nil {
		var pointer uintptr
        return pointer, fmt.Errorf("error al cargar el grammar de javascript como .so")
    }

	var lenguaje func() uintptr
	purego.RegisterLibFunc(&lenguaje, lib, "tree_sitter_javascript")

	return lenguaje(), nil
}