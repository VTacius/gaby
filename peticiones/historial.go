package peticiones

import (
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/net/html"
)

func obtenerEnlace(token html.Token) (enlace string) {
	for _, a := range token.Attr {
		if a.Key == "href" {
			enlace = a.Val
		}
	}

	return
}

func EncontrarEnlaceListaHistorial(contenido *http.Response) string {
	cuerpo := html.NewTokenizer(contenido.Body)
	for {
		token := cuerpo.Next()
		switch {
		case token == html.ErrorToken:
			return ""
		case token == html.StartTagToken:
			tag := cuerpo.Token()
			if isEnlace := tag.Data == "a"; isEnlace {
				return obtenerEnlace(tag)
			}
		}
	}

}

func ObtenerEnlaceHistorial(endpoint string, uri string) (string, error) {
	url := fmt.Sprintf("%s/%s", endpoint, uri)
	respuesta, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer respuesta.Body.Close()

	if respuesta.StatusCode != 200 {
		mensaje := fmt.Sprintf("%s devuelve %d: %s", endpoint, respuesta.StatusCode, http.StatusText(respuesta.StatusCode))
		return "", errors.New(mensaje)
	}

	resultado := EncontrarEnlaceListaHistorial(respuesta)
	if resultado == "" {
		mensaje := fmt.Sprintf("No se encontró el contenido adecuado en la página %s. \nRevise que se esté apuntando al dispositivo adecuado", url)
		return "", errors.New(mensaje)
	}

	return resultado, nil
}
