package storage

import (
	models "ITK_Code/m/v2/internal/domain/models"
	"sync"
)

type UserRepository struct {
	mu sync.RWMutex

	users map[string]*models.User
}

func NewUserStorage() *UserRepository {
	return &UserRepository{
		users: make(map[string]*models.User),
	}
}

func (r *UserRepository) Create(user models.User) {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.users[user.ID] = &user
}

func (r *UserRepository) Get(id string) (
	*models.User,
	bool,
) {

	r.mu.RLock()

	defer r.mu.RUnlock()

	user, ok :=
		r.users[id]

	return user, ok
}

func (r *UserRepository) Update(
	id string,
	name *string,
	login *string,
	password *string,
) (
	string,
	bool,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return "", false
	}

	if name != nil {
		user.Name = *name
	}

	if login != nil {
		user.Login = *login
	}

	return id, true
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
	models.Balance,
	bool,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[user_id]
	if !ok {
		return models.Balance{}, false
	}

	if user.Balances == nil {
		user.Balances = make(map[string]*models.Balance)
	}

	userbalance, ok := user.Balances[asset]
	if !ok {
		userbalance = &models.Balance{
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
