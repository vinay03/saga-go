package saga

type InMemoryCarrier struct {
	IsInit  bool
	Active  bool
	Options CarrierConfig
	Data    map[string][]string
}
type InMemoryCarrierConfig struct {
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

func (memOpt *InMemoryCarrierConfig) Verify() error {
	return nil
}
