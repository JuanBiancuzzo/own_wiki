package dependencias

import (
	"fmt"
	"strconv"
	"strings"
)

type Date struct {
	Dia  int
	Mes  int
	Anio int
}

func NewDate(representacion string) (Date, error) {
	var date Date

	separacion := strings.Split(representacion, "-")
	if len(separacion) != 3 {
		return date, fmt.Errorf("la fecha no tiene el formato dd-mm-aaaa")
	}

	if dia, err := strconv.Atoi(separacion[2]); err != nil || dia < 1 {
		return date, fmt.Errorf("el dia de la fecha no es un numero representativo de un dia")

	} else if mes, err := strconv.Atoi(separacion[1]); err != nil || mes < 1 || mes > 12 {
		return date, fmt.Errorf("el mes de la fecha no es un numero representativo de un mes")

	} else if anio, err := strconv.Atoi(separacion[0]); err != nil {
		return date, fmt.Errorf("el anio de la fecha no es un numero representativo de un anio")

	} else {
		return Date{
			Dia:  dia,
			Mes:  mes,
			Anio: anio,
		}, nil
	}
}

func (d Date) Representacion() string {
	return fmt.Sprintf("%d-%d-%d", d.Anio, d.Mes, d.Dia)
}
