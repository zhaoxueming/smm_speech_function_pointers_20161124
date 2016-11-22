package share

import (
    "github.com/gin-gonic/gin"
)

var (


    posts                   map[string]func(*gin.Context)
    index_gets              map[string]func(*gin.Context)
    public_gets             map[string]func(*gin.Context)
    not_founds              map[string]func(*gin.Context)

)

func init(){

    public_gets = map[string](func(*gin.Context)){
        "/"                                 : index,
        "/robots.txt"                       : robots,
        "/page/:name"                       : pages,
    }

    posts = map[string](func(*gin.Context)){
        "/fileupload/:name/upload.json"     : fileupload,
    }

    not_founds = map[string](func(*gin.Context)){
        "*" : page_not_found,
    }

    actions = []ServicAction{
        ServicAction{
            Method : METHOD_NOT_FOUND,
            Auth   : auth_no_limit,
            Actions : not_founds,
        },
        ServicAction{
            Method : METHOD_GET,
            Auth   : auth_no_limit,
            Actions : public_gets,
        },
        ServicAction{
            Method : METHOD_POST,
            Auth   : auth_no_limit,
            Actions : posts,
        },
    }
    service_internal_server_error = internal_server_error
}
