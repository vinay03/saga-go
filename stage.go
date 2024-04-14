package saga

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type ActionFunction func(context.Context, interface{}) (interface{}, error)

type Stage struct {
	ID               string
	Action           ActionFunction
	CompensateAction ActionFunction
}

var (
	ErrEmptyStageID               = errors.New("blank stage id")
	ErrActionFuncInvalidParameter = errors.New("action function has invalid parameters")
)

func VerifyStageID(stageID string) error {
	if stageID == "" {
		return ErrEmptyStageID
	}
	return nil
}

func (st *Stage) Verify() error {
	// TODO: validate all the configurations of a stage in SAGA.

	actionValueType := reflect.TypeOf(st.Action)
	if actionValueKind := actionValueType.Kind(); actionValueKind != reflect.Func {
		return fmt.Errorf("action field should be a function. provided %s instead", actionValueKind)
	}
	if actionValueType.NumOut() != 2 {
		return errors.New("action function must return 2 values in format (error, interface{})")
	} else if actionValueType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.New("action function must return error as its first return value")
	}

	CompensateActionValueType := reflect.TypeOf(st.CompensateAction)
	if CompensateActionValueKind := CompensateActionValueType.Kind(); CompensateActionValueKind != reflect.Func {
		return fmt.Errorf("CompensateAction field should be a function. Provided %s instead", CompensateActionValueKind)
	}

	return nil
}
