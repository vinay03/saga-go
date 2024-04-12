package saga

type MessageCarrier interface {
}

type MessageCarrierOptions struct {
}

func NewInMemoryCarrier() *InMemoryMessageCarrier {
	return &InMemoryMessageCarrier{}
}

type InMemoryMessageCarrier struct {
	Data map[string][]string
}
