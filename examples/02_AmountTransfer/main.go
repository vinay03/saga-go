package main

import (
	"context"
	"fmt"
	"log"
	"time"

	saga "github.com/vinay03/saga-go"
)

var TransferSaga *saga.Saga

type TransferRequest struct {
	SenderUserID        int     `json:"sender_user_id"`
	SenderUserBalance   float32 `json:"sender_user_balance"`
	ReceiverUserId      int     `json:"receiver_user_id"`
	ReceiverUserBalance float32 `json:"receiver_user_balance"`
	TransferAmount      float32 `json:"transfer_amount"`
	SenderNotified      bool    `json:"sender_notified"`
	ReceiverNotified    bool    `json:"receiver_notified"`
}

var Users = make(map[int]User)

type User struct {
	ID      int
	Balance float32
}

type AccountService struct {
}

func (s *AccountService) DeduceBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.SenderUserBalance -= transferPayload.TransferAmount
	return transferPayload, nil
}
func (s *AccountService) CompensateDeduceBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.SenderUserBalance += transferPayload.TransferAmount
	fmt.Printf("%+v", transferPayload)
	return transferPayload, nil
}

func (s *AccountService) IncreaseBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.ReceiverUserBalance += transferPayload.TransferAmount
	return transferPayload, nil
}
func (s *AccountService) CompensateIncreaseBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.ReceiverUserBalance -= transferPayload.TransferAmount
	return transferPayload, nil
}

type NotificationService struct {
}

func (s *NotificationService) NotifySenderWithUpdatedBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.SenderNotified = true
	return transferPayload, nil
}
func (s *NotificationService) CompensateNotifySenderWithUpdatedBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.SenderNotified = false
	return payload, nil
}
func (s *NotificationService) NotifyReceiverWithUpdatedBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.ReceiverNotified = true
	fmt.Printf("%+v  \n", transferPayload)
	return transferPayload, nil
}
func (s *NotificationService) CompensateNotifyReceiverWithUpdatedBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.ReceiverNotified = false
	return payload, nil
}

func main() {
	accountServ := AccountService{}
	notificationServ := NotificationService{}

	var err error
	TransferSaga, err = saga.NewSaga("Transfer").Transactions(
		"DeduceSenderBalancer",
		"IncreaseReceiverBalance",
		"NotifySenderWithUpdatedBalance",
		"NotifyReceiverWithUpdatedBalance",
	)
	if err != nil {
		log.Fatal(err)
	}

	TransferSaga.DefineActions(
		"DeduceSenderBalancer",
		accountServ.DeduceBalance,
		accountServ.CompensateDeduceBalance,
	)

	TransferSaga.DefineActions(
		"IncreaseReceiverBalance",
		accountServ.IncreaseBalance,
		accountServ.CompensateIncreaseBalance,
	)

	TransferSaga.DefineActions(
		"NotifySenderWithUpdatedBalance",
		notificationServ.NotifySenderWithUpdatedBalance,
		notificationServ.CompensateNotifySenderWithUpdatedBalance,
	)

	TransferSaga.DefineActions(
		"NotifyReceiverWithUpdatedBalance",
		notificationServ.NotifyReceiverWithUpdatedBalance,
		notificationServ.CompensateNotifyReceiverWithUpdatedBalance,
	)

	coord := saga.GetCoordinatorInstance()
	coord.SetupCarriers(
		&saga.InMemoryCarrierConfig{},
	)
	coord.RegisterSaga(TransferSaga, coord.Carrier.InMem)

	transferPayload := TransferRequest{
		SenderUserID:        1,
		SenderUserBalance:   1000,
		ReceiverUserId:      2,
		ReceiverUserBalance: 1000,
		TransferAmount:      50,
	}
	coord.Start("Transfer", transferPayload)

	time.Sleep(1 * time.Second)
}
