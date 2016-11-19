package trade

import(
    "github.com/gin-gonic/gin"
    "github.com/flosch/pongo2"
    "github.com/axgle/mahonia"
    "github.com/smmit/smmbase/qiniu"
    // "github.com/levigross/grequests"
    "crypto/md5"
    "log"
    "strings"
    // "strconv"
    // "math"
    "path/filepath"
    "errors"
    "net/http"
    "io/ioutil"
    "bufio"
    "fmt"
    // "time"
    // "strings"

)

const (
    TRADE_CENTER        = "trade_center"
    PAY_CENTER          = "pay_center"
    MSG_CENTER          = "msg_center"
    USER_CENTER         = "user_center"
    BCHAINTRADE         = "bchaintrade"

    MSG_TRADE_CODE      = 2

    JSON_SUCCESS_CODE   = "0"
    JSON_SUCCESS_MSG    = "ok"
    JSON_ERROR_CODE     = "100"
    JSON_MSG            = "msg"
    JSON_CODE           = "code"
    JSON_DATA           = "data"

    PAY_AUTH_PASS       = "100"
    BILL_CREATE_ORDER_TYPE_SELL = "1"
    BILL_CREATE_ORDER_TYPE_BUY  = "0"
    BCHAIN_AUTH_KEY     = "auth_token"
)

var(
    enc = mahonia.NewEncoder("gbk")
)

func index (c *gin.Context) {
    base := Base{}.Parse(c)
    c.HTML(http.StatusOK, "index.html",pongo2.Context{
        "base"  : base,
    })
}

/*
    /page/:name
*/
func pages(c *gin.Context) {
    base := Base{}.Parse(c)
    req := map[string]string{}
    err := render_param([]string{
        "name",
    },c,req)
    if err != nil {
        page_not_found(c)
        return
    }
    c.HTML(http.StatusOK, req["name"]+ ".html" ,pongo2.Context{
        "base"  : base,
    })
}

func internal_server_error(c *gin.Context) {
    c.HTML(http.StatusInternalServerError, "500.html",pongo2.Context{})
}

func page_not_found(c *gin.Context) {
    c.HTML(http.StatusNotFound, "404.html",pongo2.Context{})
}

func robots(c *gin.Context){
    text := strings.Join( config.Robots ,"\r\n")
    c.String(http.StatusOK, text)
}

func fileupload(c *gin.Context) {
    var err error
    var backurl string
    req := map[string]string{}

    _err := func () {
        c.JSON(http.StatusOK, map[string]interface{}{
            JSON_CODE  : JSON_ERROR_CODE,
            JSON_MSG   : err.Error(),
        })
    }

    _success:= func () {
        c.Header("Content-Type", "text/plain")
        c.Header("Access-Control-Allow-Origin","*")
        c.JSON(http.StatusOK, map[string]interface{}{
            JSON_CODE  : JSON_SUCCESS_CODE,
            JSON_MSG   : JSON_SUCCESS_MSG,
            JSON_DATA  : backurl,
        })
    }

    err = render_param([]string{
        "name",
    },c,req)

    if err != nil {
        _err()
        return
    }

    file, header , err := c.Request.FormFile(req["name"])
    if err != nil {
        log.Println(err)
        err = errors.New(ERROR_FILE_UPLOAD_NAME)
        _err()
        return
    }
    filename := header.Filename
    extension := filepath.Ext(filename)

    if len(extension) == 0 {
        extension = ".jpg"
    }

    blob, err := ioutil.ReadAll(bufio.NewReader(file))
    if err != nil {
        log.Println(err)
        err = errors.New(ERROR_FILE_UPLOAD_BYTES)
        _err()
        return
    }

    key := config.Qiniu_file_upload_path + fmt.Sprintf("%x",md5.Sum(blob)) + extension

    backurl, err = qiniu.QiniuUploadImage(key, blob)
    if err!=nil {
        log.Println(err)
        err = errors.New(ERROR_FILE_UPLOAD_TO_QINIU)
        _err()
        return
    }
    _success()

}

