package utils

import (
    "os"
    "fmt"
    "time"
    "strings"
    "strconv"
)

func ParsearMarcaTiempo(layout string, estampa string) (time.Time, error) {

    localidad, err := time.LoadLocation("Local")
    if err != nil {
        return time.Time{}, err
    }

    tiempo, errorParse := time.ParseInLocation(layout, estampa, localidad)

    return tiempo, errorParse
}

func GuardarFechaEnArchivo(ruta string, ts time.Time) (error){
    redondo := ts.Round(time.Second * 60)
    datos := []byte(fmt.Sprintf("%d", redondo.UnixMilli()))
    err := os.WriteFile(ruta, datos, 0700)
    return err
}

func LeerFechaEnArchivo(ruta string) (int64, error) {
    datos, err := os.ReadFile(ruta)
    if err != nil {
        return 0, err
    }
    contenido := strings.Trim(string(datos), "\n")
    entero, err := strconv.ParseInt(contenido, 10, 64)
    
    return entero, err
}

