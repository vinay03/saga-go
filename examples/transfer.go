package main

import (
	"context"
	"log"

	"github.com/vinay03/saga-go"
)

var TransferSaga *saga.Saga

type TransferRequest struct {
	SenderUserID        int     `json:"sender_user_id"`
	SenderUserBalance   float32 `json:"sender_user_balance"`
	ReceiverUserId      int     `json:"receiver_user_id"`
	ReceiverUserBalance float32 `json:"receiver_user_balance"`
	TransferAmount      float32 `json:"transfer_amount"`
	SenderNotified      bool
	ReceiverNotified    bool
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
	return payload, nil
}
func (s *NotificationService) NotifyReceiverWithUpdatedBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	transferPayload := payload.(TransferRequest)
	transferPayload.ReceiverNotified = true
	return transferPayload, nil
}
func (s *NotificationService) CompensateNotifyReceiverWithUpdatedBalance(ctx context.Context, payload interface{}) (interface{}, error) {
	return payload, nil
}

func init() {
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

}
