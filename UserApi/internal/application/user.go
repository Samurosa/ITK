package application

import (
	userModels "ITK_Code/m/v2/internal/core/user/models"
	"context"
	"time"

	"go.uber.org/zap"
)

func (u *User) GetUser(ctx context.Context,
	id string,
) (
	user userModels.User,
	err error,
) {
	log := u.log.Named("GetUser")
	log.Info("getting user", zap.String("id", id))

	user, err = u.userProvider.Get(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return userModels.User{}, err
	}

	log.Info("got user", zap.String("id", id))

	return user, nil
}

func (u *User) UpdateUser(ctx context.Context,
	id string,
	name string,
	email string,
	password string,
) (
	updated bool,
	updatedAt time.Time,
	err error,
) {
	log := u.log.Named("UpdateUser")
	log.Info("updating user", zap.String("id", id))

	user, err := u.userProvider.Get(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, time.Time{}, err
	}

	updateUser := userModels.UpdateUser{
		Name:  name,
		Email: email,
	}

	done, err := u.userProvider.Update(ctx, &user, updateUser)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, time.Time{}, err
	}

	log.Info("user update", zap.String("id", id))

	return done, updatedAt, err
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

func (u *User) Deposit(ctx context.Context,
	id string,
	asset string,
	amount userModels.Money,
) (
	bool,
	userModels.Balance,
	error,
) {
	log := u.log.Named("Deposit")
	log.Info("depositing user", zap.String("id", id))

	balance, err := u.userProvider.GetBalance(ctx, id, asset)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, userModels.Balance{}, err
	}

	_ = depositLock(true, balance, amount)

	newBalance, err := u.userProvider.Deposit(ctx, balance, amount)
	if err != nil {
		_ = depositLock(false, balance, amount)
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, userModels.Balance{}, err
	}
	//добавить персонализованную ошибку под баланс

	log.Info("deposit is successful", zap.String("id", id))

	return true, newBalance, nil
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
	user userModels.User,
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

func (u *User) GetBalance(ctx context.Context,
	id string,
) (
	balances map[string]userModels.Balance,
	err error,
) {
	panic("implement me")
}

func depositLock(deposit bool, target userModels.Balance, amount userModels.Money) error {
	if deposit {
		target.Asset = amount.Currency

		target.Locked.Currency = amount.Currency
		target.Locked.Units += amount.Units
		target.Locked.Nanos += amount.Nanos
		return nil
	}

	target.Asset = amount.Currency

	target.Locked.Currency = amount.Currency
	target.Locked.Units -= amount.Units
	target.Locked.Nanos -= amount.Nanos

	return nil
}
