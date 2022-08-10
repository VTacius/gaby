package main

import (
    "sanidad/alortiz/gaby/peticiones"
    "os"
    "fmt"
)

func main() {
    URL_ORIGEN_DATOS := os.Getenv("URL_ORIGEN_DATOS")
    URI_HISTORIAL := "PageHislog.html"

    
    enlaceDatos, err := peticiones.ObtenerEnlaceHistorial(URL_ORIGEN_DATOS, URI_HISTORIAL)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println(enlaceDatos)
    
    var resultado []string
    resultado, err = peticiones.ObtenerDatosAmbientales(URL_ORIGEN_DATOS, enlaceDatos)
    if err != nil {
        fmt.Println(err)
        return
    }
    
    fmt.Println(resultado)
}
