package saga

import "errors"

type Stage struct {
	ID               string
	Action           interface{}
	CompensateAction interface{}
}

var (
	EmptyStageIDError               = errors.New("Blank Stage ID")
	ActionFuncInvalidParameterError = errors.New("Action function has invalid parameters")
)

func (st *Stage) Verify() error {
	return nil
}
