package tests

import (
	"context"
	"errors"
	"fmt"
	"log"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/vinay03/saga-go"
)

func TestSagaGo(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "SagaGo Suite")
}

func InitInMemSaga() {
	sampleSaga, err := saga.NewSaga("SampleSaga").Transactions(
		"Step1",
		"Step2",
		"Step3",
	)
	if err != nil {
		log.Fatal(err)
	}

	sampleSaga.DefineSubTransactions(
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
	sampleSaga.DefineSubTransactions(
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
	sampleSaga.DefineSubTransactions(
		"Step3",
		func(ctx context.Context, data interface{}) (interface{}, error) {
			fmt.Println("3->")
			// log.WithFields(log.Fields{
			// 	"service": "BasicService",
			// }).Error("Error in Step3")
			return nil, nil
			// return nil, errors.New("Aborted")
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
}
