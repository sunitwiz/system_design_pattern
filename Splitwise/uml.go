package main

import "fmt"

func printClassDiagram() {
	fmt.Println(`classDiagram

    class SplitType {
        <<enumeration>>
        EqualSplit
        ExactSplit
        PercentSplit
    }

    class Split {
        <<interface>>
        Calculate(totalAmount float64, participants []string, details map[string]float64) (map[string]float64, error)
        GetType() SplitType
    }

    class equalSplit {
        func (e *equalSplit) Calculate(totalAmount float64, participants []string, _ map[string]float64) (map[string]float64, error)
        func (e *equalSplit) GetType() SplitType
    }

    class exactSplit {
        func (ex *exactSplit) Calculate(totalAmount float64, _ []string, details map[string]float64) (map[string]float64, error)
        func (ex *exactSplit) GetType() SplitType
    }

    class percentSplit {
        func (ps *percentSplit) Calculate(totalAmount float64, _ []string, details map[string]float64) (map[string]float64, error)
        func (ps *percentSplit) GetType() SplitType
    }

    class User {
        ID       string
        Name     string
        Email    string
        balances map[string]float64
        func NewUser(id, name, email string) *User
        func (u *User) GetBalance(userID string) float64
        func (u *User) UpdateBalance(userID string, amount float64)
        func (u *User) GetAllBalances() map[string]float64
    }

    class Group {
        ID       string
        Name     string
        Members  []string
        Expenses []*expense.Expense
        func NewGroup(id, name string) *Group
        func (g *Group) AddMember(userID string)
        func (g *Group) AddExpense(exp *expense.Expense)
    }

    class Expense {
        ID           string
        PayerID      string
        Amount       float64
        SplitType    split.SplitType
        Participants []string
        SplitDetails map[string]float64
        func NewExpense(id, payerID string, amount float64, splitType split.SplitType, participants []string, splitDetails map[string]float64) *Expense
    }

    class SplitwiseOperations {
        <<interface>>
        AddUser(name, email string) string
        CreateGroup(name string, memberIDs []string) (string, error)
        AddExpenseToGroup(groupID, payerID string, amount float64, splitType split.SplitType, participantIDs []string, details map[string]float64) error
        GetBalances(userID string) map[string]float64
        GetGroupExpenses(groupID string) []*expense.Expense
        SettleUp(fromUserID, toUserID string, amount float64) error
        ViewStatus()
    }

    class SplitwiseService {
        mu             sync.Mutex
        users          map[string]*user.User
        groups         map[string]*group.Group
        expenses       []*expense.Expense
        expenseCounter int
        func GetInstance() *SplitwiseService
        func (s *SplitwiseService) AddUser(name, email string) string
        func (s *SplitwiseService) CreateGroup(name string, memberIDs []string) (string, error)
        func (s *SplitwiseService) AddExpenseToGroup(groupID, payerID string, amount float64, splitType split.SplitType, participantIDs []string, details map[string]float64) error
        func (s *SplitwiseService) GetBalances(userID string) map[string]float64
        func (s *SplitwiseService) GetGroupExpenses(groupID string) []*expense.Expense
        func (s *SplitwiseService) SettleUp(fromUserID, toUserID string, amount float64) error
        func (s *SplitwiseService) ViewStatus()
    }

    equalSplit ..|> Split : implements
    exactSplit ..|> Split : implements
    percentSplit ..|> Split : implements
    Expense --> SplitType
    Group *-- Expense : contains
    Group --> User : references members by ID
    SplitwiseService ..|> SplitwiseOperations : implements
    SplitwiseService o-- User : manages
    SplitwiseService o-- Group : manages
    SplitwiseService o-- Expense : tracks
    SplitwiseService ..> Split : uses for calculation`)
}
