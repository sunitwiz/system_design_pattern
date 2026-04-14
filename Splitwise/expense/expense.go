package expense

import "splitwise/split"

type Expense struct {
	ID           string
	PayerID      string
	Amount       float64
	SplitType    split.SplitType
	Participants []string
	SplitDetails map[string]float64
}

func NewExpense(id, payerID string, amount float64, splitType split.SplitType, participants []string, splitDetails map[string]float64) *Expense {
	return &Expense{
		ID:           id,
		PayerID:      payerID,
		Amount:       amount,
		SplitType:    splitType,
		Participants: participants,
		SplitDetails: splitDetails,
	}
}
