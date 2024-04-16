package saga

import (
	"crypto/md5"
	"fmt"
	"math/rand"
	"time"
)

const (
	RANDOM_NUMBER_SEED_LIMIT = 1048576
)

func NewOperation() *Operation {
	return &Operation{
		ID: generateOperationId(),
	}
}

func generateOperationId() string {
	randNum := rand.Intn(RANDOM_NUMBER_SEED_LIMIT)
	now := time.Now()
	hash := md5.Sum([]byte(fmt.Sprintf("%d-%d-%d", randNum, now.UnixNano(), now.Nanosecond())))
	return string(hash[:])
}

type Operation struct {
	ID                 string
	Saga               *Saga
	Stage              *Stage
	CurrentEventAction string
	currentEventKey    string
}
