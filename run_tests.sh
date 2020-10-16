#!/bin/sh
export PATH=$PATH:/usr/local/go/bin
export GOPATH="$WORKSPACE"

#mkdir -p $GOPATH/src/github.com/mercimat
#ln -s $WORKSPACE $GOPATH/src/github.com/mercimat/instavote
cd $GOPATH/src/github.com/mercimat/instavote

if [ ! -f ./go.mod ]; then
    go mod init github.com/mercimat/instavote
fi

go get -d -v go.mongodb.org/mongo-driver/mongo
go get -d -v github.com/go-redis/redis/v8
go get -d -v github.com/google/uuid

cd core/
go test -v -vet=off .
