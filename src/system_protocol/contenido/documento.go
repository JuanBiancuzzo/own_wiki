package contenido

import (
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strings"

	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/ast"
	p "github.com/gomarkdown/markdown/parser"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Hijos []any

type Documento struct {
	ID        bson.ObjectID `bson:"_id"`
	Elementos Hijos         `bson:"elementos"`
}

type WikiLinkNode struct {
	ast.Leaf
	Target  string
	Text    string
	IsImage bool
}

func wikiLinkParserHook(data []byte) (ast.Node, []byte, int) {
	// Regex: optional !, then [[target|text]]
	re := regexp.MustCompile(`^!?(\[\[([^\]\|]+)(\|([^\]]+))?\]\])`)
	matches := re.FindSubmatchIndex(data)
	if matches == nil {
		return nil, nil, 0
	}

	isImage := data[0] == '!'
	target := string(data[matches[4]:matches[5]])
	text := target

	if matches[8] != -1 && matches[9] != -1 {
		text = string(data[matches[8]:matches[9]])
	}

	node := &WikiLinkNode{
		Target:  strings.TrimSpace(target),
		Text:    strings.TrimSpace(text),
		IsImage: isImage,
	}

	return node, data[matches[0]:matches[1]], matches[1]
}

func NewDocumento(markdownBytes []byte) (Documento, error) {
	extensions := p.CommonExtensions | p.AutoHeadingIDs | p.SuperSubscript

	parser := p.NewWithExtensions(extensions)
	parser.Opts.ParserHook = wikiLinkParserHook

	nodoRaiz := markdown.Parse(markdownBytes, parser)
	if documento, err := RecorrerArbol(nodoRaiz); err != nil {
		var temp Documento
		return temp, err

	} else if _, ok := documento.(Documento); !ok {
		var temp Documento
		return temp, fmt.Errorf("el primer nodo no es un documento")

	} else {
		return documento.(Documento), nil
	}
}

func ObtenerHijos(nodos []ast.Node) ([]any, error) {
	hijos := make([]any, len(nodos))
	contador := 0

	for _, hijo := range nodos {
		if elemento, err := RecorrerArbol(hijo); err != nil {
			return nil, fmt.Errorf("se tuvo el error: %v", err)
		} else if elemento != nil {
			hijos[contador] = elemento
			contador++
		}
	}

	return hijos[:contador], nil
}

func GenerarContenedor(generador func(elementos Hijos) Contenedor, nodos []ast.Node) (Contenedor, error) {
	if elementos, err := ObtenerHijos(nodos); err != nil {
		var c Contenedor
		return c, err

	} else {
		return generador(elementos), nil
	}
}

func RecorrerArbol(nodo ast.Node) (any, error) {
	switch valor := nodo.(type) {
	case *ast.Document:
		if elementos, err := ObtenerHijos(valor.Children); err != nil {
			return nil, err

		} else {
			return Documento{
				ID:        bson.NewObjectID(),
				Elementos: elementos,
			}, nil
		}

	// Contenedor
	case *ast.Aside:
		if contenedor, err := GenerarContenedor(NewAside, valor.Children); err != nil {
			return nil, err
		} else {
			return contenedor.Reducir(), nil
		}
	case *ast.BlockQuote:
		return GenerarContenedor(NewBlockQuote, valor.Children)
	case *ast.Emph:
		return GenerarContenedor(NewItalic, valor.Children)
	case *ast.Paragraph:
		if contenedor, err := GenerarContenedor(NewParrafo, valor.Children); err != nil {
			return nil, err
		} else {
			return contenedor.Reducir(), nil
		}
	case *ast.Strong:
		return GenerarContenedor(NewBold, valor.Children)

	// Hojas sin texto
	case *ast.Hardbreak:
		return NewHardbreak(), nil
	case *ast.HorizontalRule:
		return NewHorizontalRule(), nil

	// Hojas con texto
	case *ast.Code:
		return NewCode(string(valor.Literal)), nil
	case *ast.HTMLBlock:
		return NewHtmlBlock(string(valor.Literal)), nil
	case *ast.HTMLSpan:
		return NewHtmlSpan(string(valor.Literal)), nil
	case *ast.Math:
		return NewMath(string(valor.Literal)), nil
	case *ast.MathBlock:
		return NewMathBlock(string(valor.Literal)), nil
	case *ast.Subscript:
		return NewSubscript(string(valor.Literal)), nil
	case *ast.Superscript:
		return NewSuperscript(string(valor.Literal)), nil
	case *ast.Text:
		return NewText(string(valor.Literal)), nil

	// el resto
	case *ast.Heading:
		if elementos, err := ObtenerHijos(valor.Children); err != nil {
			return nil, err

		} else {
			return NewHeader(uint8(valor.Level), elementos)
		}

	case *ast.List:
		elementoDesordenado := false
		items := make([]any, len(valor.Children))
		for i, elemento := range valor.Children {
			if itemAST, ok := elemento.(*ast.ListItem); !ok {
				return nil, fmt.Errorf("los elementos de una lista no son un item, son %T", elemento)

			} else if item, err := RecorrerArbol(elemento); err != nil {
				return nil, err

			} else {
				elementoDesordenado = slices.Contains([]byte{'*', '+', '-'}, itemAST.BulletChar)
				items[i] = item
			}
		}

		if elementoDesordenado {
			return NewListaDesordenada(items), nil
		}
		return NewListaOrdenada(TL_Numerico, items, valor.Start), nil

	case *ast.ListItem:
		if elementos, err := ObtenerHijos(valor.Children); err != nil {
			return nil, err

		} else {
			return NewItem(elementos), nil
		}

	case *ast.CodeBlock:
		return NewCodeBlock(string(valor.Info), string(valor.Literal)), nil

	case *ast.Image:
		return NewImagen(string(valor.Destination), string(valor.Title)), nil

	case *ast.Link:
		if elementos, err := ObtenerHijos(valor.Children); err != nil {
			return nil, err

		} else if elemento, ok := elementos[0].(HojaTexto); !ok {
			return nil, fmt.Errorf("el elemento del link no es un HojaTexto, es: %T", elementos[0])

		} else {
			return NewLink(string(valor.Destination), elemento.Texto), nil
		}

	case *ast.Callout:
		// TODO:
		fmt.Println("Esto es un callout")
		return nil, nil
	case *WikiLinkNode:

		return nil, nil

	case *ast.Table:
		fmt.Println("Esto es un table")
		return nil, nil

	case *ast.TableBody:
		fmt.Println("Esto es un table body")
		return nil, nil

	case *ast.TableCell:
		fmt.Println("Esto es un table cell")
		return nil, nil

	case *ast.TableFooter:
		fmt.Println("Esto es un tabla footer")
		return nil, nil

	case *ast.TableHeader:
		fmt.Println("Esto es un table header")
		return nil, nil

	case *ast.TableRow:
		fmt.Println("Esto es un table row")
		return nil, nil

	}

	return nil, fmt.Errorf("el tipo %T no se pudo procesar", nodo)
}

func (d Documento) CrearJson() (string, error) {
	if bytes, err := json.MarshalIndent(d, "", "   "); err != nil {
		return "", fmt.Errorf("no se pudo crear json, con error: %v", err)
	} else {
		return string(bytes), nil
	}
}
