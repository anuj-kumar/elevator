package lift

import (
	"context"
	"fmt"
	"time"

	"github.com/oleiade/lane/v2"

	"anujkumar.com/elevator/models"
)

type LiftState int

const (
	Idle      LiftState = 0
	GoingUp   LiftState = 1
	GoingDown LiftState = 2
)

type LiftImpl struct {
	Id            uint
	Capacity      uint
	Floors        uint
	stopsUpward   *lane.PriorityQueue[models.Floor, int64]
	stopsDownward *lane.PriorityQueue[models.Floor, int64]
	State         LiftState
	CurrentFloor  models.Floor
}

func NewLiftImpl(id uint, floors uint) ILift {
	return &LiftImpl{
		Id:            id,
		Floors:        floors,
		stopsUpward:   lane.NewMinPriorityQueue[models.Floor, int64](),
		stopsDownward: lane.NewMaxPriorityQueue[models.Floor, int64](),
		State:         Idle,
		CurrentFloor:  0,
	}
}

func (m *LiftImpl) MoveToNextStop(ctx context.Context) error {
	stop, state := m.GetNextStop(ctx)
	m.State = state
	return m.moveToFloor(ctx, stop)
}

func (m *LiftImpl) GetCurrentFloor(ctx context.Context) models.Floor {
	return m.CurrentFloor
}

func (m *LiftImpl) GetState(ctx context.Context) LiftState {
	return m.State
}

func (m *LiftImpl) moveToFloor(_ context.Context, floor models.Floor) error {
	var incr models.Floor = 1
	if m.State == Idle {
		return NewIdleError(m.Id)
	}
	if m.State == GoingDown {
		incr = -1
	}
	fmt.Printf("Lift %d Dir %d Destination %d\n", m.Id, m.State, floor)
	for ; m.CurrentFloor != floor; m.CurrentFloor += incr {
		fmt.Printf("Lift %d moving to floor %d\n", m.Id, m.CurrentFloor)
		time.Sleep(500 * time.Millisecond)
	}
	fmt.Printf("Lift %d reached floor %d\n", m.Id, floor)
	return nil
}

func (m *LiftImpl) GetNextStop(ctx context.Context) (models.Floor, LiftState) {
	return m.NextStop(ctx, true)
}

func (m *LiftImpl) NextStop(ctx context.Context, pop bool) (models.Floor, LiftState) {
	var stop models.Floor
	var state LiftState
	if m.State != GoingDown && !m.stopsUpward.Empty() {
		state = GoingUp
		stop, _, _ = m.stopsUpward.Pop()
	} else if m.State != GoingUp && !m.stopsDownward.Empty() {
		state = GoingDown
		stop, _, _ = m.stopsDownward.Pop()
	} else {
		return 0, Idle
	}
	return stop, state
}

func (m *LiftImpl) Assign(ctx context.Context, req models.Request) error {
	if req.Dest > req.Source {
		m.stopsUpward.Push(req.Dest, int64(req.Dest))
	} else {
		m.stopsDownward.Push(req.Dest, int64(req.Dest))
	}
	fmt.Printf("Lift %d assigned to request: %v, upward len: %d, downward len: %d", m.Id, req, m.stopsUpward.Size(), m.stopsDownward.Size())
	return nil
}

func (m *LiftImpl) Operate(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		default:
			err := m.MoveToNextStop(ctx)
			if err != nil {
				if _, ok := err.(*IdleError); ok {
					time.Sleep(1 * time.Second)
				} else {
					fmt.Printf("Lift: %d, Error: %v\n", m.Id, err)
				}
			}
		}
	}
}
