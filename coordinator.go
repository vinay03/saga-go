package saga

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
)

var Coord *Coordinator
var CoordLock sync.Mutex

const (
	EventKeyFormat = "%s|%s|%s"

	ErrSagaNotFound  = "[%s] Saga Id not found"
	ErrStageNotFound = "[%s] Stage Id '%s' not found"
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

// SetupCarriers sets up the carriers for the SAGA Execution Coordinator (SEC)
//
// It takes in a variadic parameter of CarrierConfig types, which can be any of the carrier configuration classes the implements the CarrierConfig interface.
// CarrierOptions are set to the corresponding carrier and the EventHandler is added as a listener to the carrier events.
//
// Returns:
// - error: an error if the carrier option is invalid, otherwise nil.
func (coord *Coordinator) SetupCarriers(options ...CarrierConfig) error {
	for _, opts := range options {
		switch v := opts.(type) {
		case *InMemoryCarrierConfig:
			err := coord.Carrier.InMem.SetOptions(v)
			if err != nil {
				return err
			}
			coord.Carrier.InMem.AddListener(coord.EventHandler)
		case *RedisCarrierOption:
			err := coord.Carrier.Redis.SetOptions(v)
			if err != nil {
				return err
			}
			coord.Carrier.InMem.AddListener(coord.EventHandler)
		default:
			return errors.New("invalid carrier option")
		}
	}
	return nil
}

// EventHandler handles incoming event messages and performs actions based on the event key.
//
// Parameters:
// - eventkey: the event message received (string)
// - data: the data associated with the event (interface{})
//
// Returns: None
func (coord *Coordinator) EventHandler(eventkey string, data interface{}) {
	// log.Println("Received eventkey: ", eventkey, " Data: ", data)

	sagaId, stageId, eventAction := decodeEventKey(eventkey)
	value, ok := coord.Sagas[sagaId]
	if !ok {
		log.Fatal(fmt.Errorf(ErrSagaNotFound, sagaId))
	}

	// blank stage Id means it is the first event of the SAGA
	if stageId == "" {
		stageId = value.Saga.GetFirstStage().ID
	}

	stage, found := value.Saga.StagesNameRef[stageId]
	if !found {
		log.Fatalf(ErrStageNotFound, sagaId, stageId)
	}

	switch eventAction {
	case "start":
		data, err := stage.Action(context.Background(), data)
		if err != nil {
			eventKey := generateEventKey(sagaId, stageId, "abort")
			value.Carrier.Push(eventKey, data)
		}
		// Call start action of the next stage or complete the SAGA.
		nextStage := value.Saga.GetNextStage(stage)
		if nextStage != nil {
			eventKey := generateEventKey(sagaId, nextStage.ID, "start")
			value.Carrier.Push(eventKey, data)
		} else {
			// End of SAGA
		}
	case "abort":
		data, err := stage.CompensateAction(context.Background(), data)
		if err != nil {
			eventKey := generateEventKey(sagaId, stageId, "abort")
			value.Carrier.Push(eventKey, data)
		}
		// call abort action of previous stage or abort the SAGA completely
		prevStage := value.Saga.GetPrevStage(stage)
		if prevStage != nil {
			eventKey := generateEventKey(sagaId, prevStage.ID, "abort")
			value.Carrier.Push(eventKey, data)
		} else {
			// End of SAGA Abortion sequence
		}
	default:
		log.Fatalf("[%s] Invalid event action: %s", eventkey, eventAction)
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

func decodeEventKey(eventKey string) (string, string, string) {
	parts := strings.Split(eventKey, "|")
	return parts[0], parts[1], parts[2]
}

func (coord *Coordinator) Start(sagaId string, data interface{}) (interface{}, error) {

	val, ok := coord.Sagas[sagaId]
	if !ok {
		return data, fmt.Errorf(ErrSagaNotFound, sagaId)
	}

	eventKey := generateEventKey(sagaId, "", "start")
	val.Carrier.Push(eventKey, data)

	return data, nil
}
func (tr *Coordinator) Abort(data interface{}) interface{} {
	tr.IsAborted = true
	return data
}
