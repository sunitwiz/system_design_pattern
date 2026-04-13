package scheduler

import (
	"elevator_system/elevator"
	"elevator_system/request"
)

type RoundRobinScheduler struct {
	lastIndex int
}

func (rr *RoundRobinScheduler) AssignElevator(elevators []*elevator.Elevator, req request.Request) *elevator.Elevator {
	if len(elevators) == 0 {
		return nil
	}

	n := len(elevators)
	for i := 0; i < n; i++ {
		idx := (rr.lastIndex + 1 + i) % n
		if elevators[idx].Status != elevator.StatusMaintenance {
			rr.lastIndex = idx
			return elevators[idx]
		}
	}

	return nil
}

func (rr *RoundRobinScheduler) String() string {
	return "RoundRobinScheduler"
}

// Compile-time interface check.
var _ ElevatorScheduler = (*RoundRobinScheduler)(nil)
