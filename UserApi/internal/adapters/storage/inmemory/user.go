package inmemory

import (
	userCore "ITK_Code/m/v2/internal/core/user"
	"context"

	"github.com/google/uuid"

	"sync"
)

type UserRepository struct {
	mu sync.RWMutex

	users map[string]*userCore.User

	usersLoginById map[string]string
}

func NewUserStorage() *UserRepository {
	return &UserRepository{
		users: make(map[string]*userCore.User),
	}
}

func (r *UserRepository) SaveUser(ctx context.Context,
	user userCore.User,
) (
	string,
	error,
) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id := uuid.New().String()
	user.ID = id

	r.users[id] = &user

	return id, nil
}

func (r *UserRepository) IsExistsUserByEmail(ctx context.Context,
	login string,
) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	isUserExist := false

	for _, userInRep := range r.users {
		if userInRep.Email == login {
			isUserExist = true
		}
	}

	return isUserExist
}

func (r *UserRepository) Get(ctx context.Context,
	uid string,
) (
	userCore.User,
	error,
) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user, ok := r.users[uid]
	if !ok {
		return userCore.User{}, userCore.ErrUserNotFound
	}

	return *user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context,
	email string,
) (
	userCore.User,
	error,
) {

	r.mu.RLock()
	defer r.mu.RUnlock()

	user := userCore.User{}

	for _, userInRep := range r.users {
		if userInRep.Email == email {
			user = *userInRep
		}
	}

	if user.ID == "" {
		return userCore.User{}, userCore.ErrUserNotFound
	}

	return user, nil
}

func (r *UserRepository) Update(ctx context.Context, user *userCore.User, update userCore.UpdateUser) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	panic("implement me")
}

func (r *UserRepository) Delete(ctx context.Context,
	uid string,
) error {

	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[uid]; !ok {
		return userCore.ErrUserNotFound
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
		return false, userCore.ErrUserNotFound
	}

	if user.Role != userCore.AdminRole {
		return false, nil
	}

	return true, nil
}

// возвращать только копии
// переходим на постргрес
