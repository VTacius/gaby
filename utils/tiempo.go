package utils

import (
    "os"
    "fmt"
    "time"
    "strings"
    "strconv"
)

func ParsearHora(timezone string, layout string, datos Datos) (time.Time) {
    marcaTemporal := fmt.Sprintf("%s %s", datos.Fecha, datos.Hora)
    
    localidad, _ := time.LoadLocation(timezone)
    tiempo, _ := time.ParseInLocation(layout, marcaTemporal, localidad) 
    
    return tiempo
}

func GuardarFechaEnArchivo(ruta string, ts time.Time){
    redondo := ts.Round(time.Second * 60)
    datos := []byte(fmt.Sprintf("%d", redondo.UnixMilli()))
    err := os.WriteFile(ruta, datos, 0700)
    if err != nil {
        fmt.Println(err)
    }
}

func LeerFechaEnArchivo(ruta string) (int64) {
    
    datos, _ := os.ReadFile(ruta)
    contenido := strings.Trim(string(datos), "\n")
    entero, _ := strconv.ParseInt(contenido, 10, 64)
    
    return entero
}

