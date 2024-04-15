package saga

import "log"

type InMemoryCarrier struct {
	Active        bool
	Options       CarrierConfig
	Data          map[string][]string
	EventListener func(Message string, Data interface{})
}

func getInMemoryCarrierInstance() *InMemoryCarrier {
	return &InMemoryCarrier{
		Active: false,
		Data:   make(map[string][]string),
	}
}

func (mem *InMemoryCarrier) IsActive() bool {
	return mem.Active
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

func (mem *InMemoryCarrier) Push(Message string, Data interface{}) error {
	log.Println("Pushing Message: ", Message, " Data: ", Data)
	mem.EventListener(Message, Data)
	return nil
}
func (mem *InMemoryCarrier) AddListener(handlerFunc func(Message string, Data interface{})) error {
	mem.EventListener = handlerFunc
	return nil
}

/* In-Memory Carrier Configuration */
type InMemoryCarrierConfig struct {
}

func (memOpt *InMemoryCarrierConfig) Verify() error {
	return nil
}
