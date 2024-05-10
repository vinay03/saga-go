package saga

import "errors"

type InMemoryCarrier struct {
	Active        bool
	Options       CarrierConfig
	EventListener func(Message string, Data interface{})
	EventReceiver chan EventData
}

// A factory method to return pointer to instance of `InMemoryCarrier`
func getInMemoryCarrierInstance() *InMemoryCarrier {
	return &InMemoryCarrier{
		Active: false,
	}
}

func (mem *InMemoryCarrier) IsActive() bool {
	return mem.Active
}

func (mem *InMemoryCarrier) Activate() error {
	if mem.EventListener == nil {
		return errors.New("Event Listener not set")
	}
	mem.Active = true
	return nil
}

func (mem *InMemoryCarrier) SetOptions(opts CarrierConfig) error {
	val, _ := opts.(*InMemoryCarrierConfig)
	err := val.Verify()
	if err != nil {
		return err
	}
	mem.Options = opts
	return nil
}

// func (mem *InMemoryCarrier) RegisterEvent(key string, func)

func (mem *InMemoryCarrier) Push(Message string, Data interface{}) error {
	go mem.EventListener(Message, Data)
	return nil
}

func (mem *InMemoryCarrier) AddEventsListener(handlerFunc func(Message string, Data interface{})) error {
	mem.EventListener = handlerFunc
	return nil
}

// ************************************************************
/* In-Memory Carrier Configuration */
type InMemoryCarrierConfig struct {
}

func (memCfg *InMemoryCarrierConfig) Verify() error {
	return nil
}
