package application

import (
	tokenGenerate "ITK_Code/m/v2/internal/adapters/jwt"
	"ITK_Code/m/v2/internal/adapters/storage"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user/models"
	"context"
	"fmt"
	"time"

	"ITK_Code/m/v2/internal/core/user"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	log *zap.SugaredLogger

	jwtConfig auth.JWTConfig

	sessionStorage auth.SessionRepository

	userSaver    user.Save
	userProvider user.Service
}

func New(
	log *zap.SugaredLogger,
	jwtConfig auth.JWTConfig,
	sessionStorage auth.SessionRepository,
	userSaver user.Save,
	userProvider user.Service,
) *User {
	return &User{
		log:            log,
		jwtConfig:      jwtConfig,
		sessionStorage: sessionStorage,
		userSaver:      userSaver,
		userProvider:   userProvider,
	}
}
func (u *User) Registration(ctx context.Context,
	email string,
	password string,
	name string,
) (
	id string,
	tokenPairs auth.TokensModel,
	createdAt time.Time,
	err error,
) {
	log := u.log.Named("RegisterNewUser")
	log.Info("registering new user")

	if u.userSaver.IsExistsUserByEmail(ctx, email) {
		log.Info("user with email already exists")
		return "", auth.TokensModel{}, time.Time{}, user.ErrUserExists
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, err
	}

	defaultUserName := "default name"

	balances := map[string]*models2.Balance{}
	balances["USD"] = &models2.Balance{
		Asset:     "USD",
		Available: "0",
		Locked:    "0",
	}

	now := time.Now()

	uid, err := u.userSaver.SaveUser(ctx, email, passHash, defaultUserName, balances, models2.UserRole, now, now)
	if err != nil {
		log.Error("error saving user", zap.Error(err))
		return "", time.Time{}, err
	}

	log.Info("user created", zap.String("uid", uid))

	return uid, now, nil
}

func (u *User) email(ctx context.Context,
	email string,
	password string,
) (
	tokensPairs auth.TokensModel,
	err error,
) {
	panic("implement me")
}

func (u *User) Logout(ctx context.Context,
	refreshToken string,
	deviceID string,
) (
	success bool,
	loggedOutAt time.Time,
	err error,
) {
	panic("implement me")
}

func (u *User) RefreshToken(ctx context.Context,
	refreshToken string,
	deviceID string,
) (
	tokensPairs auth.TokensModel,
	err error,
) {
	panic("implement me")
}

func (u *User) GetUser(ctx context.Context,
	id string,
) (
	user models2.User,
	err error,
) {
	log := u.log.Named("GetUser")
	log.Info("getting user", zap.String("id", id))

	user, err := u.userProvider.GetUser(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return models2.User{}, err
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

	user, err := u.userProvider.GetUser(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, time.Time{}, err
	}

	passHash := []byte(password)
	if password != "" {
		passHash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("error generating password hash", zap.Error(err))
			return false, time.Time{}, err
		}
	}

	updateUser := models.Update{
		Name:     name,
		email:    email,
		Password: passHash,
	}

	done, err := u.userProvider.UpdateUser(ctx, user, updateUser)
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
	error,
) {
	log := u.log.Named("DeleteUser")
	log.Info("deleting user", zap.String("id", id))

	err := u.userProvider.DeleteUser(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, err
	}

	log.Info("user deleted", zap.String("id", id))

	return true, nil
}

func (u *User) Deposit(ctx context.Context,
	id string,
	asset string,
	amount string,
) (
	bool,
	*models2.Balance,
	error,
) {
	log := u.log.Named("Deposit")
	log.Info("depositing user", zap.String("id", id))

	balance, err := u.userProvider.GetBalanceUser(ctx, id, asset)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return false, nil, err
	}

	balance.Available += amount
	balance.Locked += amount

	log.Info("deposit is successful", zap.String("id", id))

	return true, balance, nil
}

func (u *User) Authorization(ctx context.Context,
	email string,
	password string,
) (
	string,
	error,
) {
	log := u.log.Named("Authorization")
	log.Info("authorizing user", zap.String("email", email))

	user, err := u.userProvider.GetUserByemail(ctx, email)
	if err != nil {
		log.Error("email failed", zap.Error(err))
		return "", storage.ErrUserNotFound
	}

	secret, err := u.app.GetSecret()
	if err != nil {
		log.Error("secret failed", zap.Error(err))
		return "", storage.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Error("email failed", zap.Error(err))
		return "", storage.ErrUserNotFound
	}

	log.Info("user is authorized", zap.String("email", email))

	token, err := tokenGenerate.NewToken(*user, secret, u.tokenTTL)
	if err != nil {
		log.Error("error generating token", zap.Error(err))
		return "", err
	}

	return token, nil
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
