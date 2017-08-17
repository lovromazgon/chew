#!/bin/bash

#version=`git describe --tags $(git rev-list --tags --max-count=1)`
version=$1
curdate=`date '+%F'`

go build -o chew-${version} -ldflags \
"-X github.com/lovromazgon/chew.Version=${version}
 -X github.com/lovromazgon/chew.ReleaseDate=${curdate}" chew/main.go
