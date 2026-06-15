package services

import (
	"ITK_Code/m/v2/internal/storage"
	"context"
	"time"

	"ITK_Code/m/v2/internal/domain/models"
	tokenGenerate "ITK_Code/m/v2/internal/lib/jwt"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	log          *zap.Logger
	userSaver    UserSave
	userProvider UserProvider
	app          AppProvider
	tokenTTL     time.Duration
}

type UserSave interface {
	SaveUser(ctx context.Context,
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
	)
	IsExistsUserByLogin(ctx context.Context,
		login string,
	) bool
}

type UserProvider interface {
	GetUser(ctx context.Context, uid string) (*models.User, error)
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
	GetBalanceUser(ctx context.Context, uid string, asset string) (*models.Balance, error)
	UpdateUser(ctx context.Context, user *models.User, name string, login string, password []byte) (bool, error)
	DeleteUser(ctx context.Context, uid string) error
	IsAdmin(ctx context.Context, uid string) (bool, error)
}

type AppProvider interface {
	GetSecret() (models.App, error)
}

func New(
	log *zap.Logger,
	userSaver UserSave,
	userProvider UserProvider,
	app AppProvider,
	tokenTTL time.Duration,
) *User {
	return &User{
		log:          log,
		userSaver:    userSaver,
		userProvider: userProvider,
		app:          app,
		tokenTTL:     tokenTTL,
	}
}
func (u *User) Registration(ctx context.Context,
	login string,
	password string,
) (
	id string,
	createdAt time.Time,
	err error,
) {
	log := u.log.Named("RegisterNewUser")
	log.Info("registering new user")

	if u.userSaver.IsExistsUserByLogin(ctx, login) {
		log.Info("user with login already exists")
		return "", time.Time{}, storage.ErrUserExists
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, err
	}

	defaultUserName := "default name"

	balances := map[string]*models.Balance{}
	balances["USD"] = &models.Balance{
		Asset:     "USD",
		Available: "0",
		Locked:    "0",
	}

	now := time.Now()

	uid, err := u.userSaver.SaveUser(ctx, login, passHash, defaultUserName, balances, models.UserRole, now, now)
	if err != nil {
		log.Error("error saving user", zap.Error(err))
		return "", time.Time{}, err
	}

	log.Info("user created", zap.String("uid", uid))

	return uid, now, nil
}

func (u *User) GetUser(ctx context.Context,
	id string,
) (
	name string,
	login string,
	balances map[string]*models.Balance,
	role models.Role,
	createdAt time.Time,
	updatedAt time.Time,
	err error,
) {
	log := u.log.Named("GetUser")
	log.Info("getting user", zap.String("id", id))

	user, err := u.userProvider.GetUser(ctx, id)
	if err != nil {
		log.Error("user not found", zap.String("id", id), zap.Error(err))
		return "", "", nil, "", time.Time{}, time.Time{}, err
	}

	log.Info("got user", zap.String("id", id))

	return user.Name, user.Login, user.Balances, user.Role, user.CreateTime, user.UpdateTime, nil
}

func (u *User) UpdateUser(ctx context.Context,
	id string,
	name string,
	login string,
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

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return false, time.Time{}, err
	}

	done, err := u.userProvider.UpdateUser(ctx, user, name, login, passHash)
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
	*models.Balance,
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
	login string,
	password string,
) (
	string,
	error,
) {
	log := u.log.Named("Authorization")
	log.Info("authorizing user", zap.String("login", login))

	user, err := u.userProvider.GetUserByLogin(ctx, login)
	if err != nil {
		log.Error("login failed", zap.Error(err))
		return "", storage.ErrUserNotFound
	}

	secret, err := u.app.GetSecret()
	if err != nil {
		log.Error("secret failed", zap.Error(err))
		return "", storage.ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(password)); err != nil {
		log.Error("login failed", zap.Error(err))
		return "", storage.ErrUserNotFound
	}

	log.Info("user is authorized", zap.String("login", login))

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
