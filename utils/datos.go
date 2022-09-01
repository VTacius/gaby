package utils

import (
	"fmt"
	"strconv"
)

type Datos struct {
	Hostname    string
	Timestamp   marcadorTemporal
	Temperatura []float64
	Humedad     []float64
}

func conversor(lista []string, indice int) float64 {
	if len(lista)-1 >= indice {
		if resultado, err := strconv.ParseFloat(lista[indice], 64); err == nil {
			return resultado
		}
		return 0.0
	}

	return 0.0
}

func NewDatos(resultado []string, Hostname string, Timestamp marcadorTemporal) *Datos {

	temp1 := conversor(resultado, 2)
	temp2 := conversor(resultado, 3)
	hum1 := conversor(resultado, 4)
	hum2 := conversor(resultado, 5)

	Temperatura := []float64{temp1, temp2}
	Humedad := []float64{hum1, hum2}
	return &Datos{Hostname, Timestamp, Temperatura, Humedad}

}

func (d *Datos) String() string {
	return fmt.Sprintf(`Fecha: %s 
    Temperatura: %.2f - %.2f
    Humedad: %.2f - %.2f`, d.Timestamp.String(), d.Temperatura[0], d.Temperatura[1], d.Humedad[0], d.Humedad[1])
}

func (d *Datos) Mensaje() string {
	return fmt.Sprintf(`Origen: %s en %s - Temperatura: %v - Humedad: %v`, d.Hostname, d.Timestamp.String(), d.Temperatura, d.Humedad)
}

func (d *Datos) EsNuevo() bool {
	//fmt.Printf("\n\n%s es after %s\n", d.Timestamp.Ts.Format("2006-01-02 15:04:05"), d.Timestamp.Anterior.Format("2006-01-02 15:04:05"))
	resultado := d.Timestamp.Ts.After(d.Timestamp.Anterior)
	//fmt.Printf("Veredicto %t\n\n", resultado)
	return resultado
}
