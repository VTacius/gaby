package utils

import (
	"fmt"
	"sanidad/alortiz/gaby/peticiones"
)

type Origen struct {
	Url      string
	Hostname string
}

func NewOrigen(origen string, usuario string, password string) *Origen {
	Url := fmt.Sprintf("http://%s:%s@%s", usuario, password, origen)
	return &Origen{Url: Url, Hostname: origen}
}

func (o *Origen) ObtenerDatos(uri_historial string, estampa marcadorTemporal) (*Datos, error) {
	// De la página inicial, obtenemos el enlace a los datos
	enlace, err := peticiones.ObtenerEnlaceHistorial(o.Url, uri_historial)
	if err != nil {
		return &Datos{}, err
	}

	// Buscamos el último dato en la página de datos
	var resultado []string
	resultado, err = peticiones.ObtenerDatosAmbientales(o.Url, enlace)
	if err != nil {
		return &Datos{}, err
	}

	// TODO: Manejar este posible error
	estampa.ActualizarEstampa(resultado[0], resultado[1])

	return NewDatos(resultado, o.Hostname, estampa), nil

}
