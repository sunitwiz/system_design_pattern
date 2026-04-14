package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class Direction {
        <<enumeration>>
        Up
        Down
        Idle
        func (d Direction) String() string
    }

    class RequestType {
        <<enumeration>>
        External
        Internal
        func (r RequestType) String() string
    }

    class Request {
        SourceFloor      int
        DestinationFloor int
        Direction        Direction
        Type             RequestType
        func NewRequest(sourceFloor, destinationFloor int, reqType RequestType) Request
        func (r Request) String() string
    }

    class Status {
        <<enumeration>>
        StatusIdle
        StatusMoving
        StatusMaintenance
        func (s Status) String() string
    }

    class Elevator {
        ID           int
        CurrentFloor int
        Direction    request.Direction
        Status       Status
        Requests     []int
        func NewElevator(id int) *Elevator
        func (e *Elevator) AddRequest(floor int)
        func (e *Elevator) MoveOneStep()
        func (e *Elevator) GetDirection() request.Direction
        func (e *Elevator) IsIdle() bool
        func (e *Elevator) String() string
        func (e *Elevator) sortRequests()
        func (e *Elevator) updateDirection()
        func (e *Elevator) hasArrived() bool
        func (e *Elevator) removeCurrentFloor()
    }

    class SchedulerType {
        <<enumeration>>
        Nearest
        RoundRobin
        func (s SchedulerType) String() string
    }

    class ElevatorScheduler {
        <<interface>>
        AssignElevator(elevators []*elevator.Elevator, req request.Request) *elevator.Elevator
        String() string
    }

    class NearestElevatorScheduler {
        func (n *NearestElevatorScheduler) AssignElevator(elevators []*elevator.Elevator, req request.Request) *elevator.Elevator
        func (n *NearestElevatorScheduler) String() string
    }

    class RoundRobinScheduler {
        lastIndex int
        func (rr *RoundRobinScheduler) AssignElevator(elevators []*elevator.Elevator, req request.Request) *elevator.Elevator
        func (rr *RoundRobinScheduler) String() string
    }

    class ElevatorOperations {
        <<interface>>
        RequestElevator(sourceFloor, destFloor int) (*elevator.Elevator, error)
        StepAll()
        ViewStatus()
    }

    class AdminOperations {
        <<interface>>
        AddElevator(id int)
        RemoveElevator(id int) error
        SetScheduler(s scheduler.ElevatorScheduler)
    }

    class ElevatorController {
        mu        sync.Mutex
        Elevators []*elevator.Elevator
        Scheduler scheduler.ElevatorScheduler
        func GetInstance() *ElevatorController
        func (ec *ElevatorController) RequestElevator(sourceFloor, destFloor int) (*elevator.Elevator, error)
        func (ec *ElevatorController) StepAll()
        func (ec *ElevatorController) ViewStatus()
        func (ec *ElevatorController) AddElevator(id int)
        func (ec *ElevatorController) RemoveElevator(id int) error
        func (ec *ElevatorController) SetScheduler(s scheduler.ElevatorScheduler)
    }

    Request --> Direction
    Request --> RequestType
    Elevator --> Direction
    Elevator --> Status
    NearestElevatorScheduler ..|> ElevatorScheduler : implements
    RoundRobinScheduler ..|> ElevatorScheduler : implements
    ElevatorController ..|> ElevatorOperations : implements
    ElevatorController ..|> AdminOperations : implements
    ElevatorController o-- Elevator : manages
    ElevatorController --> ElevatorScheduler : uses`)
}
