package views

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

type Templates struct {
	templates map[string]*template.Template
}

func templatesDeArchivos(pathsTemplate []string) ([]string, error) {
	templUsuario := []string{}

	for _, pathTemplate := range pathsTemplate {
		archivo, err := os.Open(pathTemplate)
		if err != nil {
			return nil, fmt.Errorf("el template dado por '%s' no existe, se obtuvo: %v", pathTemplate, err)
		}

		infoArchvio, err := archivo.Stat()
		if err != nil {
			return nil, fmt.Errorf("el se pudo leer la informacion de '%s' no existe, se obtuvo: %v", pathTemplate, err)
		}

		if !infoArchvio.IsDir() {
			templUsuario = append(templUsuario, pathTemplate)
			continue
		}

		if templ, err := filepath.Glob(fmt.Sprintf("%s/*.html", pathTemplate)); err != nil {
			return nil, fmt.Errorf("error al obtener templates del usuario, con error: %v", err)
		} else {
			templUsuario = append(templUsuario, templ...)
		}
	}

	return templUsuario, nil
}

func NewTemplate(views []View, pathView *PathView) (*Templates, error) {
	templ := template.New("").Funcs(template.FuncMap{
		"PathViewURL":    pathView.CreateURLPathView,
		"PedirElementos": CreateURLPedido,
		"EndpointURL":    func(_ ...any) string { return "[[ERROR]]" },
	})

	templates := make(map[string]*template.Template)
	for _, view := range views {
		if templUsuario, err := templatesDeArchivos(view.PathTemplates); err != nil {
			return nil, err

		} else if templ, err := templ.ParseFiles(templUsuario...); err != nil {
			return nil, fmt.Errorf("error al parsear archivos, con error: %v", err)

		} else {
			templates[view.Nombre] = templ
		}
	}

	return &Templates{templates: templates}, nil
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	separacion := strings.SplitN(name, "/", 2)
	if separacion == nil {
		return fmt.Errorf("el nombre del bloque esta mal formateado, se dio esto: %s", name)
	}

	if templates, ok := t.templates[separacion[0]]; !ok {
		return fmt.Errorf("no existe la view dada por '%s'", separacion[0])
	} else {
		return templates.ExecuteTemplate(w, separacion[1], data)
	}
}
