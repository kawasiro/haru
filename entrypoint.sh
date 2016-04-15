#!/bin/bash

git pull origin master

go get github.com/tools/godep
godep restore
go build -v

./haru
