package split

import (
	"fmt"
	"math"
)

type SplitType int

const (
	EqualSplit SplitType = iota
	ExactSplit
	PercentSplit
)

func (s SplitType) String() string {
	switch s {
	case EqualSplit:
		return "Equal"
	case ExactSplit:
		return "Exact"
	case PercentSplit:
		return "Percent"
	default:
		return "Unknown"
	}
}

type Split interface {
	Calculate(totalAmount float64, participants []string, details map[string]float64) (map[string]float64, error)
	GetType() SplitType
}

func NewSplit(splitType SplitType) (Split, error) {
	switch splitType {
	case EqualSplit:
		return &equalSplit{}, nil
	case ExactSplit:
		return &exactSplit{}, nil
	case PercentSplit:
		return &percentSplit{}, nil
	default:
		return nil, fmt.Errorf("unknown split type: %d", splitType)
	}
}

type equalSplit struct{}

func (e *equalSplit) Calculate(totalAmount float64, participants []string, _ map[string]float64) (map[string]float64, error) {
	if len(participants) == 0 {
		return nil, fmt.Errorf("no participants provided")
	}
	share := math.Round(totalAmount/float64(len(participants))*100) / 100
	result := make(map[string]float64)
	for _, p := range participants {
		result[p] = share
	}
	return result, nil
}

func (e *equalSplit) GetType() SplitType { return EqualSplit }

type exactSplit struct{}

func (ex *exactSplit) Calculate(totalAmount float64, _ []string, details map[string]float64) (map[string]float64, error) {
	if len(details) == 0 {
		return nil, fmt.Errorf("no split details provided")
	}
	var sum float64
	for _, amt := range details {
		sum += amt
	}
	if math.Abs(sum-totalAmount) > 0.01 {
		return nil, fmt.Errorf("exact split amounts (%.2f) do not equal total (%.2f)", sum, totalAmount)
	}
	result := make(map[string]float64)
	for k, v := range details {
		result[k] = v
	}
	return result, nil
}

func (ex *exactSplit) GetType() SplitType { return ExactSplit }

type percentSplit struct{}

func (ps *percentSplit) Calculate(totalAmount float64, _ []string, details map[string]float64) (map[string]float64, error) {
	if len(details) == 0 {
		return nil, fmt.Errorf("no split details provided")
	}
	var totalPercent float64
	for _, pct := range details {
		totalPercent += pct
	}
	if math.Abs(totalPercent-100.0) > 0.01 {
		return nil, fmt.Errorf("percentages sum to %.2f, expected 100", totalPercent)
	}
	result := make(map[string]float64)
	for k, pct := range details {
		result[k] = math.Round(totalAmount*pct) / 100
	}
	return result, nil
}

func (ps *percentSplit) GetType() SplitType { return PercentSplit }
