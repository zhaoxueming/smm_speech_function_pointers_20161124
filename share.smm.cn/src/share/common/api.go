package trade

import(
    "github.com/levigross/grequests"
    "errors"
    "net/http"
    "strconv"
    "strings"
    "log"
)

type api_Response interface  {
    Get_code () string
    Get_msg  () string
    Get_data () interface{}
    Set_data ( interface{} )
}

type api_sc_response struct {
    Code        string
    Msg         string
    Data        interface{}
}

func (res *api_sc_response) Get_code() string{
    return res.Code
}
func (res *api_sc_response) Get_msg () string{
    return res.Msg
}
func (res *api_sc_response) Get_data() interface{} {
    return res.Data
}
func (res *api_sc_response) Set_data( data interface{} ){
    res.Data = data
}

type api_ic_response struct {
    Code        int
    Msg         string
    Data        interface{}
}

func (res *api_ic_response) Get_code() string{
    return strconv.Itoa( res.Code )
}
func (res *api_ic_response) Get_msg () string{
    return res.Msg
}
func (res *api_ic_response) Get_data() interface{} {
    return res.Data
}
func (res *api_ic_response) Set_data( data interface{} ){
    res.Data = data
}

type api_Option struct{
    Url             string
    Request         *grequests.RequestOptions
    Response        interface{}
    Method          string
    Host_type       int
    CodeErrorIgnore bool
    Code            string
}

const(
    API_SC      = iota + 1
    API_IC
    API_GET     = "get"
    API_POST    = "post"
    API_EMPTY   = ""
    API_SUCCESS = "0"
    API_UNKNOWN_METHOD  = "unknown request method : "
)

func (opt *api_Option) method() func( string , *grequests.RequestOptions )(*grequests.Response, error) {

    if(opt.Method == API_EMPTY){
        opt.Method = API_POST
    }
    switch(opt.Method){
    case API_GET :
        return grequests.Get
    case API_POST :
        return grequests.Post
    default :
        panic( API_UNKNOWN_METHOD + opt.Method)
    }
}

func ( opt *api_Option )response_obj() api_Response {
    switch opt.Host_type {
        case API_SC :
            return &api_sc_response{}
        case API_IC :
            return &api_ic_response{}
        default :
            return &api_ic_response{}
    }
}

func ( opt *api_Option ) Get() error{

    if(opt.Request == nil){
        opt.Request = &grequests.RequestOptions{}
    }

    if(opt.Request.Data == nil){
        opt.Request.Data = opt.Request.Params
    }else if(opt.Request.Params == nil){
        opt.Request.Params = opt.Request.Data
    }

    if(strings.Index(opt.Url , "https") == 0){
        opt.Request.InsecureSkipVerify = true
    }

    result, err := opt.method()(opt.Url, opt.Request)
    if err != nil {
        log.Println(err,opt)
        return errors.New( ERROR_GOREQUEST )
    }

    if !result.Ok {
        log.Println(result.StatusCode, http.StatusText(result.StatusCode), opt)
        return errors.New( ERROR_SERVICE_REQUEST_CODE )
    }

    response := opt.response_obj()
    if(opt.Response != nil){
        response.Set_data ( opt.Response )
    }

    err = result.JSON(response)
    if err != nil {
        log.Println(err ,result, opt)
        return errors.New( ERROR_SERVICE_REQUEST_DATA )
    }
    opt.Code = response.Get_code()
    if opt.Code != API_SUCCESS && !opt.CodeErrorIgnore {
        return errors.New(response.Get_msg())
    }
    return nil
}
