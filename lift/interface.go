package lift

import (
	"context"

	"anujkumar.com/elevator/models"
)

type ILift interface {
	MoveToNextStop(context.Context) error
	Assign(context.Context, models.Request) error
	GetNextStop(context.Context) (models.Floor, LiftState)
	Operate(context.Context) error
	GetCurrentFloor(context.Context) models.Floor
	GetState(context.Context) LiftState
}
