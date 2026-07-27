package application

import (
	context2 "ITK_Code/m/v2/internal/adapters/outbound/context"
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/hash"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
	"context"
	"time"

	"go.uber.org/zap"
)

func (u *User) GetUser(ctx context.Context,
) (
	dto.User,
	error,
) {
	log := u.log.Named("GetUser")

	id, err := context2.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return dto.User{}, errors.ErrInvalidContext
	}
	log.Info("user id from context", zap.String("id", id))

	current, err := u.userProvider.Get(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return dto.User{}, errors.ErrUserNotFound
	}
	log.Info("got user id:", zap.String("id", id))

	return current, nil
}

func (u *User) DeleteUser(ctx context.Context,
) (
	bool,
	time.Time,
	error,
) {
	log := u.log.Named("DeleteUser")

	id, err := context2.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return false, time.Time{}, errors.ErrInvalidContext
	}
	log.Info("user id from context", zap.String("id", id))

	err = u.userProvider.Delete(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, time.Time{}, errors.ErrUserNotFound
	}
	log.Info("user deleted", zap.String("id", id))

	err = u.sessionStorage.DeleteByUser(ctx, id)
	if err != nil {
		log.Error("session not found", zap.String("id", id), zap.Error(err))
		return false, time.Time{}, errors.ErrUserNotFound
	}
	log.Info("user sessions deleted", zap.String("id", id))

	return true, time.Now(), nil
}

func (u *User) IsAdmin(ctx context.Context,
) (
	bool,
	error,
) {
	log := u.log.Named("IsAdmin")

	id, err := context2.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return false, errors.ErrInvalidContext
	}
	log.Info("user id from context", zap.String("id", id))

	isAdmin, err := u.userProvider.IsAdmin(ctx, id)
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
	dto.User,
	error,
) {
	log := u.log.Named("GetUserByEmail")

	current, err := u.userProvider.GetByEmail(ctx, email)
	if err != nil {
		log.Error("user not found", zap.String("email", email), zap.Error(err))
		return current, err
	}
	log.Info("got user by email", zap.String("id", current.ID))

	return current, nil
}

func (u *User) UpdateUserInfo(ctx context.Context,
	name string,
	email string,
) (
	bool,
	time.Time,
	error,
) {
	log := u.log.Named("update user")

	id, err := context2.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return false, time.Time{}, errors.ErrInvalidContext
	}
	log.Info("user id from context", zap.String("id", id))

	updated := dto.UpdateUser{}
	if name != "" {
		updated.Name = &name
	}
	if email != "" {
		updated.Email = &email
	}

	success, err := u.userProvider.Update(ctx, id, updated)
	if err != nil {
		log.Error("error updating user", zap.Error(err))
		return false, time.Time{}, errors.ErrUpdateUser
	}
	log.Info("user updated", zap.String("id", id))

	return success, time.Now(), nil
}

func (u *User) ChangePassword(ctx context.Context,
	oldPassword string,
	newPassword string,
) (
	bool,
	time.Time,
	error,
) {
	log := u.log.Named("change Password")

	id, err := context2.GetUserIDFromContext(ctx)
	if err != nil {
		log.Error("context is not valid", zap.Error(err))
		return false, time.Time{}, errors.ErrInvalidContext
	}
	log.Info("user id from context", zap.String("id", id))

	current, err := u.userProvider.Get(ctx, id)
	if err != nil {
		log.Error("error getting user", zap.String("id", id), zap.Error(err))
		return false, time.Time{}, errors.ErrUserNotFound
	}
	log.Info("health check user successful, got user:", zap.String("id", current.ID))

	err = hash.VerifyPasswordHash(oldPassword, current.PasswordHash)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(err))
		return false, time.Time{}, errors.ErrInvalidLoginCredentials
	}
	log.Info("verify password successful")

	newPassHash, err := hash.GeneratePasswordHash(newPassword)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return false, time.Time{}, errors.ErrPassGenHash
	}
	log.Info("password hash generated")

	success, err := u.userProvider.UpdatePassword(ctx, current, string(newPassHash))
	if err != nil {
		log.Error("error updating user", zap.Error(err))
		return false, time.Time{}, errors.ErrUpdateUser
	}
	log.Info("success updated password", zap.String("id", current.ID))

	return success, time.Now(), nil
}
