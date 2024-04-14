package main

import (
	"context"
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
		"CheckDiscount",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered `CheckDiscount` phase")
			return data, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered compensating `CheckDiscount` phase")
			return data, nil
		},
	)
	OrderCreatedSaga.DefineActions(
		"NotifySeller",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered `NotifySeller` phase")
			return data, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered compensating `NotifySeller` phase")
			return data, nil
		},
	)
	OrderCreatedSaga.DefineActions(
		"NotifySeller",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered `NotifySeller` phase")
			return data, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("Entered compensating `NotifySeller` phase")
			return data, nil
		},
	)

	coord := saga.GetCoordinatorInstance()
	coord.Carrier.SetupCarriers(
		&saga.InMemoryCarrierConfig{},
		// saga.RedisCarrierOption{},
	)

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
