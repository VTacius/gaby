package main

import (
    "os"
    "fmt"
    "time"
    "sanidad/alortiz/gaby/utils"
    "github.com/jaffee/commandeer"
)

// TODO: Hay un problema con esta estructura de datos, y es que estas uniendo acceso con datos
//  y sobre todo, necesitas poner lo del archivo en 
type Gaby struct {
    Envio bool `help:"Hace el envío de los datos hacia un sistema influxDB"`
    Origen string `help:"URL del sistema EATON"`
    UsuarioOrigen string `help:"Usuario para acceder al sistema EATON"`
    PasswordOrigen string `help:"Password para acceder al sistema EATON"`
    Acceso utils.Acceso `flag:"-"`
    DirectorioTemporal string `flag:"-"`
    UriHistorial string `flag:"-"`
}

func NewGaby(acceso utils.Acceso, directorioTemporal string, uriHistorial string) *Gaby {
    return &Gaby{
        Envio: false,
        Origen: "10.10.20.25", 
        UsuarioOrigen: "EATON",
        PasswordOrigen: "admin",
        Acceso: acceso,
        DirectorioTemporal: directorioTemporal,
        UriHistorial: uriHistorial,
    }
}

func (cfg *Gaby) Run() error {
    UrlOrigenDatos := fmt.Sprintf("http://%s:%s@%s", cfg.UsuarioOrigen, cfg.PasswordOrigen, cfg.Origen)

    ficheroMarcaTiempo := fmt.Sprintf("%s/%s", cfg.DirectorioTemporal, cfg.Origen)

    // Obtenemos los datos en un estructura compleja
    datos, err := utils.ObtenerDatos(UrlOrigenDatos, cfg.UriHistorial)
    if err != nil {
        return err
    }
    fmt.Println(datos.String())

    /* Esto se correponde con el envio de datos a influxDB*/
    if cfg.Envio {

        horaAnterior, errHora := utils.LeerFechaEnArchivo(ficheroMarcaTiempo)
        if errHora != nil {
            return errHora
        }

        // TODO: Creo que esto no esta funcionando como se supone que debe hacerlo 
        for i := 0; i < 10; i++ {

            if datos.Timestamp.UnixMilli() > horaAnterior {
                err = utils.EnviarDatos(cfg.Acceso, datos, cfg.Origen)
                if err != nil {
                    return err
                }
                fmt.Println(fmt.Sprintf("Enviado %+v", datos))
                // TODO: Trabajar en como debe manejarse este error
                err = utils.GuardarFechaEnArchivo(ficheroMarcaTiempo, datos.Timestamp)
                return err
            } else {
                fmt.Println("Todavía no se puede guardar")
                time.Sleep(1 * time.Second)
                // Obtenemos los datos de nuevo 
                datos, err = utils.ObtenerDatos(UrlOrigenDatos, cfg.UriHistorial)
                if err != nil {
                    return err
                }
            }
        }
    }

    return nil
}

//Tenés la fecha del sistema y del nodo. La fecha del nodo es la que vamos a guardar para comparar
// pero la que vamos a enviar al sistema es la del sistema, que debería ser más actual
// Ambas tienen que ser redondeadas

func main() {
    // Asumimos que si bien pueden haber múltiples origenes de datos, solo hay un destino
    Token := os.Getenv("GABY_TOKEN")
    Bucket := os.Getenv("GABY_BUCKET")
    Endpoint := os.Getenv("GABY_ENDPOINT")
    Organization := os.Getenv("GABY_ORGANIZACION")
    acceso := utils.Acceso{Endpoint, Token, Organization, Bucket}
    
    // TODO: Supongo que también podrían parametrizarse 
    DIRECTORIO_TEMPORAL := "/var/lib/gaby"
    
    // Esto es constante respecto al funcionamiento del sensor 
    URI_HISTORIAL := "PageHislog.html"

    // La operación, propiamente dicha
    err := commandeer.Run(NewGaby(acceso, DIRECTORIO_TEMPORAL, URI_HISTORIAL))
    if err !=nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
