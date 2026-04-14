package service

import (
	"splitwise/expense"
	"splitwise/split"
)

type SplitwiseOperations interface {
	AddUser(name, email string) string
	CreateGroup(name string, memberIDs []string) (string, error)
	AddExpenseToGroup(groupID, payerID string, amount float64, splitType split.SplitType, participantIDs []string, details map[string]float64) error
	GetBalances(userID string) map[string]float64
	GetGroupExpenses(groupID string) []*expense.Expense
	SettleUp(fromUserID, toUserID string, amount float64) error
	ViewStatus()
}

var _ SplitwiseOperations = (*SplitwiseService)(nil)
