package main

import (
	saga "github.com/vinay03/saga-go"
)

var SagaCoordinator *saga.Coordinator

func main() {
	SagaCoordinator = saga.GetCoordinatorInstance()

	// sagas.TransferAmountSaga.DefineSubTransactions("DeduceSenderBalancer")
	// sagas.TransferAmountSaga.DefineSubTransactions("IncreaseReceiverBalance")

	// SagaCoordinator.RegisterSaga("TransferAmountSaga")

	// // Connect to Redis
	// rdb := redis.NewClient(&redis.Options{
	// 	Addr:     "localhost:6379",
	// 	Password: "", // no password set
	// 	DB:       0,  // use default DB
	// })

	// // Subscribe to the channel
	// pubsub := rdb.Subscribe(context.Background(), "events")
	// defer pubsub.Close()

	// // Wait for confirmation that subscription is created before publishing anything
	// _, err := pubsub.Receive(context.Background())
	// if err != nil {
	// 	panic(err)
	// }

	// // Wait for messages
	// channel := pubsub.Channel()
	// for msg := range channel {
	// 	fmt.Println("Received message:", msg.Payload)
	// }
}
