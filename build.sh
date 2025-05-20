#!/usr/bin/env sh

BIN="cli"
LDFLAGS=""
#LDFLAGS="-s -w"
export GOOS=linux
export GOARCH=amd64
stripped=false

for arg in "$@"; do
    case $arg in
        "stripped")
            stripped=true
            echo here
            LDFLAGS="-s -w"
            ;;
        "deploy")
            export GOARCH="arm64"
            ;;
    esac
done


go build -ldflags "$LDFLAGS" -o $BIN

# strip further if 'upx' command is available
if [ $stripped = true ];then
    if $(which upx>/dev/null);then
        upx --brute $BIN
    fi
fi
