package lift

import (
	"context"
	"fmt"
	"time"

	"anujkumar.com/elevator/models"
	pq "gopkg.in/dnaeon/go-priorityqueue.v1"
)

type LiftState int

const (
	Idle      LiftState = 0
	GoingUp             = 1
	GoingDown           = 2
)

type Model struct {
	Id            uint
	Capacity      uint
	Floors        uint
	stopsUpward   *pq.PriorityQueue[models.Floor, int64]
	stopsDownward *pq.PriorityQueue[models.Floor, int64]
	State         LiftState
	CurrentFloor  models.Floor
}

func NewModel(id uint, floors uint) *Model {
	return &Model{
		Id:            id,
		Floors:        floors,
		stopsUpward:   pq.New[models.Floor, int64](pq.MinHeap),
		stopsDownward: pq.New[models.Floor, int64](pq.MaxHeap),
		State:         Idle,
		CurrentFloor:  0,
	}
}

func (m *Model) MoveToNextStop(ctx context.Context) error {
	stop, state := m.GetNextStop(ctx)
	return m.moveToFloor(ctx, stop, state)
}

func (m *Model) GetCurrentFloor(ctx context.Context) models.Floor {
	return m.CurrentFloor
}

func (m *Model) GetState(ctx context.Context) LiftState {
	return m.State
}

func (m *Model) moveToFloor(_ context.Context, floor models.Floor, direction LiftState) error {
	var incr models.Floor = 1
	if direction == Idle {
		return NewIdleError(m.Id)
	}
	if direction == GoingDown {
		incr = -1
	}
	for ; m.CurrentFloor < floor; m.CurrentFloor += incr {
		fmt.Printf("Lift %d moving to floor %d\n", m.Id, m.CurrentFloor)
		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

func (m *Model) GetNextStop(ctx context.Context) (models.Floor, LiftState) {
	return m.NextStop(ctx, true)
}

func (m *Model) NextStop(ctx context.Context, pop bool) (models.Floor, LiftState) {
	var stop *pq.Item[models.Floor, int64]
	var state LiftState
	if m.State != GoingDown && m.stopsUpward.Len() > 0 {
		state = GoingUp
		stop = m.stopsUpward.Pop().(*pq.Item[models.Floor, int64])
	} else if m.State != GoingUp && m.stopsDownward.Len() > 0 {
		state = GoingDown
		stop = m.stopsDownward.Pop().(*pq.Item[models.Floor, int64])
	} else {
		return 0, m.State
	}
	return stop.Value, state
}

func (m *Model) Assign(ctx context.Context, req models.Request) error {
	if req.Source > m.CurrentFloor {
		m.stopsUpward.Put(req.Source, int64(req.Source))
	} else {
		m.stopsDownward.Put(req.Source, int64(req.Source))
	}
	if req.Dest > req.Source {
		m.stopsUpward.Put(req.Dest, int64(req.Dest))
	} else {
		m.stopsDownward.Put(req.Dest, int64(req.Dest))
	}
	return nil
}

func (m *Model) Operate(ctx context.Context) error {
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
