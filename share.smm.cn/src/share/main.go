package main

import(
    "share/common"
)

const (

    default_path     = "./config.default.json"

    user_path        = "./config.user.json"
)

func main() {
    config_file := []string{ default_path, user_path }
    share.Service(config_file);

}


