package saga

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
)

var Coord *Coordinator
var CoordLock sync.Mutex

const (
	EventKeyFormat = "%s|%s|%s"

	ErrSagaNotFound     = "[%s] Saga Id not found"
	ErrStagesNotDefined = "[%s] Stages not defined"
	ErrStageNotFound    = "[%s] Stage Id '%s' not found"
)

func GetCoordinatorInstance(logger Logger) *Coordinator {
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
		Sagas:  make(map[string]CoordinatorSaga),
		Logger: logger,
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
	Logger    Logger

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
			coord.Carrier.InMem.AddEventsListener(coord.EventHandler)
			coord.Logger.Info("InMemory Carrier is active")
		case *RedisCarrierOption:
			err := coord.Carrier.Redis.SetOptions(v)
			if err != nil {
				return err
			}
			coord.Carrier.InMem.AddEventsListener(coord.EventHandler)
			coord.Logger.Info("Redis Carrier is active")
		default:
			return errors.New("invalid carrier option")
		}
	}
	return nil
}

func (coord *Coordinator) DecodeEventKey(eventkey string) (sagaRecord *CoordinatorSaga, stage *Stage, eventAction string, err error) {
	var sagaId, stageId string
	sagaId, stageId, eventAction = splitEventKey(eventkey)

	value, ok := coord.Sagas[sagaId]
	if !ok {
		err = fmt.Errorf(ErrSagaNotFound, sagaId)
		return
	}
	sagaRecord = &value

	// blank stage Id means it is the first event of the SAGA
	if stageId == "" {
		stageId = sagaRecord.Saga.GetFirstStage().ID
	}

	stage, found := sagaRecord.Saga.StagesNameRef[stageId]
	if !found {
		err = fmt.Errorf(ErrStageNotFound, sagaId, stageId)
		return
	}
	return sagaRecord, stage, eventAction, err
}

// EventHandler handles incoming event messages and performs actions based on the event key.
//
// Parameters:
// - eventkey: the event message received (string)
// - data: the data associated with the event (interface{})
//
// Returns: None
func (coord *Coordinator) EventHandler(eventkey string, data interface{}) {
	coord.Logger.Info(fmt.Sprintf("[%s] Event received", eventkey))
	sagaRecord, stage, eventAction, err := coord.DecodeEventKey(eventkey)
	if err != nil {
		coord.Logger.Error(err)
	}

	sagaId := sagaRecord.Saga.ID
	stageId := stage.ID

	switch eventAction {
	case "start":
		data, err := stage.Action(context.Background(), data)
		if err != nil {
			eventKey := generateEventKey(sagaId, stageId, "abort")
			coord.PushEvent(sagaRecord.Carrier, eventKey, data)
		}
		// Call start action of the next stage or complete the SAGA.
		nextStage := sagaRecord.Saga.GetNextStage(stage)
		if nextStage != nil {
			eventKey := generateEventKey(sagaId, nextStage.ID, "start")
			coord.PushEvent(sagaRecord.Carrier, eventKey, data)
		} else {
			// End of SAGA
		}
	case "abort":
		data, err := stage.CompensateAction(context.Background(), data)
		if err != nil {
			eventKey := generateEventKey(sagaId, stageId, "abort")
			coord.PushEvent(sagaRecord.Carrier, eventKey, data)
		}
		// call abort action of previous stage or abort the SAGA completely
		prevStage := sagaRecord.Saga.GetPrevStage(stage)
		if prevStage != nil {
			eventKey := generateEventKey(sagaId, prevStage.ID, "abort")
			coord.PushEvent(sagaRecord.Carrier, eventKey, data)
		} else {
			// End of SAGA Abortion sequence
		}
	default:
		coord.Logger.Error(fmt.Sprintf("[%s] Invalid event action: %s", eventkey, eventAction))
	}
}

func (coord *Coordinator) PushEvent(carrier Carrier, eventKey string, data interface{}) {
	coord.Logger.Info(fmt.Sprintf("[%s] Pushing event", eventKey))
	carrier.Push(eventKey, data)
}

// RegisterSaga registers a new saga with the coordinator.
// It associates the given saga and carrier with the saga ID in the coordinator's Sagas map.
func (coord *Coordinator) RegisterSaga(saga *Saga, carr Carrier) {
	coord.Sagas[saga.ID] = CoordinatorSaga{
		Saga:    saga,
		Carrier: carr,
	}
}

// generateEventKey generates a unique key for an event based on the given saga ID, stage ID, and event action.
func generateEventKey(sagaId, stageId, eventAction string) string {
	return fmt.Sprintf(EventKeyFormat, sagaId, stageId, eventAction)
}

// splitEventKey splits the given eventKey into three parts.
// It expects the eventKey to be in the format "part1|part2|part3".
// It returns three strings representing the three parts of the eventKey.
//
// Parameters:
// - eventKey: The eventKey to be split into three parts.
//
// Returns:
// - string: Saga Id.
// - string: Stage Id.
// - string: Event Action.
func splitEventKey(eventKey string) (string, string, string) {
	parts := strings.Split(eventKey, "|")
	return parts[0], parts[1], parts[2]
}

// Start initiates the execution of a saga with the given sagaId and data.
// The sagaId is a unique identifier for the saga, and data is the input data
// required for the saga execution.
//
// Parameters:
// - sagaId: The unique identifier for the saga.
// - data: The input data required for the saga execution.
//
// Returns:
// - interface{}: The result of the saga execution.
// - error: An error, if any, occurred during the saga execution.
func (coord *Coordinator) Start(sagaId string, data interface{}) error {

	val, ok := coord.Sagas[sagaId]
	if !ok {
		return fmt.Errorf(ErrSagaNotFound, sagaId)
	}

	if val.Saga.StagesCount == 0 {
		return fmt.Errorf(ErrStagesNotDefined, sagaId)
	}

	stage := val.Saga.GetFirstStage()

	eventKey := generateEventKey(sagaId, stage.ID, "start")
	coord.PushEvent(val.Carrier, eventKey, data)

	return nil
}

// func (tr *Coordinator) Abort(data interface{}) interface{} {
// 	tr.IsAborted = true
// 	return data
// }
