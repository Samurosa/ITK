package inmemory

import (
	"ITK_Code/m/v2/internal/adapters/storage"
	models2 "ITK_Code/m/v2/internal/core/user/models"
	"ITK_Code/m/v2/internal/domain/models"
	"bytes"
	"context"
	"time"

	"github.com/google/uuid"

	"sync"
)

type UserRepository struct {
	mu sync.RWMutex

	users map[string]*models2.User

	usersLoginById map[string]string
}

func NewUserStorage() *UserRepository {
	return &UserRepository{
		users: make(map[string]*models2.User),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context,
	login string,
	passwordHash []byte,
	name string,
	balances map[string]*models2.Balance,
	role models2.Role,
	createTime time.Time,
	updateTime time.Time,
) (
	string,
	error,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	user := models2.User{
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

func (r *UserRepository) GetUser(ctx context.Context,
	uid string,
) (
	models2.User,
	error,
) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[uid]
	if !ok {
		return &models2.User{}, storage.ErrUserNotFound
	}

	return user, nil
}
func (r *UserRepository) GetUserByLogin(ctx context.Context,
	login string,
) (
	models2.User,
	error,
) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user := models2.User{}

	for _, userInRep := range r.users {
		if userInRep.Login == login {
			user = *userInRep
		}
	}

	if user.ID == "" {
		return &models2.User{}, storage.ErrUserNotFound
	}

	return &user, nil
}
func (r *UserRepository) GetBalanceUser(ctx context.Context,
	uid string,
	asset string,
) (
	models2.Balance,
	error,
) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[uid]
	if !ok {
		return nil, storage.ErrUserNotFound
	}

	return user.Balances[asset], nil
}

func (r *UserRepository) UpdateUser(ctx context.Context,
	user *models2.User,
	update models.Update,
) (
	bool,
	error,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var nilByte []byte
	if update.Name != "" {
		user.Name = update.Name
	}
	if update.Login != "" {
		user.Login = update.Login
	}
	if bytes.Contains(update.Password, nilByte) {
		user.PasswordHash = update.Password
	}

	return true, nil
}

func (r *UserRepository) DeleteUser(ctx context.Context,
	uid string,
) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[uid]; !ok {
		return storage.ErrUserNotFound
	}

	delete(r.users, uid)
	return nil
}
func (r *UserRepository) IsAdmin(ctx context.Context,
	uid string,
) (
	bool,
	error,
) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[uid]
	if !ok {
		return false, storage.ErrUserNotFound
	}

	if user.Role != models2.AdminRole {
		return false, nil
	}

	return true, nil
}

// возвращать только копии
// переходим на постргрес
