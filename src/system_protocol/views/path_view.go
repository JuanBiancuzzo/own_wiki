package views

import (
	"fmt"
	"strings"
)

type claves []string

type PathView struct {
	Views map[string]claves
}

func NewPathView() *PathView {
	return &PathView{
		Views: make(map[string]claves),
	}
}

func (pv *PathView) AgregarView(view string, claves []string) error {
	if _, ok := pv.Views[view]; ok {
		return fmt.Errorf("ya se cargo ese parametro")
	}

	pv.Views[view] = claves
	return nil
}

func (pv *PathView) CreateURLPathView(view string, valores ...any) string {
	if claves, ok := pv.Views[view]; !ok {
		return "ERROR - No existe view"

	} else if len(claves) > len(valores) {
		return "ERROR - No suficientes parametros"

	} else {
		claveValor := make([]string, len(claves))
		for i, clave := range claves {
			valor := valores[i]
			claveValor = append(claveValor, fmt.Sprintf("%s=%v", clave, valor))
		}

		parametros := ""
		if len(claveValor) > 0 {
			parametros = fmt.Sprintf("?%s", strings.Join(claveValor, "&"))
		}

		return fmt.Sprintf("/%s%s", view, parametros)
	}
}
