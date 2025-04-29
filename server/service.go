package server

import (
	"context"
	"math"

	"anujkumar.com/elevator/lift"
	"anujkumar.com/elevator/models"
)

type simpleAssignmentStrategy struct {
	Lifts []lift.ILift
}

func NewService(lifts []lift.ILift) *simpleAssignmentStrategy {
	return &simpleAssignmentStrategy{
		Lifts: lifts,
	}
}

func (s *simpleAssignmentStrategy) Assign(ctx context.Context, req models.Request) (lift.ILift, error) {
	// Assign the nearest lift
	var target lift.ILift
	for _, l := range s.Lifts {
		if target == nil || math.Abs(float64(l.GetCurrentFloor(ctx)-req.Source)) < math.Abs(float64(target.GetCurrentFloor(ctx)-req.Source)) {
			target = l
		}
	}
	target.Assign(ctx, req)
	return target, nil
}
