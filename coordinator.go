package saga

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
)

var Coord *Coordinator
var CoordLock sync.Mutex

const (
	EventKeyFormat = "%s|%s|%s"

	ErrSagaNotFound = "[%s] Saga Id not found"
)

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
		Sagas: make(map[string]CoordinatorSaga),
	}
	Coord.Carrier = &CarrierLineup{
		InMem: getInMemoryCarrierInstance(),
		Redis: getRedisCarrierInstance(),
	}
	return Coord
}

type CoordinatorSaga struct {
	Saga    *Saga
	Carrier Carrier
}

type Coordinator struct {
	Template  *Saga
	IsAborted bool
	Carrier   *CarrierLineup

	Sagas map[string]CoordinatorSaga
}

func (coord *Coordinator) SetupCarriers(options ...CarrierConfig) error {
	for _, opts := range options {
		switch v := opts.(type) {
		case *InMemoryCarrierConfig:
			coord.Carrier.InMem.SetOptions(v)
			coord.Carrier.InMem.AddListener(coord.MessageHandler)
		case *RedisCarrierOption:
			coord.Carrier.Redis.SetOptions(v)
			coord.Carrier.InMem.AddListener(coord.MessageHandler)
		default:
			return errors.New("invalid carrier option")
		}
	}
	return nil
}

func (coord *Coordinator) MessageHandler(message string, data interface{}) {
	log.Println("Received Message: ", message, " Data: ", data)
	var sagaId, stageId, eventAction string
	fmt.Sscanf(message, EventKeyFormat, sagaId, stageId, eventAction)

	value, ok := coord.Sagas[sagaId]
	if !ok {
		log.Fatal(fmt.Errorf(ErrSagaNotFound, sagaId))
	}

	// blank stage Id means it is the first event of the SAGA
	if stageId == "" {
		stageId = value.Saga.Stages[0].ID
	}

	stage := value.Saga.StagesNameRef[stageId]

	switch eventAction {
	case "start":
		data, err := stage.Action(context.Background(), data)
		if err != nil {
			eventKey := generateEventKey(sagaId, stageId, "abort")
			value.Carrier.Push(eventKey, data)
		}
		// Call start action of the next stage or complete the SAGA.
	case "abort":
		data, err := stage.CompensateAction(context.Background(), data)
		if err != nil {
			eventKey := generateEventKey(sagaId, stageId, "abort")
			value.Carrier.Push(eventKey, data)
		}
		// call abort action of previous stage or abort the SAGA completely
	default:
		log.Fatalf("[%s] Invalid event action: %s", message, eventAction)
	}

}

func (coord *Coordinator) RegisterSaga(saga *Saga, carr Carrier) {
	coord.Sagas[saga.ID] = CoordinatorSaga{
		Saga:    saga,
		Carrier: carr,
	}
}

func generateEventKey(sagaId, stageId, eventAction string) string {
	return fmt.Sprintf(EventKeyFormat, sagaId, stageId, eventAction)
}

func (coord *Coordinator) Start(sagaId string, data interface{}) (interface{}, error) {

	val, ok := coord.Sagas[sagaId]
	if !ok {
		return data, fmt.Errorf(ErrSagaNotFound, sagaId)
	}

	eventKey := generateEventKey(sagaId, "", "start")
	val.Carrier.Push(eventKey, data)

	// var err error
	// for _, st := range val.Saga.Stages {
	// 	data, err = st.Action(context.Background(), data)
	// 	if err != nil {
	// 		return coord.Abort(data), fmt.Errorf("[%s] Saga aborted", sagaId)
	// 	}
	// }
	return data, nil
}
func (tr *Coordinator) Abort(data interface{}) interface{} {
	tr.IsAborted = true
	return data
}
