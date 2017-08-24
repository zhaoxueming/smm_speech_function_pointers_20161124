package share

import(
    "github.com/gin-gonic/gin"
    "github.com/flosch/pongo2"
    "github.com/axgle/mahonia"
    "github.com/smmit/smmbase/qiniu"
    "gitHub.com/zhaoxueming/smm_speech_function_pointers_20161124/share.smm.cn/src/share/common/gomoku"
    // "github.com/levigross/grequests"
    "crypto/md5"
    "log"
    "strings"
    "strconv"
    // "math"
    "path/filepath"
    "errors"
    "net/http"
    "io/ioutil"
    "os/exec"
    "bufio"
    "fmt"
    // "time"
    "encoding/json"
    // "strings"

)

const (

    JSON_SUCCESS_CODE   = "0"
    JSON_SUCCESS_MSG    = "ok"
    JSON_ERROR_CODE     = "100"
    JSON_MSG            = "msg"
    JSON_CODE           = "code"
    JSON_DATA           = "data"

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
    "/tool/gorun/:value"
*/
func go_run(c *gin.Context){
    tool_run(c,"./tmp/main.go","go")
}
/*
    "/tool/noderun/:value"
*/
func node_run(c *gin.Context){
    tool_run(c,"./tmp/main.js","node")
}

func tool_run(c *gin.Context,file string, name string,) {
    req := map[string]string{}
    res := ""

    err := render_error_control(func()error{
        return render_param_base64([]string{
            "value",
        },c,req)
    },func ()error {
        return ioutil.WriteFile(file,[]byte( req["value"] ),0666)
    },func ()error {
        out ,err := exec.Command("./tool_run.sh",name).Output()
        if err != nil {
            return err
        }
        res = string(out)
        return nil
    })
    if err != nil {
        api_error(c,err)
        return
    }

    c.Header("Content-Type", "text/plain")
    c.Header("Access-Control-Allow-Origin","*")
    c.JSON(http.StatusOK, map[string]interface{}{
        JSON_CODE  : JSON_SUCCESS_CODE,
        JSON_MSG   : JSON_SUCCESS_MSG,
        JSON_DATA  : res,
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

func api_error(c *gin.Context, err error) {
    c.JSON(http.StatusOK, map[string]interface{}{
            JSON_CODE  : JSON_ERROR_CODE,
            JSON_MSG   : err.Error(),
    })
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

var (
    action_gomoku_style = map[string](func(gomoku.GomokuBoard,int,*gomoku.GomokuPoint) error){
        "free" : gomoku.Gomoku_free,
    }
)

func action_gomoku(c *gin.Context){
    c.Header("Content-Type", "application/json; charset=UTF-8")
    c.Header("Access-Control-Allow-Origin", "*")
    c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept")
    c.Header("Access-Control-Allow-Methods", "HEAD, GET, POST, DELETE, PUT, OPTIONS")
    var err error
    var point gomoku.GomokuPoint
    var board gomoku.GomokuBoard
    var role int
    req := map[string]string{}
    _err := func () {
        c.JSON(http.StatusOK, map[string]interface{}{
            JSON_CODE  : JSON_ERROR_CODE,
            JSON_MSG   : err.Error(),
        })
    }

    _success:= func () {

        c.JSON(http.StatusOK, map[string]interface{}{
            JSON_CODE  : JSON_SUCCESS_CODE,
            JSON_MSG   : JSON_SUCCESS_MSG,
            JSON_DATA  : &point,
        })
    }

    err = render_error_control(
    func()error{
        return render_param([]string{
            "player_role",
            "chessboard",
            "style",
        },c,req)
    },func ()error {
        role,err = strconv.Atoi(req["player_role"])
        return err
    },func ()error {
        return json.Unmarshal([]byte(req["chessboard"]) , &board)
    },func ()error {
        var action = action_gomoku_style[req["style"]]
        if(action == nil){
            return errors.New("unknown style")
        }
        return action(board, role, &point)
    })

    if( err != nil ){
        log.Println(err)
        _err();
        return
    }
    _success();
}
