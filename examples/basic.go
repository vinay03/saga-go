package main

import (
	"context"
	"errors"
	"fmt"
	"log"

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

	coord := saga.GetCoordinatorInstance()
	coord.SetupCarriers(
		&saga.InMemoryCarrierConfig{},
		// saga.RedisCarrierOption{},
	)
	coord.RegisterSaga(OrderCreatedSaga, coord.Carrier.InMem)
	coord.RegisterSaga(TransferSaga, coord.Carrier.InMem)

	// data := struct {
	// 	Testdata string `json:"testdata"`
	// }{
	// 	"Test Sample string",
	// }

	// coord.Start("OrderCreated", data)

	transferPayload := TransferRequest{
		SenderUserID:        1,
		SenderUserBalance:   1000,
		ReceiverUserId:      2,
		ReceiverUserBalance: 1000,
		TransferAmount:      50,
	}
	coord.Start("Transfer", transferPayload)

	// coord.SetupCarrierOptions

	// mem := saga.NewInMemoryCarrier()
	// trans := saga.NewCoordinator(mem)

	// output := trans.Start()
	// fmt.Println(output)

	/* orderCreatedSaga := saga.NewSaga("Order_Created")
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
	trans := saga.NewOperator(orderCreatedSaga, mem)

	output := trans.Start()
	fmt.Println(output) */
}
