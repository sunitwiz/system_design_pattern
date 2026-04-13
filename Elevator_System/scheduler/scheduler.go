package scheduler

import "fmt"

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
