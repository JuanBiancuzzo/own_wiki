package webview

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"

	v "own_wiki/system_protocol/views"

	"github.com/labstack/echo/v4"
)

const PATH_TEMPL = "../templates"

type Templates struct {
	templates *template.Template
}

func NewTemplate(pathTemplate string) (*Templates, error) {
	funcMaps := template.FuncMap{
		"PathViewURL": v.CreateURL,
	}

	if templUsuario, err := filepath.Glob(fmt.Sprintf("%s/*.html", pathTemplate)); err != nil {
		return nil, fmt.Errorf("error al obtener templates del usuario, con error: %v", err)

	} else if template, err := template.ParseFiles(templUsuario...); err != nil {
		return nil, fmt.Errorf("error al parsear archivos, con error: %v", err)
	} else {
		return &Templates{
			templates: template.Funcs(funcMaps),
		}, nil
	}
}

func (t *Templates) Render(w io.Writer, name string, data any, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
