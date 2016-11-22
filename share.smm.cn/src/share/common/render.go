package share

import(
    "fmt"
    "time"
    "strings"
    "log"
    "strconv"
    // "net/http"
    "encoding/base64"
    "errors"
    "github.com/gin-gonic/gin"

)

const(
    RENDER_USER_HTTP_PROTOCOL      = "http://"
    RENDER_USER_HTTPS_PROTOCOL     = "https://"

    dateformat                  = "2006-01-02 15:04:05"
    dateformat_onlytime         = "15:04:05"
    dateformat_onlydate         = "2006-01-02"
    dateformat_seconds_per_day  = 24 * 60 * 60
    numberformat_radix_point    = "."
    numberformat_decollator     = ","
    numberformat_decollator_len = 3
    numberformat_empty          = "0"
    render_param_type_num       = 4
    render_map_base64           = "base64"
)



type Render struct{
    Counter int
    Base64_encode   func(string)(string)
    Base64_decode   func(string)(string,error)
    Source          func(string)(string)
    Links           func(string)(string)
    Param_base64            func([]string,*gin.Context,map[string]string)(error)
    Param_base64_ignore     func([]string,*gin.Context,map[string]string)(error)
    Param                   func([]string,*gin.Context,map[string]string)(error)
    Param_ignore            func([]string,*gin.Context,map[string]string)(error)
    Money_format    func(float64)(string)
    Quantity_format func(float64)(string)
    Datetime_format func(int64)(string)
    Date_format     func(int64)(string)
    Time_format     func(int64)(string)
    Map_get         func(map[string]string,string)(string)
    Array_get       func([]interface{},int)(interface{})
    To_string       func(int)string
    Connect_string  func(...string)string
}

var (
    render Render
)

func init(){
    render = Render{
        Counter         : 1,
        Base64_encode   : render_base64_encode,
        Base64_decode   : render_base64_decode,
        Source          : render_source,
        Links           : render_links,
        Param_base64    : render_param_base64,
        Param_base64_ignore  : render_param_base64_ignore,
        Param           : render_param,
        Param_ignore    : render_param_ignore,
        Money_format    : render_money_format,
        Quantity_format : render_quantity_format,
        Datetime_format : render_datetime_format,
        Date_format     : render_date_format,
        Time_format     : render_time_format,
        Map_get         : render_map_get,
        Array_get       : render_array_get,
        To_string       : render_to_string,
        Connect_string  : render_connect_string,
    }
}

func render_source( key string ) string {
    keys := config.Source[key]
    render.Counter++
    return keys[ render.Counter % len(keys) ]
}

func render_links( key string) string {
    return config.Links[key]
}

func render_base64_encode( key string ) string {

    return strings.Replace(base64.StdEncoding.EncodeToString([]byte(key)),"/","*", -1)
}

func render_base64_decode( key string ) (string, error) {

    value ,err := base64.StdEncoding.DecodeString(strings.Replace(key,"*","/",-1))
    return string(value) ,err
}

func render_param_base64( keys []string, c *gin.Context , req map[string]string) error {
    for i := 0; i < len(keys); i++ {

        param, ok := c.Params.Get(keys[i])
        if !ok {
            log.Println( keys[i], req , c)
            return errors.New( ERROR_PARAM_FORMAT )
        }
        if param == "_" || param == "" {
            continue
        }
        data, err := render_base64_decode(param)
        if err != nil {
            log.Println( err , keys[i], req , c )
            return errors.New( ERROR_PARAM_FORMAT )
        }
        req[keys[i]] = data
    }
    return nil
}

func render_param_base64_ignore( keys []string, c *gin.Context , req map[string]string) error {
    for i := 0; i < len(keys); i++ {

        param, ok := c.Params.Get(keys[i])
        if !ok || param == "_" || param == "" {
            continue
        }
        data, err := render_base64_decode(param)
        if err != nil {
            log.Println( err , keys[i], req , c )
            return errors.New( ERROR_PARAM_FORMAT )
        }
        req[keys[i]] = data
    }
    return nil
}


func render_param( keys []string, c *gin.Context , req map[string]string) error {
    for i := 0; i < len(keys); i++ {

        param, ok := c.Params.Get(keys[i])
        if !ok {
            log.Println( keys[i], req , c)
            return errors.New( ERROR_PARAM_FORMAT )
        }
        if param == "_" {
            continue
        }
        req[keys[i]] = param
    }
    return nil
}

func render_param_ignore( keys []string, c *gin.Context , req map[string]string) error {
    for i := 0; i < len(keys); i++ {

        param, ok := c.Params.Get(keys[i])
        if !ok || param == "_" || param == "" {
            continue
        }
        req[keys[i]] = param
    }
    return nil
}

func render_param_set( c *gin.Context,key,value string) {
    for index, _ := range c.Params {
        if c.Params[index].Key == key {
            c.Params[index].Value = value
            return
        }
    }
    c.Params = append(c.Params,gin.Param{
        Key : key,
        Value : value,
    })
}

func render_param_all(  p []string ,
                        pb []string ,
                        pi []string ,
                        pbi []string ,
                        c *gin.Context ,
                        req map[string]string  ) error {
    var err error
    keys := [render_param_type_num]([]string){
        p,
        pb,
        pi,
        pbi,
    }
    param := [render_param_type_num](func([]string,*gin.Context,map[string]string)error){
        render_param,
        render_param_base64,
        render_param_ignore,
        render_param_base64_ignore,
    }
    for i := 0; i < render_param_type_num; i++ {
        if keys[i] != nil {
            err = param[i](keys[i],c,req)
            if(err != nil){
                return err
            }
        }
    }
    return nil
}

func render_error_control(errfuncs ...func()(error)) error {
    var err error
    var lens = len(errfuncs)
    for i := 0; i < lens; i++ {
        err = errfuncs[i]()
        if err != nil {
            return err
        }
    }
    return nil
}

type Render_map struct {
    Map     map[string]string
    Get     func (ma map[string]string,key string)string
    Value   func (ma map[string]string,key string)string
    Url     func (req *Render_map,keys ...string)string
}

func render_map(ma map[string]string) (*Render_map){
    return &Render_map{
        Map  :   ma,
        Get  :   func (ma map[string]string,key string)string {
            value := ma[key]
            if value == "" {
                return "_"
            }
            return value
        },
        Value :  func (ma map[string]string,key string) string{
            return ma[key]
        },
        Url : func (req *Render_map, keys ...string)string {
            return req.param(keys)
        },
    }
}

func (r *Render_map)get(key string)(string) {
    return r.Get(r.Map,key)
}

func (r *Render_map)value(key string)(string) {
    return r.Value(r.Map,key)
}

func (r *Render_map)param(keys []string)(string) {
    lens := len(keys)
    base64 := false
    values := []string{}
    for i := 0; i < lens; i++ {
        var key = keys[i]
        if(key == render_map_base64){
            base64 = true
            continue
        }
        value := r.get(key)
        if(base64 && value != "_"){
            value = render_base64_encode(value)
        }
        base64 = false
        values = append(values,value)
    }
    return strings.Join(values,"/")
}

func render_string_section(value string,lens int, front bool) []string {

    var vals []string

    bits    := []rune( value )
    bit_len := len( bits )

    tail    := bit_len % lens

    start   := 0
    end     := bit_len

    if tail != 0 {
        if front {
            end = bit_len - tail
        }else{
            vals = append(vals,string(bits[:tail]))
            start = tail
        }
    }

    for i := start ; i < end; i += lens {
        vals = append(vals,string(bits[i:i+lens]))
    }

    if tail != 0 && front {
        vals = append(vals,string(bits[end:]))

    }

    return vals
}

func render_money_format (number float64) string{
    value := strings.Split( fmt.Sprintf("%.2f",number) , numberformat_radix_point )

    if len(value) != 2 {
        panic( fmt.Sprintf("strings.Split error? float : %d" , number))
    }

    value[0] = strings.Join(render_string_section(value[0],numberformat_decollator_len,false),numberformat_decollator)

    return strings.Join( value, numberformat_radix_point )
}

func render_quantity_format (number float64) string{
    value := strings.Split( fmt.Sprintf( "%.5f", number ) , numberformat_radix_point )
    if len(value) != 2 {
        panic( fmt.Sprintf("strings.Split error? float : %d" , number))
    }
    value[1] = strings.TrimRight( value[1] , numberformat_empty )

    value[0] = strings.Join(render_string_section(value[0],numberformat_decollator_len,false),numberformat_decollator)

    if value[1] == "" {
        return value[0]
    }
    return strings.Join( value, numberformat_radix_point )
}

func render_datetime_format(number int64)string {
    if number == 0 {
        return ""
    }
    return time.Unix(number,0).Format(dateformat)
}

func render_date_format(number int64)string {
    if number == 0 {
        return ""
    }
    return time.Unix(number,0).Format(dateformat_onlydate)
}

func render_time_format(number int64)string {
    if number == 0 {
        return ""
    }
    return time.Unix(number,0).Format(dateformat_onlytime)
}

func render_map_get (m map[string]string,key string) string {
    return m[key]
}

func render_array_get (a []interface{},index int) interface{} {
    return a[index]
}

func render_to_string (index int) string{
    return strconv.Itoa(index)
}

func render_connect_string(values ...string)string{
    return strings.Join(values,"")
}
/*
    2016-05-15 to 1432412520
*/
func render_date_to_unix_timestamp(value string) ( int64 , error ) {
    ti , err := time.ParseInLocation(dateformat_onlydate,value,time.Local)
    if  err != nil {
        return 0 , err
    }
    return ti.Unix() , nil
}

func render_start_end_time_to_str_unix_timestamp(start, end string) (string,string,error) {
    var start_str , end_str string
    if start != "" {
        ti, err := render_date_to_unix_timestamp(start)
        if err != nil {
            log.Println(err , start, end )
            return "","",errors.New(ERROR_DATE_FORMAT)
        }
        start_str = fmt.Sprintf("%d", ti)
    }
    if end != "" {
        ti, err := render_date_to_unix_timestamp(end)
        if err != nil {
            log.Println(err , start, end )
            return "","",errors.New(ERROR_DATE_FORMAT)
        }
        end_str = fmt.Sprintf("%d", ti + dateformat_seconds_per_day - 1 )
    }
    return start_str , end_str , nil
}

func render_protocol() string {
    if(config.Https){
        return RENDER_USER_HTTPS_PROTOCOL
    }
    return RENDER_USER_HTTP_PROTOCOL
}

