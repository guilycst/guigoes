package domain

import "fmt"

type FSResourceNotFoundError struct {
	Msg string
	Err error
}

func (e *FSResourceNotFoundError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Resource not found: %s\nCause: %s", e.Msg, e.Err.Error())
	}

	return fmt.Sprintln("Resource not found:", e.Msg)
}
