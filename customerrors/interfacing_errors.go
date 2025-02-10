package customerrors

import "fmt"

type InterfacingError struct {
	InterfaceName string
	Err           error
}

func (e *InterfacingError) Error() string {
	return fmt.Sprintf("interface setup error function %s, error : %s", e.InterfaceName, e.Err.Error())
}

func WrapinterfacingError(interfaceName string, err error) *InterfacingError {
	if err == nil {
		return nil
	}

	return &InterfacingError{
		InterfaceName: interfaceName,
		Err:           err,
	}
}

func (e *InterfacingError) UnWrap() error { return e.Err }
