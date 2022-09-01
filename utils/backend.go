package utils

import (
	"context"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

type Backend struct {
	Endpoint     string
	Token        string
	Organizacion string
	Bucket       string
}

func NewBackend(Endpoint string, Token string, Organizacion string, Bucket string) *Backend {
	return &Backend{Endpoint, Token, Organizacion, Bucket}
}

func (b *Backend) Enviar(datos Datos) error {
	client := influxdb2.NewClient(b.Endpoint, b.Token)
	defer client.Close()

	writeAPI := client.WriteAPIBlocking(b.Organizacion, b.Bucket)

	marcaTiempo := time.Now().Round(time.Second * 60)
	// Create point using fluent style
	temperatura := influxdb2.NewPointWithMeasurement("temperatura").
		AddTag("host", datos.Hostname).
		AddField("temp1", datos.Temperatura[0]).
		AddField("temp2", datos.Temperatura[1]).
		SetTime(marcaTiempo)

	humedad := influxdb2.NewPointWithMeasurement("humedad").
		AddTag("host", datos.Hostname).
		AddField("hum1", datos.Humedad[0]).
		AddField("hum2", datos.Humedad[1]).
		SetTime(marcaTiempo)

	err := writeAPI.WritePoint(context.Background(), temperatura)
	if err != nil {
		return err
	}

	err = writeAPI.WritePoint(context.Background(), humedad)
	return err
}
