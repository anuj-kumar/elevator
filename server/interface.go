package server

import (
	"anujkumar.com/elevator/lift"
)

type Server interface {
	Assign() (lift.ILift, error)
}
