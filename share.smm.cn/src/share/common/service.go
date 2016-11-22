package share

import (

    "github.com/gin-gonic/gin"
    "github.com/robvdl/pongo2gin"
    "github.com/fvbock/endless"
)

const (
    METHOD_GET          = iota
    METHOD_POST
    METHOD_NOT_FOUND
    SERVICE_NO_ACTION   = "no action in router when service init!"
    SERVICE_PANIC       = "panic"
)

var (

    router      *gin.Engine
    actions     []ServicAction
    service_internal_server_error func(*gin.Context)
)

func Service( config_files []string ){
    config_init (config_files)
    service_init   ()
    service_static ()
    service_action ()
    service_listen ()

}

func service_init(){
    if config.Release {
        gin.SetMode(gin.ReleaseMode)
    }
    router = gin.Default()
    router.HTMLRender = pongo2gin.Default()
}

func service_static (){
    if( !config.Static_service){
        return
    }
    for url_path , file_path := range config.Static_service_path{

        router.Static( url_path , file_path )

    }
}

type ServicAction struct{
    Method  int
    Auth    func(*gin.Context)bool
    Actions map[string](func(*gin.Context))
}

func service_action (){
    if(actions == nil){
        panic(SERVICE_NO_ACTION)
    }
    for _, sa := range actions {
        for _path , _action := range sa.Actions {
            func(){
                _inner_action := _action
                _inner_sa     := sa
                auth_action := func(c * gin.Context){
                    defer func() {
                        if(config.Internal_server_error && service_internal_server_error != nil){
                            err := recover()
                            if err != nil {
                                c.Set(SERVICE_PANIC,err)
                                service_internal_server_error(c)
                                panic(err)
                            }
                        }
                    }()
                    if _inner_sa.Auth(c) {
                        _inner_action(c)
                    }
                }
                switch sa.Method {
                case METHOD_GET :
                    router.GET( _path , auth_action )
                case METHOD_POST :
                    router.POST( _path , auth_action )
                case METHOD_NOT_FOUND :
                    router.NoRoute( auth_action )
                }
            }()
        }
    }
}

func service_listen (){
    var err error
    if config.Https_service {
        err = endless.ListenAndServeTLS(":" + config.Port, config.Https_cert ,config.Https_key , router)
    }else{
        if config.Endless {
            err = endless.ListenAndServe(":" + config.Port, router);
        }else{
            err = router.Run(":" + config.Port)
        }
    }
    if(err != nil){
        panic(err)
    }
}




