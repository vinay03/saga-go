package saga

type Carrier interface {
	IsActive() bool
	Activate()
	SetOptions(CarrierConfig) error
	Push(Message string, Data interface{}) error
	AddEventsListener(func(Message string, Data interface{})) error
}

type CarrierConfig interface {
	Verify() error
}

type CarrierLineup struct {
	InMem Carrier
	Redis Carrier
}
