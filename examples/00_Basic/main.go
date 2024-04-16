package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/vinay03/saga-go"
)

func main() {
	sampleSaga, err := saga.NewSaga("SampleSaga").Transactions(
		"Step1",
		"Step2",
		"Step3",
	)
	if err != nil {
		log.Fatal(err)
	}

	sampleSaga.DefineActions(
		"Step1",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("1->")
			return nil, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("1<-")
			return nil, nil
		},
	)
	sampleSaga.DefineActions(
		"Step2",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("2->")
			return errors.New("some"), nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("2<-")
			return nil, nil
		},
	)
	sampleSaga.DefineActions(
		"Step3",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("3->")
			return nil, nil
		},
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("3<-")
			return nil, nil
		},
	)

	coord := saga.GetCoordinatorInstance()
	coord.SetupCarriers(
		&saga.InMemoryCarrierConfig{},
	)
	coord.RegisterSaga(sampleSaga, coord.Carrier.InMem)

	data := struct{}{}
	start := time.Now()
	for i := 0; i < 100; i++ {
		coord.Start("SampleSaga", data)
	}

	fmt.Println("Time Elapsed: ", time.Since(start))
	time.Sleep(2 * time.Second)
}
