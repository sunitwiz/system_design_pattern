package service

import (
	"fmt"
	"splitwise/expense"
	"splitwise/group"
	"splitwise/split"
	"splitwise/user"
	"sync"
)

type SplitwiseService struct {
	mu             sync.Mutex
	users          map[string]*user.User
	groups         map[string]*group.Group
	expenses       []*expense.Expense
	expenseCounter int
}

var (
	instance *SplitwiseService
	once     sync.Once
)

func GetInstance() *SplitwiseService {
	once.Do(func() {
		instance = &SplitwiseService{
			users:  make(map[string]*user.User),
			groups: make(map[string]*group.Group),
		}
	})
	return instance
}

func ResetInstance() {
	once = sync.Once{}
	instance = nil
}

func (s *SplitwiseService) AddUser(name, email string) string {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := fmt.Sprintf("U%d", len(s.users)+1)
	s.users[id] = user.NewUser(id, name, email)
	return id
}

func (s *SplitwiseService) CreateGroup(name string, memberIDs []string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, mid := range memberIDs {
		if _, exists := s.users[mid]; !exists {
			return "", fmt.Errorf("user %s not found", mid)
		}
	}

	id := fmt.Sprintf("G%d", len(s.groups)+1)
	g := group.NewGroup(id, name)
	for _, mid := range memberIDs {
		g.AddMember(mid)
	}
	s.groups[id] = g
	return id, nil
}

func (s *SplitwiseService) AddExpenseToGroup(groupID, payerID string, amount float64, splitType split.SplitType, participantIDs []string, details map[string]float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	g, exists := s.groups[groupID]
	if !exists {
		return fmt.Errorf("group %s not found", groupID)
	}

	if _, exists := s.users[payerID]; !exists {
		return fmt.Errorf("payer %s not found", payerID)
	}

	splitter, err := split.NewSplit(splitType)
	if err != nil {
		return err
	}

	shares, err := splitter.Calculate(amount, participantIDs, details)
	if err != nil {
		return err
	}

	s.expenseCounter++
	expID := fmt.Sprintf("E%d", s.expenseCounter)
	exp := expense.NewExpense(expID, payerID, amount, splitType, participantIDs, shares)

	g.AddExpense(exp)
	s.expenses = append(s.expenses, exp)

	for uid, share := range shares {
		if uid == payerID {
			continue
		}
		s.users[uid].UpdateBalance(payerID, share)
		s.users[payerID].UpdateBalance(uid, -share)
	}

	return nil
}

func (s *SplitwiseService) GetBalances(userID string) map[string]float64 {
	s.mu.Lock()
	defer s.mu.Unlock()

	u, exists := s.users[userID]
	if !exists {
		return nil
	}
	return u.GetAllBalances()
}

func (s *SplitwiseService) GetGroupExpenses(groupID string) []*expense.Expense {
	s.mu.Lock()
	defer s.mu.Unlock()

	g, exists := s.groups[groupID]
	if !exists {
		return nil
	}
	return g.Expenses
}

func (s *SplitwiseService) SettleUp(fromUserID, toUserID string, amount float64) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	from, exists := s.users[fromUserID]
	if !exists {
		return fmt.Errorf("user %s not found", fromUserID)
	}

	to, exists := s.users[toUserID]
	if !exists {
		return fmt.Errorf("user %s not found", toUserID)
	}

	from.UpdateBalance(toUserID, -amount)
	to.UpdateBalance(fromUserID, amount)
	return nil
}

func (s *SplitwiseService) ViewStatus() {
	s.mu.Lock()
	defer s.mu.Unlock()

	fmt.Println("========================================")
	fmt.Println("       SPLITWISE STATUS")
	fmt.Println("========================================")

	fmt.Printf("\n  Users: %d\n", len(s.users))
	fmt.Printf("  %-6s %-12s %s\n", "ID", "Name", "Email")
	fmt.Printf("  %-6s %-12s %s\n", "------", "------------", "--------------------")
	for _, u := range s.users {
		fmt.Printf("  %-6s %-12s %s\n", u.ID, u.Name, u.Email)
	}

	fmt.Printf("\n  Groups: %d\n", len(s.groups))
	for _, g := range s.groups {
		memberNames := s.resolveNames(g.Members)
		fmt.Printf("  [%s] %s → Members: %v | Expenses: %d\n",
			g.ID, g.Name, memberNames, len(g.Expenses))
	}

	fmt.Println("\n  Balances:")
	fmt.Printf("  %-12s %-12s %10s\n", "User", "With", "Amount")
	fmt.Printf("  %-12s %-12s %10s\n", "------------", "------------", "----------")
	for _, u := range s.users {
		balances := u.GetAllBalances()
		for otherID, amt := range balances {
			otherName := s.userName(otherID)
			if amt > 0 {
				fmt.Printf("  %-12s %-12s %10s\n", u.Name, otherName, fmt.Sprintf("owes ₹%.2f", amt))
			} else if amt < 0 {
				fmt.Printf("  %-12s %-12s %10s\n", u.Name, otherName, fmt.Sprintf("gets ₹%.2f", -amt))
			}
		}
	}

	fmt.Println("========================================")
}

func (s *SplitwiseService) resolveNames(ids []string) []string {
	names := make([]string, 0, len(ids))
	for _, id := range ids {
		if u, ok := s.users[id]; ok {
			names = append(names, u.Name)
		}
	}
	return names
}

func (s *SplitwiseService) userName(id string) string {
	if u, ok := s.users[id]; ok {
		return u.Name
	}
	return id
}
