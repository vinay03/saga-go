package saga

import (
	"context"
	"errors"
	"fmt"
	"reflect"
)

type ActionFunction func(context.Context, interface{}) (error, interface{})

type Stage struct {
	ID               string
	Action           ActionFunction
	CompensateAction ActionFunction
}

var (
	EmptyStageIDError               = errors.New("Blank Stage ID")
	ActionFuncInvalidParameterError = errors.New("Action function has invalid parameters")
)

func (st *Stage) Verify() error {

	actionValueType := reflect.TypeOf(st.Action)
	if actionValueKind := actionValueType.Kind(); actionValueKind != reflect.Func {
		return fmt.Errorf("Action field should be a function. Provided %s instead", actionValueKind)
	}
	if actionValueType.NumOut() != 2 {
		return errors.New("Action function must return 2 values in format (error, interface{}).")
	} else if actionValueType.Out(0) != reflect.TypeOf((*error)(nil)).Elem() {
		return errors.New("Action function must return error as its first return value.")
	}

	CompensateActionValueType := reflect.TypeOf(st.CompensateAction)
	if CompensateActionValueKind := CompensateActionValueType.Kind(); CompensateActionValueKind != reflect.Func {
		return fmt.Errorf("CompensateAction field should be a function. Provided %s instead", CompensateActionValueKind)
	}

	return nil
}
