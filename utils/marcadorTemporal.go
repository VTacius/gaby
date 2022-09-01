package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

type marcadorTemporal struct {
	Ruta     string
	Ts       time.Time
	Layout   string
	Anterior time.Time
}

func NewMarcadorTemporal(directorio string, hostname string, Estampa string) *marcadorTemporal {
	Ruta := fmt.Sprintf("%s/%s", directorio, hostname)
	Ts := time.Now()
	Anterior, err := leerUltimoTS(Ruta)
	// Este será el manejo del error: Si no puedo marcar el error, le mando la hora actual
	if err != nil {
		Anterior = time.Now().Round(time.Second * 60).Add(-2 * time.Minute)
	}
	return &marcadorTemporal{Ruta, Ts, Estampa, Anterior}
}

func (t *marcadorTemporal) ActualizarEstampa(fecha string, hora string) (*marcadorTemporal, error) {
	estampa := fmt.Sprintf("%s %s", fecha, hora)
	zonaHoraria, err := time.LoadLocation("Local")
	if err != nil {
		return t, err
	}

	Ts, errorParse := time.ParseInLocation(t.Layout, estampa, zonaHoraria)
	if err != errorParse {
		return t, err
	}

	t.Ts = Ts.Round(time.Second * 60)

	return t, nil
}

func (t *marcadorTemporal) GuardarActual() error {
	redondeo := t.Ts.Round(time.Second * 60)
	datos := []byte(fmt.Sprintf("%d", redondeo.Unix()))
	return os.WriteFile(t.Ruta, datos, 0700)
}

func (t *marcadorTemporal) String() string {
	return t.Ts.Format(t.Layout)
}

func leerUltimoTS(ruta string) (time.Time, error) {
	contenido, err := os.ReadFile(ruta)
	if err != nil {
		return time.Time{}, err
	}

	datos := strings.Trim(string(contenido), "\n")
	if ums, err := strconv.ParseInt(datos, 10, 64); err != nil {
		return time.Time{}, err
	} else {
		return time.Unix(ums, 0), nil
	}

}
