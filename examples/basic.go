package main

import (
	"context"
	"errors"
	"fmt"

	saga "github.com/vinay03/saga-go"
)

func main() {
	orderCreatedSaga := saga.NewSaga("Order_Created")
	orderCreatedSaga.AddStages(
		&saga.Stage{
			ID: "Step-1",
			Action: func(ctx context.Context, data interface{}) (error, interface{}) {
				fmt.Println("1->")
				return nil, nil
			},
			CompensateAction: func(ctx context.Context, data interface{}) (error, interface{}) {
				fmt.Println("1<-")
				return nil, nil
			},
		},
		&saga.Stage{
			ID: "Step-2",
			Action: func(ctx context.Context, data interface{}) (error, interface{}) {
				fmt.Println("2->")
				return errors.New("some"), nil
			},
			CompensateAction: func(ctx context.Context, data interface{}) (error, interface{}) {
				fmt.Println("2<-")
				return nil, nil
			},
		},
		&saga.Stage{
			ID: "Step-3",
			Action: func(ctx context.Context, data interface{}) (error, interface{}) {
				fmt.Println("3->")
				return nil, nil
			},
			CompensateAction: func(ctx context.Context, data interface{}) (error, interface{}) {
				fmt.Println("3<-")
				return nil, nil
			},
		},
	)
	mem := saga.NewInMemoryCarrier()
	trans := saga.NewTransaction(orderCreatedSaga, mem)

	output := trans.Start()
	fmt.Println(output)
}
