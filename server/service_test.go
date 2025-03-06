package server

import (
	"context"
	"testing"

	"anujkumar.com/elevator/lift"
	"anujkumar.com/elevator/lift/mocks"
	"anujkumar.com/elevator/models"
	"github.com/stretchr/testify/assert"
)

func TestAssign(t *testing.T) {
	ctx := context.Background()

	req := models.Request{Source: 3, Dest: 7}
	lift1 := mocks.NewILift(t)
	lift1.On("GetCurrentFloor", ctx).Return(models.Floor(2)).Once()
	lift2 := mocks.NewILift(t)
	lift2.On("GetCurrentFloor", ctx).Return(models.Floor(3)).Once()
	lifts := []lift.ILift{lift1, lift2}
	lift2.On("Assign", ctx, req).Return(nil).Once()
	// lift2.AssertCalled(t, "Assign", ctx, req)

	service := NewService(lifts)

	assignedLift, err := service.Assign(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, assignedLift)
	assert.Equal(t, lift2, assignedLift)
}
