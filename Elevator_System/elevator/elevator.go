package elevator

import (
	"elevator_system/request"
	"fmt"
	"sort"
)

type Status int

const (
	StatusIdle Status = iota
	StatusMoving
	StatusMaintenance
)

func (s Status) String() string {
	switch s {
	case StatusIdle:
		return "Idle"
	case StatusMoving:
		return "Moving"
	case StatusMaintenance:
		return "Maintenance"
	default:
		return "Unknown"
	}
}

type Elevator struct {
	ID           int
	CurrentFloor int
	Direction    request.Direction
	Status       Status
	Requests     []int
}

func NewElevator(id int) *Elevator {
	return &Elevator{
		ID:           id,
		CurrentFloor: 1,
		Direction:    request.Idle,
		Status:       StatusIdle,
	}
}

func (e *Elevator) AddRequest(floor int) {
	for _, f := range e.Requests {
		if f == floor {
			return
		}
	}

	e.Requests = append(e.Requests, floor)

	if e.Status == StatusIdle {
		e.Status = StatusMoving
	}

	e.sortRequests()
	e.updateDirection()
}

func (e *Elevator) MoveOneStep() {
	if len(e.Requests) == 0 || e.Status != StatusMoving {
		return
	}

	if e.Direction == request.Up {
		e.CurrentFloor++
	} else if e.Direction == request.Down {
		e.CurrentFloor--
	}

	if e.hasArrived() {
		e.removeCurrentFloor()
	}

	if len(e.Requests) == 0 {
		e.Status = StatusIdle
		e.Direction = request.Idle
	} else {
		e.updateDirection()
	}
}

func (e *Elevator) GetDirection() request.Direction {
	return e.Direction
}

func (e *Elevator) IsIdle() bool {
	return e.Status == StatusIdle && len(e.Requests) == 0
}

func (e *Elevator) String() string {
	return fmt.Sprintf("Elevator %d [Floor: %d, Dir: %s, Status: %s, Pending: %v]",
		e.ID, e.CurrentFloor, e.Direction, e.Status, e.Requests)
}

func (e *Elevator) sortRequests() {
	if e.Direction == request.Down {
		sort.Sort(sort.Reverse(sort.IntSlice(e.Requests)))
	} else {
		sort.Ints(e.Requests)
	}
}

func (e *Elevator) updateDirection() {
	if len(e.Requests) == 0 {
		e.Direction = request.Idle
		return
	}

	next := e.Requests[0]
	if next > e.CurrentFloor {
		e.Direction = request.Up
	} else if next < e.CurrentFloor {
		e.Direction = request.Down
	} else if len(e.Requests) > 1 {
		if e.Requests[1] > e.CurrentFloor {
			e.Direction = request.Up
		} else {
			e.Direction = request.Down
		}
	}
}

func (e *Elevator) hasArrived() bool {
	for _, f := range e.Requests {
		if f == e.CurrentFloor {
			return true
		}
	}
	return false
}

func (e *Elevator) removeCurrentFloor() {
	filtered := e.Requests[:0]
	for _, f := range e.Requests {
		if f != e.CurrentFloor {
			filtered = append(filtered, f)
		}
	}
	e.Requests = filtered
}
