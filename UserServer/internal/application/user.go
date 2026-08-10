package application

import (
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/hash"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/errors"
	"ITK_Code/m/v2/internal/core/user"
	"context"

	"go.uber.org/zap"
)

func (u *User) GetUser(ctx context.Context,
	id string,
) (
	user.User,
	error,
) {
	log := u.log.Named("GetUser")

	current, err := u.userRepository.Get(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return user.User{}, user.ErrUserNotFound
	}
	log.Info("got user id:", zap.String("id", id))

	return current, nil
}

func (u *User) DeleteUser(ctx context.Context,
	id string,
) error {
	log := u.log.Named("DeleteUser")

	err := u.userRepository.Delete(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return user.ErrUserNotFound
	}
	log.Debug("user deleted", zap.String("id", id))

	err = u.sessionStorage.DeleteByUser(ctx, id)
	if err != nil {
		log.Error("session not found", zap.String("id", id), zap.Error(err))
		return user.ErrUserNotFound
	}
	log.Info("user deleted", zap.String("id", id))

	return nil
}

func (u *User) IsAdmin(ctx context.Context,
	id string,
) (
	bool,
	error,
) {
	log := u.log.Named("IsAdmin")

	isAdmin, err := u.userRepository.IsAdmin(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, err
	}
	log.Info("checked role user", zap.Bool("admin role", isAdmin))

	return isAdmin, nil
}

func (u *User) GetUserByEmail(ctx context.Context,
	email string,
) (
	user.User,
	error,
) {
	log := u.log.Named("GetUserByEmail")

	current, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		log.Error("user not found", zap.String("email", email), zap.Error(err))
		return current, err
	}
	log.Info("got user by email", zap.String("id", current.ID))

	return current, nil
}

func (u *User) UpdateUserInfo(ctx context.Context,
	id string,
	name string,
	email string,
) error {
	log := u.log.Named("update user")

	updated := user.UpdateUser{}
	if name != "" {
		updated.Name = &name
	}
	if email != "" {
		updated.Email = &email
	}

	err := u.userRepository.Update(ctx, id, updated)
	if err != nil {
		log.Error("error updating user", zap.Error(err))
		return user.ErrUpdateUser
	}
	log.Info("user updated", zap.String("id", id))

	return nil
}

func (u *User) ChangePassword(ctx context.Context,
	id string,
	oldPassword string,
	newPassword string,
) error {
	log := u.log.Named("change Password")

	current, err := u.userRepository.Get(ctx, id)
	if err != nil {
		log.Error("error getting user", zap.String("id", id), zap.Error(err))
		return user.ErrUserNotFound
	}
	log.Debug("health check user successful, got user:", zap.String("id", current.ID))

	err = hash.VerifyPasswordHash(oldPassword, current.PasswordHash)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(err))
		return auth.ErrIncorrectPassword
	}
	log.Debug("verify password successful")

	newPassHash, err := hash.GeneratePasswordHash(newPassword)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return errors.ErrPassGenHash
	}
	log.Debug("password hash generated")

	err = u.userRepository.UpdatePassword(ctx, current, string(newPassHash))
	if err != nil {
		log.Error("error updating user", zap.Error(err))
		return user.ErrUpdateUser
	}
	log.Info("success updated password", zap.String("id", current.ID))

	return nil
}
