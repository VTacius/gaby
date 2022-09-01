package main

import (
	"fmt"
	"log"
	"os"
	"sanidad/alortiz/gaby/utils"
	"time"

	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"
)

func principal(cCtx *cli.Context) error {

	directorio := cCtx.String("directorio")
	origen := cCtx.String("origen")
	origenUsuario := cCtx.String("origen-usuario")
	origenPassword := cCtx.String("origen-password")
	origenIndex := cCtx.String("origen-index")
	envio := cCtx.Bool("envio")
	destino := cCtx.String("destino")
	destinoToken := cCtx.String("destino-token")
	destinoOrganizacion := cCtx.String("destino-organizacion")
	destinoBucket := cCtx.String("destino-bucket")

	// Acá empiezan las operaciones propiamente dichas
	estampa := utils.NewMarcadorTemporal(directorio, origen, "01/02/2006 15:04:05")
	fuente := utils.NewOrigen(origen, origenUsuario, origenPassword)
	datos, err := fuente.ObtenerDatos(origenIndex, *estampa)
	if err != nil {
		return err
	}

	if envio {
		backend := utils.NewBackend(destino, destinoToken, destinoOrganizacion, destinoBucket)

		for i := 0; i < 10; i++ {
			if datos.EsNuevo() {
				backend.Enviar(*datos)
				datos.Timestamp.GuardarActual()
				fmt.Printf("Enviados %s\n", datos.Mensaje())
				break
			} else {
				time.Sleep(1 * time.Second)
				fmt.Println("Intento de envío")
				datos, _ = fuente.ObtenerDatos(origenIndex, *estampa)
			}

		}
	} else {
		fmt.Println(datos)
	}

	return nil
}

func main() {
	directorio := &cli.StringFlag{Name: "directorio", Usage: "Directorio para guardar la fecha del último envío", Value: "/var/lib/gaby"}
	origen := &cli.StringFlag{Name: "origen", Usage: "IP del sensor a scrappear", Required: true}
	origenUsuario := &cli.StringFlag{Name: "origen-usuario", Usage: "Usuario para conectarse al sensor", Value: "EATON"}
	origenPassword := &cli.StringFlag{Name: "origen-password", Usage: "Password para conectarse al sensor", Value: "admin"}
	origenIndex := &cli.StringFlag{Name: "origen-index", Usage: "Página de inicio para el sistema", Value: "PageHislog.html"}
	envio := &cli.BoolFlag{Name: "envio", Usage: "Si debe o no enviarse la información al backend"}
	destino := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino", Usage: "Backend InfluxDB2 para almacenar datos"})
	destinoToken := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino-token", Usage: "Token para conectarse a Backend"})
	destinoOrganizacion := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino-organizacion", Usage: "Organización en Backend"})
	destinoBucket := altsrc.NewStringFlag(&cli.StringFlag{Name: "destino-bucket", Usage: "Bucket dentro de la organización"})

	banderas := []cli.Flag{directorio, origen, origenUsuario, origenPassword, origenIndex, envio, destino, destinoToken, destinoOrganizacion, destinoBucket}

	app := &cli.App{
		Name:  "gaby",
		Usage: "Scrapper para sensores EATON",
		Flags: banderas,
		Before: altsrc.InitInputSourceWithContext(banderas,
			func(context *cli.Context) (altsrc.InputSourceContext, error) {
				return altsrc.NewYamlSourceFromFile("/etc/gaby.yaml")
			}),
		Action: principal,
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
