package trade

import(
    "io/ioutil"
    "encoding/json"
)

type Config_data struct  {
    Version                 string
    Port                    string
    Static_service_path     map[string]string
    Static_service          bool
    Source                  map[string]([]string)
    Links                   map[string]string
    Cookie_domain           string
    Token_key               string
    Release                 bool
    Endless                 bool
    Https                   bool
    Https_service           bool
    Statistics              bool
    Internal_server_error   bool
    Https_cert              string
    Https_key               string
    Token_secret            string
    Qiniu_file_upload_path  string
    Robots                  []string
    Hidden_cookie           []string

}

var (
    config Config_data
)

func config_load(config_obj interface{}, file_paths []string) {

    for _ , file_path := range file_paths {

        file, err := ioutil.ReadFile(file_path)
        if( err != nil ){ continue }

        err = json.Unmarshal(file , &config_obj)
        if( err != nil ){ panic(err) }
    }
}

func config_init(file_paths []string) {
    config_load( &config , file_paths );
}



