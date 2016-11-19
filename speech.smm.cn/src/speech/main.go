package main

import(
    "speech/common"
)

const (

    default_path     = "./config.default.json"

    user_path        = "./config.user.json"
)

func main() {
    config_file := []string{ default_path, user_path }
    trade.Service(config_file);

}


