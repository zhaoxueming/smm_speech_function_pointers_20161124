#!/bin/bash

function share_pid(){
    echo `ps -ef | grep -i share.out | grep -v "grep" | awk '{print $2}'`
}

function share_gopath(){
    cd ../..
    local_go_path=`pwd | grep -i share | awk '{print $0}'`
    export GOPATH=$GOPATH:${local_go_path}
    cd ./src/share
}

function share_close(){
    sharepid=`share_pid`
    if [ x$sharepid != x""  ];then
        kill -9 ${sharepid}
        echo kill share : ${sharepid}
    fi
}

function share_build(){
    rm -f share.out
    go build -o share.out main.go

    if [ $? != "0" ];then
        exit $?
    fi
    chmod 755 share.out
}

function share_start(){
    nohup ./share.out > share.log 2>&1  &
}

function share_log(){
    tail -f share.log
}

function share_init(){
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

function share_clean(){
    rm -f  share_*.zip
#    rm -rf bower_components
    rm -f  share.log
    rm -f  share.out
#    rm -f  config.user.json
}

function share_package(){
    if [ $# != 1 ];then
        echo "package failed : version is empty"
        echo "use $0 package as"
        echo "======================="
        echo "$0 package v0.0.0"
        echo "======================="
        exit 1
    fi
    if [ -e share_$1.zip ];then
        echo package warning : share_$1.zip is exist
        echo press y to remove it and continue, or other word to exit
        read removesharezip
        if [ x$removesharezip != x"y" ];then
            exit 2
        fi
        rm -rf share_$1.zip
    fi

    if [ -e share_smm_cn ];then
        echo package warning : dir/file share_smm_cn is exist
        echo press y to remove it and continue, or other word to exit
        read removesharesmmcn
        if [ x$removesharesmmcn != x"y"  ];then
            exit 3
        fi
        rm -rf share_smm_cn
    fi

    echo share $1 build start

    export GIN_MODE=release
    mkdir share_smm_cn
    cd share_smm_cn
    mkdir service
    mkdir static
    mkdir library
    cd ../

    echo init success
    share_gopath
    echo gopath set success
    go build -o share_smm_cn/service/share_$1 main.go

    echo go build success

    #cp STAR_smm_cn.crt  share_smm_cn/service/
    #cp smmprivate.key  share_smm_cn/service/
    cp config.default.json  share_smm_cn/service/
    cp -rf templates share_smm_cn/service/templates_$1

    echo service copy success

    cp -rf static/version share_smm_cn/static/$1

    echo static copy success

    bower update -q

    echo bower success

    cp bower.json share_smm_cn/library/
    cp -rf bower_components share_smm_cn/library/

    echo library copy success

    zip -rq share_$1.zip ./share_smm_cn

    echo package success

    rm -rf share_smm_cn

    echo clean success

    echo share $1 build finished
}

function share_help(){
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

function share_run(){
    share_gopath
    share_close
    share_build
    share_start
    share_log
}

case $1 in

    "log")          share_log
    ;;
    "close")        share_close
    ;;
    "run")          share_run
    ;;
    "clean")        share_clean
    ;;
    "package")      share_package $2
    ;;
    "help")         share_help
    ;;
    "init")         share_init
    ;;
    *)
        if [ x$1 != x ];then
            echo unknown command : $1
        fi
        share_help
    ;;
esac
