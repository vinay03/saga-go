package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	log "github.com/sirupsen/logrus"

	saga "github.com/vinay03/saga-go"
)

func main() {
	OrderCreatedSaga, err := saga.NewSaga("OrderCreated").Transactions(
		"CheckProducts",
		"CheckDiscounts",
		"NotifyNewOrderToSeller",
		"NotifyOrderUpdateToBuyer",
	)
	if err != nil {
		log.Fatal(err)
	}

	OrderCreatedSaga.DefineActions(
		"CheckProducts",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered `CheckProducts` phase")
			return data, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered compensating `CheckProducts` phase")
			return data, nil
		},
	)

	OrderCreatedSaga.DefineActions(
		"CheckDiscounts",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered `CheckDiscounts` phase")
			return data, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered compensating `CheckDiscounts` phase")
			return data, nil
		},
	)
	OrderCreatedSaga.DefineActions(
		"NotifyNewOrderToSeller",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered `NotifyNewOrderToSeller` phase")
			return data, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered compensating `NotifyNewOrderToSeller` phase")
			return data, nil
		},
	)
	OrderCreatedSaga.DefineActions(
		"NotifyOrderUpdateToBuyer",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered `NotifyOrderUpdateToBuyer` phase")
			return data, errors.New("Aborted")
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered compensating `NotifyOrderUpdateToBuyer` phase")
			return data, nil
		},
	)
	loggerEntry := log.WithFields(log.Fields{
		"service": "OrderService",
	})
	coord := saga.GetCoordinatorInstance(loggerEntry)
	coord.SetupCarriers(
		&saga.InMemoryCarrierConfig{},
	)
	coord.RegisterSaga(OrderCreatedSaga, coord.Carrier.InMem)

	data := struct {
		Testdata string `json:"testdata"`
	}{
		"Test Sample string",
	}

	coord.Start("OrderCreated", data)

	time.Sleep(1 * time.Second)
}
