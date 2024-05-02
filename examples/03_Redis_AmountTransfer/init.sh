#! /bin/bash

rm ./bin/* >> /dev/null

go build -o ./bin/order order/main.go
go build -o ./bin/account account/main.go
go build -o ./bin/notification notification/main.go

sudo docker compose up -d


#  curl -H 'Content-Type: application/json' -d '{ "from":1,"to":2, "amount": 50}' -X POST http://localhost:8080/transfer