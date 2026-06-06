package storage

import (
	"sync"
	"time"
)

type UserRepository struct {
	mu sync.RWMutex

	users map[string]User
}

type Role string

const (
	UserRole Role = "USER"

	AdminRole Role = "ADMIN"
)

type Balance struct {
	Asset string

	Available string

	Locked string
}

type User struct {
	ID    string
	Name  string
	Login string
	//Password   string
	Balances   map[string]*Balance
	Role       Role
	CreateTime time.Time
	UpdateTime time.Time
}

func NewUserStorage() *UserRepository {
	return &UserRepository{
		users: make(map[string]User),
	}
}

func (r *UserRepository) Create(user User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = user
}

func (r *UserRepository) Get(id string) (
	User,
	bool,
) {

	r.mu.RLock()

	defer r.mu.RUnlock()

	user, ok :=
		r.users[id]

	return user, ok
}

func (r *UserRepository) Delete(id string) bool {

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[id]; !ok {
		return false
	}

	delete(r.users, id)

	return true
}

func (r *UserRepository) Deposit(
	user_id string,
	asset string,
	amount string,
) (
	Balance,
	bool,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[user_id]
	if !ok {
		return Balance{}, false
	}

	if user.Balances == nil {
		user.Balances = make(map[string]*Balance)
	}

	userbalance, ok := user.Balances[asset]
	if !ok {
		userbalance = &Balance{
			Asset:     asset,
			Available: "0",
			Locked:    "0",
		}
	}

	userbalance.Available += amount

	user.Balances[asset] = userbalance
	r.users[user_id] = user

	return *userbalance, true
}
