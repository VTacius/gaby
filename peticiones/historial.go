package peticiones

import (
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

func EncontrarEnlaceListaHistorial(contenido *http.Response) (string) {
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
    
    return ""
}
   
func ObtenerEnlaceHistorial(endpoint string, uri string) (string, error){
    url := fmt.Sprintf("%s/%s", endpoint, uri)
    respuesta, err := http.Get(url)

    if err != nil {
        fmt.Println(err)
        return "", err
    }
    
    defer respuesta.Body.Close()

    resultado := EncontrarEnlaceListaHistorial(respuesta)
    
    return resultado, nil

}
