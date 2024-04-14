package saga

import (
	"context"
	"sync"
)

var Coord *Coordinator
var CoordLock sync.Mutex

func GetCoordinatorInstance() *Coordinator {
	if Coord != nil {
		return Coord
	}
	CoordLock.Lock()
	defer CoordLock.Unlock()
	// Just in case if instance is created within the Mutex locking operation.
	if Coord != nil {
		return Coord
	}
	Coord = &Coordinator{
		Carrier: &CarrierLineup{
			InMem: getInMemoryCarrierInstance(),
			Redis: getRedisCarrierInstance(),
		},
	}
	return Coord
}

type Coordinator struct {
	Template  *Saga
	IsAborted bool
	Carrier   *CarrierLineup

	In  interface{}
	Out interface{}
}

func (tr *Coordinator) Start() interface{} {
	input := tr.In
	var err error
	for _, st := range tr.Template.Stages {
		input, err = st.Action(context.Background(), input)
		if err != nil {
			return tr.Abort()
		}
	}
	tr.Out = input
	return tr.Out
}
func (tr *Coordinator) Abort() interface{} {
	tr.IsAborted = true

	return nil
}
