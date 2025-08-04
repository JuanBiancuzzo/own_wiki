package parser

import (
	// Este todavia no funciona correctamente -> investigar
	tsmk "github.com/tree-sitter-grammars/tree-sitter-markdown/bindings/go"
	ts "github.com/tree-sitter/go-tree-sitter"
)

type MarkdownParser struct {
	Parser *ts.Parser
}

func NewMarkdownParser() *MarkdownParser {
	parser := ts.NewParser()
	parser.SetLanguage(ts.NewLanguage(tsmk.Language()))

	return &MarkdownParser{
		Parser: parser,
	}
}

func (mkp *MarkdownParser) Parsear(texto []byte, arbolPrevio *ts.Tree) *ts.Tree {
	return mkp.Parser.Parse(texto, arbolPrevio)
}

func (mkp *MarkdownParser) Close() {
	mkp.Parser.Close()
}
