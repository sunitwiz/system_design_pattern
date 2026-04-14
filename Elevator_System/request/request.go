package request

import "fmt"

type Direction int

const (
	Up Direction = iota
	Down
	Idle
)



type RequestType int

const (
	External RequestType = iota
	Internal
)



type Request struct {
	SourceFloor      int
	DestinationFloor int
	Direction        Direction
	Type             RequestType
}

func NewRequest(sourceFloor, destinationFloor int, reqType RequestType) Request {
	dir := Idle
	if destinationFloor > sourceFloor {
		dir = Up
	} else if destinationFloor < sourceFloor {
		dir = Down
	}

	return Request{
		SourceFloor:      sourceFloor,
		DestinationFloor: destinationFloor,
		Direction:        dir,
		Type:             reqType,
	}
}

func (r Request) String() string {
	return fmt.Sprintf("Request[Floor %d → %d, %s, %s]",
		r.SourceFloor, r.DestinationFloor, r.Direction, r.Type)
}
