package saga

import (
	"errors"
	"fmt"
	"log"
)

var (
	ErrEmptySagaID = errors.New("empty stage identifier")
)

// NewSaga initializes a new Saga with the given identifier.
//
// @Param id the identifier of the Saga (string). Same will be put in event key to identify the Saga. Alphabets and numbers are allowed only.
//
// Returns:
// - a pointer to the newly created Saga.
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
	StagesCount   int
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
		err = sg.AddStage(stage)
		if err != nil {
			log.Fatal(err)
		}
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
	st.SequenceNumber = len(sg.Stages) + 1
	sg.Stages = append(sg.Stages, st)
	sg.StagesNameRef[st.ID] = st
	sg.StagesCount = len(sg.Stages)
	return nil
}

func (sg *Saga) DefineActions(stageId string, action ActionFunction, compensateAction ActionFunction) error {
	stage, found := sg.StagesNameRef[stageId]
	if !found {
		return fmt.Errorf(ErrStageNotFound, sg.ID, stageId)
	}
	stage.Action = action
	stage.CompensateAction = compensateAction
	return nil
}

func (sg *Saga) GetFirstStage() *Stage {
	return sg.Stages[0]
}
func (sg *Saga) GetNextStage(stage *Stage) *Stage {
	if sg.StagesCount > stage.SequenceNumber {
		return sg.Stages[stage.SequenceNumber]
	}
	return nil
}
func (sg *Saga) GetPrevStage(stage *Stage) *Stage {
	if stage.SequenceNumber > 1 {
		return sg.Stages[stage.SequenceNumber-2]
	}
	return nil
}
