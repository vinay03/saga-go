package saga

import (
	"errors"
	"log"
)

var (
	ErrEmptySagaID = errors.New("empty stage identifier")
)

// Initializes a new SAGA with the given identifier
func NewSaga(id string) *Saga {
	saga := _NewSaga()
	err := VerifySagaId(id)
	if err != nil {
		log.Fatal("Invalid Saga ID: ", id)
	}
	saga.ID = id
	return saga
}

// Initializes the default state of any SAGA
func _NewSaga() *Saga {
	return &Saga{
		StagesNameRef: make(map[string]*Stage),
	}
}

type StageList []*Stage

type Saga struct {
	ID            string
	Stages        StageList
	StagesNameRef map[string]*Stage
}

func (sg *Saga) Transactions(stageIds ...string) (*Saga, error) {
	for _, stageId := range stageIds {
		err := VerifyStageID(stageId)
		if err != nil {
			return sg, err
		}
		stage := &Stage{
			ID: stageId,
		}
		sg.AddStage(stage)
	}
	return sg, nil
}

func VerifySagaId(id string) error {
	// TODO: validate id string for allowed format.
	// 	1. Should not be empty string
	// 	2. Only Alphabets and Numbers are allowed
	// 	3. Should not already exist under same SAGA.
	if id == "" {
		return ErrEmptySagaID
	}
	return nil
}

func (sg *Saga) AddStages(sl ...*Stage) error {
	for _, st := range sl {
		err := sg.AddStage(st)
		if err != nil {
			return err
		}
	}
	return nil
}

func (sg *Saga) AddStage(st *Stage) error {
	err := st.Verify()
	if err != nil {
		return err
	}
	sg.Stages = append(sg.Stages, st)
	sg.StagesNameRef[st.ID] = st
	return nil
}

func (sg *Saga) DefineActions(stageId string, action ActionFunction, compensateAction ActionFunction) error {
	for _, st := range sg.Stages {
		if st.ID == stageId {
			st.Action = action
			st.CompensateAction = compensateAction
			return nil
		}
	}
	return errors.New("transaction not found")
}
