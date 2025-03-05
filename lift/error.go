package lift

import "fmt"

type IdleError struct {
	LiftId uint
}

func NewIdleError(liftId uint) *IdleError {
	return &IdleError{LiftId: liftId}
}

func (e *IdleError) Error() string {
	return fmt.Sprintf("Lift %d is idle", e.LiftId)
}
