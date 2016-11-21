package trade
import(
    // "net/http"
    // "errors"
    // "github.com/flosch/pongo2"
    "github.com/gin-gonic/gin"
    // "github.com/levigross/grequests"
)

type Base struct{
    Frame           Frame
    Render          *Render
    Config          *Config_data
    C               *gin.Context
}

func (b Base)Parse(c *gin.Context)*Base {
    b.Frame.Parse(c)
    b.Render = &render
    b.Config = &config
    b.C = c
    return &b
}

type Frame struct{
    Hiddens    map[string]bool
}

func (f Frame)IsHidden(key string)bool {
    return f.Hiddens[key]
}

func (f *Frame)Parse(c *gin.Context) *Frame {
    f.Hiddens = map[string]bool{}
    for i := 0; i < len(config.Hidden_cookie); i++ {
        index_close, err := c.Cookie(config.Hidden_cookie[i])
        f.Hiddens[config.Hidden_cookie[i]] = ( err == nil && index_close != "" )
    }
    return f
}


func auth_no_limit(c * gin.Context)bool {
    return true
}

