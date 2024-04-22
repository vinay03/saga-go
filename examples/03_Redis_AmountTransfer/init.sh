#! /bin/bash

rm ./bin/* >> /dev/null

go build -o ./bin/order order/main.go
go build -o ./bin/account account/main.go
go build -o ./bin/notification notification/main.go

sudo docker compose up -d