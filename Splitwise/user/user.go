package user

type User struct {
	ID       string
	Name     string
	Email    string
	balances map[string]float64
}

func NewUser(id, name, email string) *User {
	return &User{
		ID:       id,
		Name:     name,
		Email:    email,
		balances: make(map[string]float64),
	}
}

func (u *User) GetBalance(userID string) float64 {
	return u.balances[userID]
}

func (u *User) UpdateBalance(userID string, amount float64) {
	u.balances[userID] += amount
}

func (u *User) GetAllBalances() map[string]float64 {
	result := make(map[string]float64)
	for k, v := range u.balances {
		if v != 0 {
			result[k] = v
		}
	}
	return result
}
