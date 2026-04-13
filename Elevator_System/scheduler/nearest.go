package scheduler

import (
	"elevator_system/elevator"
	"elevator_system/request"
	"math"
)

type NearestElevatorScheduler struct{}

func (n *NearestElevatorScheduler) AssignElevator(elevators []*elevator.Elevator, req request.Request) *elevator.Elevator {
	var best *elevator.Elevator
	bestDistance := math.MaxInt

	for _, e := range elevators {
		if e.Status == elevator.StatusMaintenance {
			continue
		}

		distance := abs(e.CurrentFloor - req.SourceFloor)

		if e.IsIdle() {
			if distance < bestDistance {
				bestDistance = distance
				best = e
			}
			continue
		}

		sameDirection := e.GetDirection() == req.Direction
		movingToward := (e.GetDirection() == request.Up && e.CurrentFloor <= req.SourceFloor) ||
			(e.GetDirection() == request.Down && e.CurrentFloor >= req.SourceFloor)

		if sameDirection && movingToward && distance < bestDistance {
			bestDistance = distance
			best = e
		}
	}

	if best != nil {
		return best
	}

	for _, e := range elevators {
		if e.Status == elevator.StatusMaintenance {
			continue
		}
		distance := abs(e.CurrentFloor - req.SourceFloor)
		if distance < bestDistance {
			bestDistance = distance
			best = e
		}
	}

	return best
}

func (n *NearestElevatorScheduler) String() string {
	return "NearestElevatorScheduler"
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Compile-time interface check.
var _ ElevatorScheduler = (*NearestElevatorScheduler)(nil)
