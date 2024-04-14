package saga

import (
	"errors"
	"fmt"
)

type Carrier interface {
	IsActive() bool
	SetOptions(CarrierConfig) error
}

type CarrierConfig interface {
	Verify() error
}

type CarrierLineup struct {
	InMem *InMemoryCarrier
	Redis *RedisCarrier
}

func (cl *CarrierLineup) SetupCarriers(options ...CarrierConfig) error {
	for _, opts := range options {
		switch v := opts.(type) {
		case *InMemoryCarrierConfig:
			cl.InMem.SetOptions(v)
			fmt.Println(v)
		case *RedisCarrierOption:
			cl.Redis.SetOptions(v)
			fmt.Println(v)
		default:
			return errors.New("invalid carrier option")
		}
	}
	return nil
}
