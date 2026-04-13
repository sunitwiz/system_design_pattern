package controller

import (
	"elevator_system/elevator"
	"elevator_system/scheduler"
)

type ElevatorOperations interface {
	RequestElevator(sourceFloor, destFloor int) (*elevator.Elevator, error)
	StepAll()
	ViewStatus()
}

type AdminOperations interface {
	AddElevator(id int)
	RemoveElevator(id int) error
	SetScheduler(s scheduler.ElevatorScheduler)
}
