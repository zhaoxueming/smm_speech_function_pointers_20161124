#!/bin/bash

function speech_pid(){
    echo `ps -ef | grep -i speech.out | grep -v "grep" | awk '{print $2}'`
}

function speech_gopath(){
    cd ../..
    local_go_path=`pwd | grep -i speech | awk '{print $0}'`
    export GOPATH=$GOPATH:${local_go_path}
    cd ./src/speech
}

function speech_close(){
    speechpid=`speech_pid`
    if [ x$speechpid != x""  ];then
        kill -9 ${speechpid}
        echo kill speech : ${speechpid}
    fi
}

function speech_build(){
    rm -f speech.out
    go build -o speech.out main.go

    if [ $? != "0" ];then
        exit $?
    fi
    chmod 755 speech.out
}

function speech_start(){
    nohup ./speech.out > speech.log 2>&1  &
}

function speech_log(){
    tail -f speech.log
}

function speech_init(){
    bower update -q
    if [ $? != "0" ];then
        echo "init failed : bower failed"
        echo "Use 'bower update' to download the remote resources."
        exit $?
    fi
    if [ ! -e config.user.json ];then
        cp -f config.default.json config.user.json
    fi
}

function speech_clean(){
    rm -f  speech_*.zip
#    rm -rf bower_components
    rm -f  speech.log
    rm -f  speech.out
#    rm -f  config.user.json
}

function speech_package(){
    if [ $# != 1 ];then
        echo "package failed : version is empty"
        echo "use $0 package as"
        echo "======================="
        echo "$0 package v0.0.0"
        echo "======================="
        exit 1
    fi
    if [ -e speech_$1.zip ];then
        echo package warning : speech_$1.zip is exist
        echo press y to remove it and continue, or other word to exit
        read removespeechzip
        if [ x$removespeechzip != x"y" ];then
            exit 2
        fi
        rm -rf speech_$1.zip
    fi

    if [ -e speech_smm_cn ];then
        echo package warning : dir/file speech_smm_cn is exist
        echo press y to remove it and continue, or other word to exit
        read removespeechsmmcn
        if [ x$removespeechsmmcn != x"y"  ];then
            exit 3
        fi
        rm -rf speech_smm_cn
    fi

    echo speech $1 build start

    export GIN_MODE=release
    mkdir speech_smm_cn
    cd speech_smm_cn
    mkdir service
    mkdir static
    mkdir library
    cd ../

    echo init success
    speech_gopath
    echo gopath set success
    go build -o speech_smm_cn/service/speech_$1 main.go

    echo go build success

    #cp STAR_smm_cn.crt  speech_smm_cn/service/
    #cp smmprivate.key  speech_smm_cn/service/
    cp config.default.json  speech_smm_cn/service/
    cp -rf templates speech_smm_cn/service/templates_$1

    echo service copy success

    cp -rf static/version speech_smm_cn/static/$1

    echo static copy success

    bower update -q

    echo bower success

    cp bower.json speech_smm_cn/library/
    cp -rf bower_components speech_smm_cn/library/

    echo library copy success

    zip -rq speech_$1.zip ./speech_smm_cn

    echo package success

    rm -rf speech_smm_cn

    echo clean success

    echo speech $1 build finished
}

function speech_help(){
    echo "usage: $0 [<command>] [<args>]"
    echo ""
    echo "These are common $0 commands used in various situations:"
    echo "      run         Close the old instance (if exist), build and run a new one."
    echo "      close       Close the running instance (if exist)."
    echo "      log         Show logs."
    echo "      clean       Clean log/binary/package files"
    echo "      package     Package the project as a zip file."
    echo "      help        Show help information."
    echo "      init        Download remote resources and init user config (if need)."
    echo ""
    echo "See '$0 help' to read about this infomation."
}

function speech_run(){
    speech_gopath
    speech_close
    speech_build
    speech_start
    speech_log
}

case $1 in

    "log")          speech_log
    ;;
    "close")        speech_close
    ;;
    "run")          speech_run
    ;;
    "clean")        speech_clean
    ;;
    "package")      speech_package $2
    ;;
    "help")         speech_help
    ;;
    "init")         speech_init
    ;;
    *)
        if [ x$1 != x ];then
            echo unknown command : $1
        fi
        speech_help
    ;;
esac
