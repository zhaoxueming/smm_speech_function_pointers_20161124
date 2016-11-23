#!/bin/bash

case $1 in

    "go")
        cd tmp
        go run main.go
    ;;
    "node")
        cd tmp
        node main.js
    ;;
esac
