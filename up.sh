#!/bin/sh

cd go_runtime
go get
go mod vendor

cd ..
docker-compose up -d --build
docker-compose logs -f