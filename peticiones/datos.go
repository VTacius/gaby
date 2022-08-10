package peticiones

import (
    "fmt"
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

    defer respuesta.Body.Close()

    if err != nil {
        return []string{}, err
    }
    
    resultado := TomarDatosUltimaFila(respuesta)
    
    return resultado, nil

}
