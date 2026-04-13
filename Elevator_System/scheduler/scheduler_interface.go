package scheduler

import (
	"elevator_system/elevator"
	"elevator_system/request"
)

type ElevatorScheduler interface {
	AssignElevator(elevators []*elevator.Elevator, req request.Request) *elevator.Elevator
	String() string
}
