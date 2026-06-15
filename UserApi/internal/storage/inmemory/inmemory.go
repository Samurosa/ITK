package inmemory

import (
	"ITK_Code/m/v2/internal/domain/models"
	"ITK_Code/m/v2/internal/storage"
	"context"
	"time"

	"github.com/google/uuid"

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

func (r *UserRepository) SaveUser(ctx context.Context,
	login string,
	passwordHash []byte,
	name string,
	balances map[string]*models.Balance,
	role models.Role,
	createTime time.Time,
	updateTime time.Time,
) (
	string,
	error,
) {
	defer ctx.Done()
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	user := models.User{
		ID:           id,
		Name:         name,
		Login:        login,
		PasswordHash: passwordHash,
		Balances:     balances,
		Role:         role,
		CreateTime:   createTime,
		UpdateTime:   updateTime,
	}

	r.users[id] = &user

	return id, nil
}

func (r *UserRepository) IsExistsUserByLogin(ctx context.Context,
	login string,
) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	isUserExist := false

	for _, userInRep := range r.users {
		if userInRep.Login == login {
			isUserExist = true
		}
	}

	return isUserExist
}

func (r *UserRepository) GetUser(ctx context.Context, uid string) (models.User, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[uid]
	if !ok {
		return models.User{}, storage.ErrUserNotFound
	}

	return *user, nil
}
func (r *UserRepository) GetUserByLogin(ctx context.Context, login string) (models.User, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user := models.User{}

	for _, userInRep := range r.users {
		if userInRep.Login == login {
			user = *userInRep
		}
	}

	if user.ID == "" {
		return user, storage.ErrUserNotFound
	}

	return user, nil
}
func (r *UserRepository) GetBalanceUser(ctx context.Context, uid string, asset string) (*models.Balance, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[uid]
	if !ok {
		return nil, storage.ErrUserNotFound
	}

	return user.Balances[asset], nil
}
func (r *UserRepository) DeleteUser(ctx context.Context, uid string) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[uid]; !ok {
		return storage.ErrUserNotFound
	}

	delete(r.users, uid)
	return nil
}
func (r *UserRepository) IsAdmin(ctx context.Context, uid string) (bool, error) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[uid]
	if !ok {
		return false, storage.ErrUserNotFound
	}

	if user.Role != "ROLE_ADMIN" {
		return false, nil
	}

	return true, nil
}
