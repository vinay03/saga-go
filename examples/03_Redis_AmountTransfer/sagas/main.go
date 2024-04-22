package sagas

import (
	log "github.com/sirupsen/logrus"
	saga "github.com/vinay03/saga-go"
)

var TransferAmountSaga *saga.Saga

func SetupSagas() {
	var err error
	TransferAmountSaga, err = saga.NewSaga("TransferAmountSaga").Transactions(
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

/*
const redisAddr = "redis:6379"
const redisPassword = "admin"

var Client *asynq.Client

func InitBroker() {
	Client = asynq.NewClient(asynq.RedisClientOpt{
		Addr:     redisAddr,
		Password: redisPassword,
	})
}

type TransferAmountPayload struct {
	SenderID   int     `json:"sender_id"`
	ReceiverID int     `json:"receiver_id"`
	Amount     float32 `json:"amount"`
}

func NewTransferAmountTask(senderID int, receiverID int, amount float32) (*asynq.Task, error) {
	payload, err := json.Marshal(TransferAmountPayload{
		SenderID:   senderID,
		ReceiverID: receiverID,
		Amount:     amount,
	})
	if err != nil {
		return nil, err
	}
	fmt.Println("Creating Transfer Amount Task")
	return asynq.NewTask(
		TypeEmailDelivery,
		payload,
	), nil
}

func CloseBroker() {
	Client.Close()
}
/*/
