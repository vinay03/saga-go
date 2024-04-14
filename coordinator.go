package saga

import (
	"context"
)

func NewTransaction(sg *Saga, cr MessageCarrier) *Transaction {
	return &Transaction{
		Template: sg,
		Carrier:  cr,
	}
}

type Transaction struct {
	Template  *Saga
	IsAborted bool
	Carrier   MessageCarrier

	In  interface{}
	Out interface{}
}

func (tr *Transaction) Start() interface{} {
	input := tr.In
	var err error
	for _, st := range tr.Template.Stages {
		err, input = st.Action(context.Background(), input)
		if err != nil {
			return tr.Abort()
		}
	}
	tr.Out = input
	return tr.Out
}
func (tr *Transaction) Abort() interface{} {
	tr.IsAborted = true

	return nil
}
