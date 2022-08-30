package main

import (
    "os"
    "fmt"
    "time"
    "sanidad/alortiz/gaby/peticiones"
    "sanidad/alortiz/gaby/utils"
    "github.com/jaffee/commandeer"
)

type Gaby struct {
    Envio bool `help:"Hace el envío de los datos hacia un sistema influxDB"`
    Origen string `help:"URL del sistema EATON"`
    UsuarioOrigen string `help:"Usuario para acceder al sistema EATON"`
    PasswordOrigen string `help:"Password para acceder al sistema EATON"`
}

func NewGaby() *Gaby {
    return &Gaby{
        Envio: false,
        Origen: "10.10.20.25", 
        UsuarioOrigen: "EATON",
        PasswordOrigen: "admin",
    }
}

func (cfg *Gaby) Run() error {
    // Asumimos que si bien pueden haber múltiples origenes de datos, solo hay un destino
    Token := os.Getenv("GABY_TOKEN")
    Bucket := os.Getenv("GABY_BUCKET")
    Endpoint := os.Getenv("GABY_ENDPOINT")
    Organization := os.Getenv("GABY_ORGANIZACION")
    
    // Supongo que también podrían parametrizarse 
    zonaHoraria := "America/El_Salvador"
    ficheroTemporal := "/var/lib/gaby"
    
    // Esta no puede cambiarse, depende es como parseamos la fecha del sistema EATON
    timeLayout := "01/02/2006 15:04:05"
    // Tampoco, porque es parte del sistema
    URI_HISTORIAL := "PageHislog.html"
    
    URL_ORIGEN_DATOS := fmt.Sprintf("http://%s:%s@%s", cfg.UsuarioOrigen, cfg.PasswordOrigen, cfg.Origen)

    // De la página inicial, obtenemos el enlace a los datos 
    enlace, err := peticiones.ObtenerEnlaceHistorial(URL_ORIGEN_DATOS, URI_HISTORIAL)
    if err != nil {
        return err
    }
   
    // Buscamos el último dato en la página de datos
    var resultado []string
    resultado, err = peticiones.ObtenerDatosAmbientales(URL_ORIGEN_DATOS, enlace)
    if err != nil {
        return err
    }
    datos := utils.NewDatos(resultado)
    
    fmt.Println(datos.String())
   
    /* Esto se correponde con el envio de datos a influxDB*/
    if cfg.Envio {
        config := utils.Configuracion{Endpoint, Token, Organization, Bucket}
        
        horaActual := utils.ParsearHora(zonaHoraria, timeLayout, datos) 
        horaAnterior := utils.LeerFechaEnArchivo(ficheroTemporal)
        
        // TODO: Creo que esto no esta funcionando como se supone que debe hacerlo 
        for i := 0; i < 10; i++ {
            if horaActual.UnixMilli() > horaAnterior {
                err = utils.EnviarDatos(config, datos, cfg.Origen)
                if err != nil {
                    return err
                }
                fmt.Println(fmt.Sprintf("Enviado %+v", datos))
                // TODO: Trabajar en como debe manejarse este error
                utils.GuardarFechaEnArchivo(ficheroTemporal, horaActual)
                return nil
            } else {
                fmt.Println("Todavía no se puede guardar")
                time.Sleep(1 * time.Second)
            }
        }
    }

    return nil
}

//Tenés la fecha del sistema y del nodo. La fecha del nodo es la que vamos a guardar para comparar
// pero la que vamos a enviar al sistema es la del sistema, que debería ser más actual
// Ambas tienen que ser redondeadas

func main() {
    err := commandeer.Run(NewGaby())
    if err !=nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
