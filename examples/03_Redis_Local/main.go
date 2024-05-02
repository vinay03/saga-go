package main

import (
	"context"

	log "github.com/sirupsen/logrus"

	"github.com/vinay03/saga-go"
)
var TransferAmountSaga *saga.Saga

func Step1(ctx context.Background, data interface{}) (interface{}, error) {
}


func main() {
	InitSagas()
	sagaCoord := saga.GetCoordinatorInstance()

	payload := []string{}
	sagaCoord.Start("TransferAmountSaga", payload)
}

func InitSagas() {
	var err error
	var TransferAmountSaga *saga.Saga
, err = saga.NewSaga("TransferAmountSaga").Transactions(
		"DeduceSenderBalancer",
		"IncreaseReceiverBalance",
		"NotifySenderWithUpdatedBalance",
		"NotifyReceiverWithUpdatedBalance",
	)
	if err != nil {
		log.Fatal(err)
	}

	coord := saga.GetCoordinatorInstance()
	coord.SetupCarriers(
		&saga.InMemoryCarrierConfig{},
		&saga.RedisCarrierConfig{
			Host:     "redis:6379",
			Password: "admin",
		},
	)
	coord.RegisterSaga(TransferAmountSaga, coord.Carrier.Redis)
}
