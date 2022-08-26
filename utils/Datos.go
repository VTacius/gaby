package utils

import (
    "fmt"
    "time"
    "context"
    "strconv" 
    "github.com/influxdata/influxdb-client-go/v2"
)

type Configuracion struct {
    Endpoint string
    Token string
    Organization string
    Bucket string
}

type Datos struct {
    Fecha string
    Hora string
    Temperatura []float64
    Humedad []float64
}

func conversor(lista []string, indice int) float64 {
    if len(lista) -1 >= indice {
        if resultado, err := strconv.ParseFloat(lista[indice], 64); err == nil {
            return resultado
        }
        return 2.0
    }

    return 1.0
}

func NewDatos(resultado []string) Datos {
    Fecha := resultado[0] 
    Hora :=  resultado[1]
    
    temp1 := conversor(resultado, 2) 
    temp2 := conversor(resultado, 3) 
    hum1 := conversor(resultado, 4) 
    hum2 := conversor(resultado, 5) 
    
    Temperatura := []float64{temp1, temp2}
    Humedad := []float64{hum1, hum2}
    datos := Datos{Fecha, Hora, Temperatura, Humedad}

    return datos
}

func EnviarDatos(config Configuracion, datos Datos, hostname string) {

    client := influxdb2.NewClient(config.Endpoint, config.Token)
    defer client.Close()
    writeAPI := client.WriteAPIBlocking(config.Organization, config.Bucket)

    // Create point using fluent style
    temperatura := influxdb2.NewPointWithMeasurement("temperatura").
        AddTag("host", hostname).
        AddField("temp1", datos.Temperatura[0]).
        AddField("temp2", datos.Temperatura[1]).
        SetTime(time.Now().Round(time.Second * 60))
    
    humedad := influxdb2.NewPointWithMeasurement("humedad").
        AddTag("host", hostname).
        AddField("hum1", datos.Humedad[0]).
        AddField("hum2", datos.Humedad[1]).
        SetTime(time.Now().Round(time.Second * 60))
    
        err:= writeAPI.WritePoint(context.Background(), temperatura)

        fmt.Println(err)
    writeAPI.WritePoint(context.Background(), humedad)
    
}
