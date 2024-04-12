package main

import (
	"errors"

	saga "github.com/vinay03/saga-go"
)

func main() {
	orderSaga := saga.NewSaga("Order_Created")
	orderSaga.AddStages(
		&saga.Stage{
			ID: "Step-1",
			Action: func() error {
				return errors.New("")
			},
			CompensateAction: func() error {
				return errors.New("")
			},
		},
	)
	orderSaga.Start()
}
