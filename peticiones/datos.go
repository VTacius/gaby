package peticiones

import (
    "fmt"
    "errors"
    "net/http"
    "golang.org/x/net/html"
) 


func TomarDatosUltimaFila(contenido *http.Response) ([]string) {
    cuerpo := html.NewTokenizer(contenido.Body)
    var datos []string
    for {
        token := cuerpo.Next()
        switch {
            case token == html.ErrorToken:
                return datos 
            case token == html.StartTagToken:
                tag := cuerpo.Token()
                switch tag.Data {
                    case "tr":
                        datos = nil
                    case "center":
                        // Nos movemos hasta el contenido del center, no lo ví a la primera
                        cuerpo.Next()
                        tag = cuerpo.Token()
                        datos = append(datos, tag.Data)
               
                }
        } 
    } 
}

func ObtenerDatosAmbientales(endpoint string, uri string) ([]string, error){
    url := fmt.Sprintf("%s/%s", endpoint, uri)
    respuesta, err := http.Get(url)
    
    if err != nil {
        return []string{}, err
    }

    defer respuesta.Body.Close()
    
    if respuesta.StatusCode != 200 {
        mensaje := fmt.Sprintf("%s devuelve %d: %s", endpoint, respuesta.StatusCode, http.StatusText(respuesta.StatusCode))
        return []string{}, errors.New(mensaje)
    }

    resultado := TomarDatosUltimaFila(respuesta)
    if len(resultado) == 0 {
        mensaje := fmt.Sprintf("No se encontró el contenido adecuado en la página %s. \nRevise que se esté apuntando al dispositivo adecuado", url)
        return []string{}, errors.New(mensaje)
    }
    
    return resultado, nil
}
