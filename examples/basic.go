package main

import (
	"context"
	"errors"

	saga "github.com/vinay03/saga-go"
)

func main() {
	orderSaga := saga.NewSaga("Order_Created")
	orderSaga.AddStage(&saga.Stage{
		ID: "check_products",
		Action: func(ctx context.Context, data interface{}) error {
			return errors.New("")
		},
		CompensateAction: func(ctx context.Context, data interface{}) error {
			return errors.New("")
		},
	})
	orderSaga.Start()
}
