package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/user"
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.uber.org/zap"
)

type User struct {
	userRepository ports.UserRepository

	sessionStorage ports.SessionRepository

	passwordHasher  ports.PasswordHasher
	passwordLimiter ports.PasswordChangeLimiter
}

func NewUserService(
	userRepository ports.UserRepository,
	sessionStorage ports.SessionRepository,
	passwordHasher ports.PasswordHasher,
	passwordLimiter ports.PasswordChangeLimiter,
) *User {
	return &User{
		userRepository:  userRepository,
		sessionStorage:  sessionStorage,
		passwordHasher:  passwordHasher,
		passwordLimiter: passwordLimiter,
	}
}

func (u *User) GetUser(ctx context.Context,
	id string,
) (
	user.User,
	error,
) {
	log := logging.FromContext(ctx).Named("GetUser").With(zap.String("user_id", id))

	current, err := u.userRepository.Get(ctx, id)
	if err != nil {
		logOperationError(log, "failed to get user", err)
		return user.User{}, err
	}
	log.Debug("user retrieved")

	return current, nil
}

func (u *User) DeleteUser(ctx context.Context,
	id string,
) error {
	log := logging.FromContext(ctx).Named("DeleteUser").With(zap.String("user_id", id))

	err := u.userRepository.Delete(ctx, id)
	if err != nil {
		logOperationError(log, "failed to delete user", err)
		return err
	}
	log.Debug("user marked as deleted")

	err = u.sessionStorage.DeleteByUser(ctx, id)
	if err != nil {
		log.Error("user deleted but session revocation failed", zap.Error(err))
		return err
	}
	log.Info("user deleted and sessions revoked")

	return nil
}

func (u *User) IsAdmin(ctx context.Context,
	id string,
) (
	bool,
	error,
) {
	log := logging.FromContext(ctx).Named("IsAdmin").With(zap.String("user_id", id))

	isAdmin, err := u.userRepository.IsAdmin(ctx, id)
	if err != nil {
		logOperationError(log, "failed to check user role", err)
		return false, err
	}
	log.Debug("user role checked", zap.Bool("is_admin", isAdmin))

	return isAdmin, nil
}

func (u *User) GetUserByEmail(ctx context.Context,
	email string,
) (
	user.User,
	error,
) {
	log := logging.FromContext(ctx).Named("GetUserByEmail")

	current, err := u.userRepository.GetByEmail(ctx, email)
	if err != nil {
		logOperationError(log, "failed to get user by email", err)
		return current, err
	}
	log.Debug("user retrieved by email", zap.String("user_id", current.ID))

	return current, nil
}

func (u *User) UpdateUserInfo(ctx context.Context,
	id string,
	name string,
) error {
	log := logging.FromContext(ctx).Named("UpdateUserInfo").With(zap.String("user_id", id))

	updated := user.UpdateUser{}
	if name != "" {
		updated.Name = &name
	}

	err := u.userRepository.Update(ctx, id, updated)
	if err != nil {
		logOperationError(log, "failed to update user", err)
		return err
	}
	log.Info("user updated")

	return nil
}

func (u *User) ChangePassword(ctx context.Context,
	id string,
	oldPassword string,
	newPassword string,
) error {
	log := logging.FromContext(ctx).Named("ChangePassword").With(zap.String("user_id", id))
	allowed, err := u.passwordLimiter.AllowPasswordChange(ctx, id)
	if err != nil {
		logOperationError(log, "failed to check password change rate limit", err)
		return corerrors.ErrTooManyRequests
	}
	if !allowed {
		log.Warn("password change rate limit exceeded")
		return corerrors.ErrTooManyRequests
	}

	current, err := u.userRepository.Get(ctx, id)
	if err != nil {
		logOperationError(log, "failed to get user", err)
		return err
	}
	log.Debug("user retrieved for password change")

	err = u.passwordHasher.VerifyPasswordHash(oldPassword, current.PasswordHash)
	if err != nil {
		logPasswordVerificationError(log, err)
		return auth.ErrIncorrectPassword
	}
	log.Debug("current password verified")

	newPassHash, err := u.passwordHasher.GeneratePasswordHash(newPassword)
	if err != nil {
		log.Error("failed to generate password hash", zap.Error(err))
		return corerrors.ErrPassGenHash
	}
	log.Debug("password hash generated")

	err = u.userRepository.UpdatePassword(ctx, current, string(newPassHash))
	if err != nil {
		logOperationError(log, "failed to update password", err)
		return err
	}
	log.Debug("password updated; revoking sessions")

	if err := u.sessionStorage.DeleteByUser(ctx, id); err != nil {
		log.Error("password updated but session revocation failed", zap.Error(err))
		return err
	}
	log.Info("password changed and sessions revoked")

	return nil
}
