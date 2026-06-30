package application

import (
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"time"

	"go.uber.org/zap"
)

func (u *User) GetUser(ctx context.Context,
	id string,
) (
	user.User,
	error,
) {
	log := u.log.Named("GetUser")
	log.Info("getting user", zap.String("id", id))

	currentUser, err := u.userProvider.Get(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return user.User{}, err
	}

	log.Info("got user", zap.String("id", id))

	return currentUser, nil
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
	user user.User,
	err error,
) {
	panic("implement me")
}

func (u *User) UpdateUserInfo(ctx context.Context,
	id string,
	name string,
	email string,
) (
	updated bool,
	updatedUserInfoAt time.Time,
	err error,
) {
	panic("implement me")
}

func (u *User) ChangePassword(ctx context.Context,
	id string,
	oldPassword string,
	newPassword string,
) (
	success bool,
	userPasswordChangedAt time.Time,
) {
	panic("implement me")
}
