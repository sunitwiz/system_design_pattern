package scheduler

import (
	"elevator_system/elevator"
	"elevator_system/request"
	"fmt"
)

type SchedulerType int

const (
	Nearest SchedulerType = iota
	RoundRobin
)

func (s SchedulerType) String() string {
	switch s {
	case Nearest:
		return "Nearest"
	case RoundRobin:
		return "RoundRobin"
	default:
		return "Unknown"
	}
}

type ElevatorScheduler interface {
	AssignElevator(elevators []*elevator.Elevator, req request.Request) *elevator.Elevator
	String() string
}

func NewScheduler(schedulerType SchedulerType) (ElevatorScheduler, error) {
	switch schedulerType {
	case Nearest:
		return &NearestElevatorScheduler{}, nil
	case RoundRobin:
		return &RoundRobinScheduler{}, nil
	default:
		return nil, fmt.Errorf("unknown scheduler type: %d", schedulerType)
	}
}
