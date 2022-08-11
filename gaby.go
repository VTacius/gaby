package main

import (
    "os"
    "fmt"
    "time"
    "sanidad/alortiz/gaby/peticiones"
    "sanidad/alortiz/gaby/utils"
)

//Tenés la fecha del sistema y del nodo. La fecha del nodo es la que vamos a guardar para comparar
// pero la que vamos a enviar al sistema es la del sistema, que debería ser más actual
// Ambas tienen que ser redondeadas

func main() {
    zonaHoraria := "America/El_Salvador"
    timeLayout := "01/02/2006 15:04:05"
    ficheroTemporal := "/var/lib/gaby"

    hostname := "10.0.0.9"
   
    URL_ORIGEN_DATOS := os.Getenv("URL_ORIGEN_DATOS")
    URI_HISTORIAL := "PageHislog.html"
    Endpoint := os.Getenv("GABY_ENDPOINT")
    Token := os.Getenv("GABY_TOKEN")
    Organization := "sanidad"
    Bucket := "ambientales"
   
    enlaceDatos, err := peticiones.ObtenerEnlaceHistorial(URL_ORIGEN_DATOS, URI_HISTORIAL)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    var resultado []string
    resultado, err = peticiones.ObtenerDatosAmbientales(URL_ORIGEN_DATOS, enlaceDatos)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    datos := utils.NewDatos(resultado)
    config := utils.Configuracion{Endpoint, Token, Organization, Bucket}
    
    horaActual := utils.ParsearHora(zonaHoraria, timeLayout, datos) 
    horaAnterior := utils.LeerFechaEnArchivo(ficheroTemporal)
    
    utils.EnviarDatos(config, datos, hostname)
       
    for i := 0; i < 10; i++ {
        if horaActual.UnixMilli() >  horaAnterior {
            utils.GuardarFechaEnArchivo(ficheroTemporal, horaActual)
        } else {
            fmt.Println("Todavía no se puede guardar")
            time.Sleep(1 * time.Second)
        }
    }
}

