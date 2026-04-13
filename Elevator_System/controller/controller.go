package controller

import (
	"elevator_system/elevator"
	"elevator_system/request"
	"elevator_system/scheduler"
	"fmt"
	"sync"
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

type ElevatorController struct {
	mu        sync.Mutex
	Elevators []*elevator.Elevator
	Scheduler scheduler.ElevatorScheduler
}

var _ ElevatorOperations = (*ElevatorController)(nil)
var _ AdminOperations = (*ElevatorController)(nil)

var (
	instance *ElevatorController
	once     sync.Once
)

func GetInstance() *ElevatorController {
	once.Do(func() {
		defaultScheduler, _ := scheduler.NewScheduler(scheduler.Nearest)
		instance = &ElevatorController{
			Scheduler: defaultScheduler,
		}
	})
	return instance
}

func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

func (ec *ElevatorController) RequestElevator(sourceFloor, destFloor int) (*elevator.Elevator, error) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	if len(ec.Elevators) == 0 {
		return nil, fmt.Errorf("no elevators available")
	}

	req := request.NewRequest(sourceFloor, destFloor, request.External)
	assigned := ec.Scheduler.AssignElevator(ec.Elevators, req)
	if assigned == nil {
		return nil, fmt.Errorf("no suitable elevator for %s", req)
	}

	assigned.AddRequest(sourceFloor)
	assigned.AddRequest(destFloor)

	return assigned, nil
}

func (ec *ElevatorController) StepAll() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	for _, e := range ec.Elevators {
		e.MoveOneStep()
	}
}

func (ec *ElevatorController) ViewStatus() {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	fmt.Println("========================================")
	fmt.Println("       ELEVATOR SYSTEM STATUS")
	fmt.Println("========================================")

	if len(ec.Elevators) == 0 {
		fmt.Println("  No elevators configured.")
		fmt.Println("========================================")
		return
	}

	fmt.Printf("  Scheduler: %s\n\n", ec.Scheduler)
	fmt.Printf("  %-4s %-8s %-10s %-12s %s\n", "ID", "Floor", "Direction", "Status", "Pending")
	fmt.Printf("  %-4s %-8s %-10s %-12s %s\n", "----", "--------", "----------", "------------", "-------")

	for _, e := range ec.Elevators {
		fmt.Printf("  %-4d %-8d %-10s %-12s %v\n",
			e.ID, e.CurrentFloor, e.Direction, e.Status, e.Requests)
	}

	fmt.Println("========================================")
}

func (ec *ElevatorController) AddElevator(id int) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	for _, e := range ec.Elevators {
		if e.ID == id {
			fmt.Printf("  Elevator %d already exists\n", id)
			return
		}
	}

	e := elevator.NewElevator(id)
	ec.Elevators = append(ec.Elevators, e)
	fmt.Printf("  Added Elevator %d\n", id)
}

func (ec *ElevatorController) RemoveElevator(id int) error {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	for i, e := range ec.Elevators {
		if e.ID == id {
			if !e.IsIdle() {
				return fmt.Errorf("cannot remove elevator %d: has pending requests", id)
			}
			ec.Elevators = append(ec.Elevators[:i], ec.Elevators[i+1:]...)
			fmt.Printf("  Removed Elevator %d\n", id)
			return nil
		}
	}
	return fmt.Errorf("elevator %d not found", id)
}

func (ec *ElevatorController) SetScheduler(s scheduler.ElevatorScheduler) {
	ec.mu.Lock()
	defer ec.mu.Unlock()

	ec.Scheduler = s
	fmt.Printf("  Scheduler set to: %s\n", s)
}
