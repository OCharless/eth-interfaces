package base

import (
	"fmt"
)

type CallError struct {
	Field string
	Err   error
}

func (e *CallError) Error() string {
	return fmt.Sprintf("call error on %s: %s", e.Field, e.Err.Error())
}

func (d *BaseInteractions) WrapCallError(abiString, field string, err error) *CallError {
	if err == nil {
		return nil
	}
	if err != nil {
		if err := d.ManageCustomContractError(abiString, err); err != nil {
			return &CallError{
				Field: field,
				Err:   err,
			}
		}
	}
	return &CallError{
		Field: field,
		Err:   err,
	}
}

func (e *CallError) UnWrap() error { return e.Err }
