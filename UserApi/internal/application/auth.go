package application

import (
	authCore "ITK_Code/m/v2/internal/core/auth"
	userCore "ITK_Code/m/v2/internal/core/user"
	userModels "ITK_Code/m/v2/internal/core/user/models"
	"context"
	"time"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func (u *User) Registration(ctx context.Context,
	email string,
	password string,
	name string,
) (
	string,
	time.Time,
	error,
) {
	now := time.Now()

	log := u.log.Named("RegisterNewUser")

	if u.userSaver.IsExistsUserByEmail(ctx, email) {
		log.Info("user with email already exists")
		return "", time.Time{}, userCore.ErrUserExists
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, err
	}

	newUser := userModels.User{
		Email:        email,
		Name:         name,
		PasswordHash: passHash,
		Balances:     GetDefaultBalance(),
		Role:         userModels.UserRole,
		CreateTime:   now,
		UpdateTime:   now,
	}

	uid, err := u.userSaver.SaveUser(ctx, newUser)
	if err != nil {
		log.Error("error saving user", zap.Error(err))
		return "", time.Time{}, err
	}
	log.Info("user created", zap.String("uid", uid), zap.String("createdAt", now.String()))

	return uid, now, nil
}

func (u *User) Login(ctx context.Context,
	email string,
	password string,
	deviceId string,
) (
	tokensPairs authCore.TokensModel,
	err error,
) {
	log := u.log.Named("login")
	log.Info("searching user with email", zap.String("email", email))

	user, err := u.userProvider.GetByEmail(ctx, email)
	if err != nil {
		log.Error("error getting user by email", zap.String("email", email), zap.Error(userCore.ErrUserNotFound))
		return authCore.TokensModel{}, userCore.ErrUserNotFound
	}

	err = authorization(user.PasswordHash, password)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	log.Info("search old session")
	_ = u.sessionStorage.DeleteByUserAndDevice(ctx, user.ID, deviceId)

	log.Info("create token")
	tokens, err := u.tokenManager.Generate(user)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	tokenHash, err := bcrypt.GenerateFromPassword([]byte(tokens.RefreshToken), bcrypt.DefaultCost)

	log.Info("session info save")
	session := authCore.SessionModel{
		UserID:           user.ID,
		DeviceID:         deviceId,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        tokens.RefreshExpiresAt,
	}

	err = u.sessionStorage.Create(ctx, session)
	if err != nil {
		log.Error("error creating session", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	return tokens, nil
}

func (u *User) Logout(ctx context.Context,
	refreshToken string,
	userId string,
	deviceID string,
) (
	success bool,
	loggedOutAt time.Time,
	err error,
) {
	now := time.Now()
	log := u.log.Named("Logout")
	sessionInfo, err := u.sessionStorage.GetByUserIdAndDeviceId(ctx, userId, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, authCore.Unauthorized
	}

	err = bcrypt.CompareHashAndPassword(sessionInfo.RefreshTokenHash, []byte(refreshToken))
	if err != nil {
		log.Error("error comparing refresh token", zap.Error(err))
		return false, time.Time{}, authCore.ErrNoAccess
	}

	err = u.sessionStorage.DeleteByUserAndDevice(ctx, userId, deviceID)
	if err != nil {
		log.Error("error deleting session", zap.Error(err))
		return false, time.Time{}, err
	}

	return true, now, nil
}

func (u *User) RefreshToken(ctx context.Context,
	refreshToken string,
	userId string,
	deviceID string,
) (
	tokensPairs authCore.TokensModel,
	err error,
) {
	log := u.log.Named("Refresh tokens")

	log.Info("searching user with db", zap.String("userId", userId))
	user, err := u.userProvider.Get(ctx, userId)
	if err != nil {
		log.Error("error getting user by id", zap.Error(err))
		return authCore.TokensModel{}, userCore.ErrUserNotFound
	}

	log.Info("searching session with userId deviceID", zap.String("userId", userId), zap.String("deviceID", deviceID))
	sessionInfo, err := u.sessionStorage.GetByUserIdAndDeviceId(ctx, userId, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return authCore.TokensModel{}, authCore.Unauthorized
	}

	log.Info("validating token expiry", zap.String("ExpireAt", sessionInfo.ExpiresAt.String()))
	if sessionInfo.ExpiresAt.Before(time.Now()) {
		return authCore.TokensModel{}, authCore.ErrRefreshExpired
	}

	log.Info("comparing refresh token")
	err = bcrypt.CompareHashAndPassword(sessionInfo.RefreshTokenHash, []byte(refreshToken))
	if err != nil {
		log.Error("error verifying refresh token", zap.Error(err))
		return authCore.TokensModel{}, authCore.Unauthorized
	}

	log.Info("generating new tokens")
	newTokens, err := u.tokenManager.Generate(user)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	log.Info("gen hash token")
	tokenHash, err := bcrypt.GenerateFromPassword([]byte(newTokens.RefreshToken), bcrypt.DefaultCost)

	newSessionInfo := authCore.SessionModel{
		UserID:           user.ID,
		DeviceID:         deviceID,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        newTokens.RefreshExpiresAt,
	}

	log.Info("session info update")
	err = u.sessionStorage.Update(ctx, newSessionInfo)
	if err != nil {
		log.Error("error updating session", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	return newTokens, nil
}

func authorization(user []byte, currentPassword string) error {
	err := bcrypt.CompareHashAndPassword(user, []byte(currentPassword))
	if err != nil {
		return authCore.Unauthorized
	}
	return nil
}

func GetDefaultBalance() map[string]userModels.Balance {
	wallet := make(map[string]userModels.Balance)

	balance := userModels.Balance{
		Asset:     "USD",
		Available: userModels.Money{},
		Locked:    userModels.Money{},
	}

	wallet[balance.Asset] = balance
	return wallet
}
