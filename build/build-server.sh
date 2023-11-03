#!/bin/bash
set -e -x

#go env -w GOPROXY=https://mirrors.aliyun.com/goproxy

go env -w GOPROXY="https://proxy.golang.org,https://goproxy.io,direct"

cd $(dirname $0)/../src/cmd/socialserver
#go build -ldflags "-linkmode external -extldflags -static -s"
go build


#go build --ldflags "-extldflags ' -L/usr/lib -ltdjson_static -ltdjson_private -ltdclient -ltdcore -ltdactor -ltddb -ltdsqlite -ltdnet -ltdutils -ldl -lm -lssl -lcrypto -lstdc++ -lz'"
