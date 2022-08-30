package utils

import (
    "fmt"
    "time"
    "context"
    "strconv" 
    "github.com/influxdata/influxdb-client-go/v2"
    
    "sanidad/alortiz/gaby/peticiones"
)

type Acceso struct {
    Endpoint string
    Token string
    Organization string
    Bucket string
}

type Datos struct {
    Timestamp time.Time
    Temperatura []float64
    Humedad []float64
}

func conversor(lista []string, indice int) float64 {
    if len(lista) -1 >= indice {
        if resultado, err := strconv.ParseFloat(lista[indice], 64); err == nil {
            return resultado
        }
        return 0.0
    }

    return 0.0
}

func NewDatos(resultado []string) (Datos, error) {
    // Procesamos la marca de tiempo
    timeLayout := "01/02/2006 15:04:05"
    Fecha := resultado[0] 
    Hora :=  resultado[1]
    estampa := fmt.Sprintf("%s %s", Fecha, Hora)
    Timestamp, err := ParsearMarcaTiempo(timeLayout, estampa)
    if err != nil {
        return Datos{}, err
    }
    
    temp1 := conversor(resultado, 2) 
    temp2 := conversor(resultado, 3) 
    hum1 := conversor(resultado, 4) 
    hum2 := conversor(resultado, 5) 
    
    Temperatura := []float64{temp1, temp2}
    Humedad := []float64{hum1, hum2}
    datos := Datos{Timestamp, Temperatura, Humedad}

    return datos, nil
}

func (d *Datos) String() string {
    timeLayout := "01/02/2006 15:04:05"
    fecha := d.Timestamp.Format(timeLayout)
    return fmt.Sprintf(`Fecha: %s 
    Temperatura: %.2f - %.2f
    Humedad: %.2f - %.2f`, fecha, d.Temperatura[0], d.Temperatura[1], d.Humedad[0], d.Humedad[1])  
}

func ObtenerDatos(url_origen string, uri_historial string) (Datos, error) {
    // De la página inicial, obtenemos el enlace a los datos 
    enlace, err := peticiones.ObtenerEnlaceHistorial(url_origen, uri_historial)
    if err != nil {
        return Datos{}, err
    }
   
    // Buscamos el último dato en la página de datos
    var resultado []string
    resultado, err = peticiones.ObtenerDatosAmbientales(url_origen, enlace)
    if err != nil {
        return Datos{}, err
    }

    return NewDatos(resultado)
}

func EnviarDatos(config Acceso, datos Datos, hostname string) error {

    client := influxdb2.NewClient(config.Endpoint, config.Token)
    defer client.Close()
    
    writeAPI := client.WriteAPIBlocking(config.Organization, config.Bucket)

    marcaTiempo := time.Now().Round(time.Second * 60)
    // Create point using fluent style
    temperatura := influxdb2.NewPointWithMeasurement("temperatura").
        AddTag("host", hostname).
        AddField("temp1", datos.Temperatura[0]).
        AddField("temp2", datos.Temperatura[1]).
        SetTime(marcaTiempo)
    
    humedad := influxdb2.NewPointWithMeasurement("humedad").
        AddTag("host", hostname).
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
