package webview

import (
	"fmt"
	"strings"
)

type PathView struct {
	View       string
	Texto      string
	Parametros map[string]string // Clave-valor
}

func NewPathView(view, texto string) *PathView {
	return &PathView{
		View:       view,
		Texto:      texto,
		Parametros: make(map[string]string),
	}
}

func (pv *PathView) AgregarParametro(clave, valor string) error {
	if _, ok := pv.Parametros[clave]; ok {
		return fmt.Errorf("ya se cargo ese parametro")
	}

	pv.Parametros[clave] = valor
	return nil
}

func CreateURL(pathView *PathView) string {
	claveValor := []string{}
	for clave := range pathView.Parametros {
		valor := pathView.Parametros[clave]
		claveValor = append(claveValor, fmt.Sprintf("%s=%s", clave, valor))
	}

	return fmt.Sprintf("/%s?%s", pathView.View, strings.Join(claveValor, "&"))
}
