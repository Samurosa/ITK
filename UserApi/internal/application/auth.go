package application

import (
	"ITK_Code/m/v2/internal/adapters/hash"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"time"

	"go.uber.org/zap"
)

func (a *Auth) Registration(ctx context.Context,
	email string,
	password string,
	name string,
) (
	string,
	time.Time,
	error,
) {
	now := time.Now()

	log := a.log.Named("RegisterNewUser")
	passHash, err := hash.GeneratePasswordHash(password)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, user.ErrPassGenHash
	}

	newUser := user.User{
		Email:        email,
		Name:         name,
		PasswordHash: passHash,
		Role:         user.UserRole,
		CreateTime:   now,
		UpdateTime:   now,
	}

	uid, err := a.userSaver.SaveUser(ctx, newUser)
	if err != nil {
		log.Error("error saving user", zap.Error(err))
		return "", time.Time{}, err
	}
	log.Info("user created", zap.String("uid", uid), zap.String("createdAt", now.String()))

	return uid, now, nil
}

func (a *Auth) Login(ctx context.Context,
	email string,
	password string,
	deviceId string,
) (
	auth.TokensModel,
	error,
) {
	log := a.log.Named("login")
	log.Info("searching user with email", zap.String("email", email))

	gotUser, err := a.userProvider.GetByEmail(ctx, email)
	if err != nil {
		log.Error("error getting user by email", zap.String("email", email), zap.Error(err))
		return auth.TokensModel{}, auth.ErrInvalidLoginCredentials
	}

	err = hash.VerifyPasswordHash(password, gotUser.PasswordHash)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(err))
		return auth.TokensModel{}, auth.ErrInvalidLoginCredentials
	}

	log.Info("search old session")
	_ = a.sessionStorage.DeleteByUserAndDevice(ctx, gotUser.ID, deviceId)

	log.Info("create token")
	tokens, err := a.tokenManager.Generate(gotUser, deviceId)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return auth.TokensModel{}, err
	}

	tokenHash := hash.GenerateHashSHA256(tokens.RefreshToken)

	log.Info("session info save")
	session := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         deviceId,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        tokens.RefreshExpiresAt,
	}

	err = a.sessionStorage.Create(ctx, session)
	if err != nil {
		log.Error("error creating session", zap.Error(err))
		return auth.TokensModel{}, err
	}

	return tokens, nil
}

func (a *Auth) Logout(ctx context.Context,
	refreshToken string,
	deviceID string,
) (
	success bool,
	loggedOutAt time.Time,
	err error,
) {
	log := a.log.Named("Logout")

	userID, err := auth.GetUserIDByContext(ctx)
	if err != nil {
		log.Error("error getting user id by context", zap.Error(err))
		return false, time.Time{}, auth.ErrInvalidContext
	}

	sessionInfo, err := a.sessionStorage.GetByUserAndDevice(ctx, userID, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, auth.Unauthorized
	}

	ok := hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash)
	if !ok {
		log.Error("error comparing refresh token", zap.Error(err))
		return false, time.Time{}, auth.ErrNoAccess
	}

	err = a.sessionStorage.DeleteByUserAndDevice(ctx, userID, deviceID)
	if err != nil {
		log.Error("error deleting session", zap.Error(err))
		return false, time.Time{}, err
	}

	return true, time.Now(), nil
}

func (a *Auth) LogoutAllDevices(ctx context.Context,
	refreshToken string,
	deviceID string,
) (
	bool,
	time.Time,
	error,
) {
	log := a.log.Named("Logout all devices")

	userID, err := auth.GetUserIDByContext(ctx)
	if err != nil {
		log.Error("error getting user id by context", zap.Error(err))
		return false, time.Time{}, auth.ErrInvalidContext
	}

	sessionInfo, err := a.sessionStorage.GetByUserAndDevice(ctx, userID, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, auth.Unauthorized
	}

	ok := hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash)
	if !ok {
		log.Error("error comparing refresh token", zap.Error(err))
		return false, time.Time{}, auth.ErrNoAccess
	}

	err = a.sessionStorage.DeleteByUser(ctx, userID)
	if err != nil {
		log.Error("error deleting sessions", zap.Error(err))
		return false, time.Time{}, err
	}

	return true, time.Now(), nil
}

func (a *Auth) RefreshToken(ctx context.Context,
	refreshToken string,
	deviceID string,
) (
	auth.TokensModel,
	error,
) {
	log := a.log.Named("Refresh tokens")

	refreshTokenHash := hash.GenerateHashSHA256(refreshToken)

	session, err := a.sessionStorage.GetByRefreshToken(ctx, refreshTokenHash)
	if err != nil {
		log.Error("error getting session", zap.Error(err))
		return auth.TokensModel{}, auth.ErrSessionNotFound
	}

	userID := session.UserID

	log.Info("searching user with db", zap.String("userId", userID))
	gotUser, err := a.userProvider.Get(ctx, userID)
	if err != nil {
		log.Error("error getting user by id", zap.Error(err))
		return auth.TokensModel{}, user.ErrUserNotFound
	}

	log.Info("searching session with userId deviceID", zap.String("userId", userID), zap.String("deviceID", deviceID))
	sessionInfo, err := a.sessionStorage.GetByUserAndDevice(ctx, userID, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return auth.TokensModel{}, auth.Unauthorized
	}

	log.Info("validating token expiry", zap.String("ExpireAt", sessionInfo.ExpiresAt.String()))
	if sessionInfo.ExpiresAt.Before(time.Now()) {
		return auth.TokensModel{}, auth.ErrRefreshExpired
	}

	log.Info("comparing refresh token")
	ok := hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash)
	if !ok {
		log.Error("error verifying refresh token", zap.Error(err))
		return auth.TokensModel{}, auth.Unauthorized
	}

	log.Info("generating new tokens")
	newTokens, err := a.tokenManager.Generate(gotUser, deviceID)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return auth.TokensModel{}, err
	}

	tokenHash := hash.GenerateHashSHA256(refreshToken)

	newSessionInfo := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         deviceID,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        newTokens.RefreshExpiresAt,
	}

	log.Info("session info update")
	err = a.sessionStorage.Update(ctx, newSessionInfo)
	if err != nil {
		log.Error("error updating session", zap.Error(err))
		return auth.TokensModel{}, err
	}

	return newTokens, nil
}
