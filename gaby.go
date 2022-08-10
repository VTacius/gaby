package main
import (
    "os"
    "fmt"
    "bytes"
    "errors"
    "net/http"
    "encoding/json"
    "github.com/gin-gonic/gin"
)

type AlertaBase struct {
    CheckName string `json:"_check_name"`
    Level string `json:"_level"`
    Host string `json:"host"`
}

func (alerta AlertaBase) String() string{
    return fmt.Sprintf(`
Nombre: %s
Nivel: %s
Hostname: %s
    `, alerta.CheckName, alerta.Level, alerta.Host) 
}

type AlertaDisco struct {
    CheckName string `json:"_check_name"`
    Level string `json:"_level"`
    Host string `json:"host"`
    Path string `json:"path"`
    UsedPercent string `json:"used_percent"`
}

func (alerta AlertaDisco) String() string{
    return fmt.Sprintf(`
Nombre: %s
Nivel: %s
Hostname: %s

Particion: %s    -    Uso: %s
`, alerta.CheckName, alerta.Level, alerta.Host, alerta.Path, alerta.UsedPercent) 
}

// enviarPeticion trata directamente con la API de telegram
// TODO: Tenemos que manejar mejor los mensajes, sobre todo los de error, que la API pudiera enviar
func enviarPeticion(token string, chat_id string, contenido string) (int, error){
    valores := map[string]interface{}{
        "chat_id": chat_id, 
        "text": contenido,
        "disable_notification": true,
        "parse_mode": "markdown",
    }
    
    mensaje, err := json.Marshal(valores)
    
    if err != nil {
        return 500, err  
    }
    
    // TODO: Sacar la construcción del endpoint a otra parte
    endpoint := fmt.Sprintf("https://api.telegram.org/%s/sendMessage", token)
    respuesta, err := http.Post(endpoint, "application/json", bytes.NewBuffer(mensaje))
    
    if err != nil {
        fmt.Println(err)
        return 500, errors.New(fmt.Sprintf("%s", err))
    }
    
    if respuesta.StatusCode != 200 {
        return respuesta.StatusCode, errors.New(fmt.Sprintf("%s", respuesta.Body))
    }

    return respuesta.StatusCode, nil
}

func enviarAlertaDisco(token string, chat_id string) func(*gin.Context) {
    return func (c *gin.Context) {
        // TODO: Cuando existan tipos génericos
        var newAlerta AlertaDisco

        if err := c.BindJSON(&newAlerta); err != nil {
            c.IndentedJSON(500, gin.H{"error": fmt.Sprintf("%s", err)})
            return 
        }

        contenidoAlerta := fmt.Sprintf("%s", newAlerta)
        codigo, err := enviarPeticion(token, chat_id, contenidoAlerta)
        
        if err != nil {
            c.IndentedJSON(codigo, gin.H{"error": fmt.Sprintf("%s", err)})
            return
        }  
        
        c.IndentedJSON(http.StatusCreated, newAlerta)
    }
}

// TODO: Copiar y pegar, para que al menos funcione. Refactorizar para después
func enviarAlertaBase(token string, chat_id string) func(*gin.Context) {
    return func (c *gin.Context) {
        // TODO: Cuando existan tipos génericos
        var newAlerta AlertaBase

        if err := c.BindJSON(&newAlerta); err != nil {
            c.IndentedJSON(500, gin.H{"error": fmt.Sprintf("%s", err)})
            return 
        }
        
        contenidoAlerta := fmt.Sprintf("%s", newAlerta)
        codigo, err := enviarPeticion(token, chat_id, contenidoAlerta)

        if err != nil {
            c.IndentedJSON(codigo, gin.H{"error": fmt.Sprintf("%s", err)})
            return
        }  
        
        c.IndentedJSON(http.StatusCreated, newAlerta)
    }
}

func main() {
    TELEGRAM_BOT_TOKEN := os.Getenv("TELEGRAM_BOT_TOKEN")
    TELEGRAM_CHAT_ID := os.Getenv("TELEGRAM_CHAT_ID")
    KATY_PROXY_IP := os.Getenv("KATY_PROXY_IP")
    KATY_SOCKET := os.Getenv("KATY_SOCKET")
    
    newEnviarAlertaDisco := enviarAlertaDisco(TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID)
    newEnviarAlertaBase := enviarAlertaBase(TELEGRAM_BOT_TOKEN, TELEGRAM_CHAT_ID)
    
    router := gin.Default()
    router.SetTrustedProxies([]string{KATY_PROXY_IP})
    
    router.POST("/alertas", newEnviarAlertaBase)
    router.POST("/alertas/disco", newEnviarAlertaDisco)

    router.Run(KATY_SOCKET)
}
