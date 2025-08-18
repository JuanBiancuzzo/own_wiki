package views

import (
	"fmt"
	"strings"
)

type PathView struct {
	View       string
	Parametros DataView // Clave-valor
}

func NewPathView(view string) *PathView {
	return &PathView{
		View:       view,
		Parametros: make(DataView),
	}
}

func (pv *PathView) AgregarParametro(clave string, valor any) error {
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
		claveValor = append(claveValor, fmt.Sprintf("%s=%v", clave, valor))
	}

	parametros := ""
	if len(claveValor) > 0 {
		parametros = fmt.Sprintf("?%s", strings.Join(claveValor, "&"))
	}

	return fmt.Sprintf("/%s%s", pathView.View, parametros)
}
