package group

import "splitwise/expense"

type Group struct {
	ID       string
	Name     string
	Members  []string
	Expenses []*expense.Expense
}

func NewGroup(id, name string) *Group {
	return &Group{
		ID:   id,
		Name: name,
	}
}

func (g *Group) AddMember(userID string) {
	g.Members = append(g.Members, userID)
}

func (g *Group) AddExpense(exp *expense.Expense) {
	g.Expenses = append(g.Expenses, exp)
}
