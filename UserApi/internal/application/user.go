package application

import (
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) GetUser(ctx context.Context,
	id string,
) (
	user.User,
	error,
) {
	log := u.log.Named("GetUser")
	log.Info("getting user", zap.String("id", id))

	current, err := u.userProvider.Get(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return user.User{}, err
	}

	log.Info("got user", zap.String("id", id))

	return current, nil
}

func (u *User) DeleteUser(ctx context.Context,
	id string,
) (
	bool,
	time.Time,
	error,
) {
	log := u.log.Named("DeleteUser")
	log.Info("deleting user", zap.String("id", id))

	err := u.userProvider.Delete(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, time.Time{}, err
	}

	log.Info("user deleted", zap.String("id", id))

	return true, time.Now(), nil
}

func (u *User) IsAdmin(ctx context.Context,
	uid string,
) (
	bool,
	error,
) {
	log := u.log.Named("IsAdmin")
	log.Info("checking user is admin", zap.String("id", uid))

	isAdmin, err := u.userProvider.IsAdmin(ctx, uid)
	if err != nil {
		log.Error("user not found", zap.String("id", uid), zap.Error(err))
		return false, err
	}

	return isAdmin, nil
}

func (u *User) GetUserByEmail(ctx context.Context,
	email string,
) (
	user.User,
	error,
) {
	log := u.log.Named("GetUserByEmail")

	current, err := u.userProvider.GetByEmail(ctx, email)
	if err != nil {
		log.Error("user not found", zap.String("email", email), zap.Error(err))
		return current, err
	}
	return current, nil
}

func (u *User) UpdateUserInfo(ctx context.Context,
	id string,
	name string,
	email string,
) (
	bool,
	time.Time,
	error,
) {
	log := u.log.Named("update user")

	updated := user.UpdateUser{}
	if name != "" {
		updated.Name = &name
	}
	if email != "" {
		updated.Email = &email
	}

	success, err := u.userProvider.Update(ctx, id, updated)
	if err != nil {
		log.Error("error updating user", zap.Error(user.ErrUpdateUser))
		return false, time.Time{}, user.ErrUpdateUser
	}

	return success, time.Now(), nil
}

func (u *User) ChangePassword(ctx context.Context,
	id string,
	oldPassword string,
	newPassword string,
) (
	bool,
	time.Time,
	error,
) {
	log := u.log.Named("change Password")
	log.Info("searching user with id", zap.String("id", id))

	current, err := u.userProvider.Get(ctx, id)
	if err != nil {
		log.Error("error getting user by email", zap.String("id", id), zap.Error(user.ErrUserNotFound))
		return false, time.Time{}, user.ErrUserNotFound
	}

	err = authorization(current.PasswordHash, oldPassword)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(user.ErrComparePassword))
		return false, time.Time{}, user.ErrComparePassword
	}

	newPassHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return false, time.Time{}, err
	}

	updated := user.UpdateUser{}
	updated.PassHash = &newPassHash

	success, err := u.userProvider.UpdatePassword(ctx, current, updated)
	if err != nil {
		log.Error("error updating user", zap.Error(user.ErrUpdateUser))
		return false, time.Time{}, user.ErrUpdateUser
	}

	return success, time.Now(), nil
}
