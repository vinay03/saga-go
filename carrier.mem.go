package saga

type InMemoryCarrier struct {
	Active        bool
	Options       CarrierConfig
	EventListener func(Message string, Data interface{})
}

func getInMemoryCarrierInstance() *InMemoryCarrier {
	return &InMemoryCarrier{
		Active: false,
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
