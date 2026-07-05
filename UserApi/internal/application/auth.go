package application

import (
	"ITK_Code/m/v2/internal/adapters/hash"
	authCore "ITK_Code/m/v2/internal/core/auth"
	userCore "ITK_Code/m/v2/internal/core/user"
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

	if a.userSaver.IsExistsUserByEmail(ctx, email) {
		log.Info("user with email already exists")
		return "", time.Time{}, userCore.ErrUserExists
	}

	passHash, err := hash.GeneratePasswordHash(password)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, err
	}

	newUser := userCore.User{
		Email:        email,
		Name:         name,
		PasswordHash: passHash,
		Role:         userCore.UserRole,
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
	authCore.TokensModel,
	error,
) {
	log := a.log.Named("login")
	log.Info("searching user with email", zap.String("email", email))

	user, err := a.userProvider.GetByEmail(ctx, email)
	if err != nil {
		log.Error("error getting user by email", zap.String("email", email), zap.Error(userCore.ErrUserNotFound))
		return authCore.TokensModel{}, userCore.ErrUserNotFound
	}

	err = hash.VerifyPasswordHash(password, user.PasswordHash)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(err))
		return authCore.TokensModel{}, authCore.Unauthorized
	}

	log.Info("search old session")
	_ = a.sessionStorage.DeleteByUserAndDevice(ctx, user.ID, deviceId)

	log.Info("create token")
	tokens, err := a.tokenManager.Generate(user)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	tokenHash := hash.GenerateHashSHA256(tokens.RefreshToken)

	log.Info("session info save")
	session := authCore.SessionModel{
		UserID:           user.ID,
		DeviceID:         deviceId,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        tokens.RefreshExpiresAt,
	}

	err = a.sessionStorage.Create(ctx, session)
	if err != nil {
		log.Error("error creating session", zap.Error(err))
		return authCore.TokensModel{}, err
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

	userID := ctx.Value(authCore.UserIDContextKey).(string)
	if userID == "" {
		return false, time.Time{}, authCore.ErrInvalidToken
	}

	sessionInfo, err := a.sessionStorage.GetByUserAndDevice(ctx, userID, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, authCore.Unauthorized
	}

	ok := hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash)
	if !ok {
		log.Error("error comparing refresh token", zap.Error(err))
		return false, time.Time{}, authCore.ErrNoAccess
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
	success bool,
	loggedOutAt time.Time,
	err error,
) {
	log := a.log.Named("Logout all devices")

	userID := ctx.Value(authCore.UserIDContextKey).(string)
	if userID == "" {
		return false, time.Time{}, authCore.ErrInvalidToken
	}

	sessions, err := a.sessionStorage.GetByUserAndDevice(ctx, userID, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, authCore.Unauthorized
	}

	ok := hash.CompareHashSHA256(refreshToken, sessions.RefreshTokenHash)
	if !ok {
		log.Error("error comparing refresh token", zap.Error(err))
		return false, time.Time{}, authCore.ErrNoAccess
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
	tokensPairs authCore.TokensModel,
	err error,
) {
	log := a.log.Named("Refresh tokens")

	userID := ctx.Value(authCore.UserIDContextKey).(string)
	if userID == "" {
		return authCore.TokensModel{}, authCore.ErrInvalidToken
	}

	log.Info("searching user with db", zap.String("userId", userID))
	user, err := a.userProvider.Get(ctx, userID)
	if err != nil {
		log.Error("error getting user by id", zap.Error(err))
		return authCore.TokensModel{}, userCore.ErrUserNotFound
	}

	log.Info("searching session with userId deviceID", zap.String("userId", userID), zap.String("deviceID", deviceID))
	sessionInfo, err := a.sessionStorage.GetByUserAndDevice(ctx, userID, deviceID)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return authCore.TokensModel{}, authCore.Unauthorized
	}

	log.Info("validating token expiry", zap.String("ExpireAt", sessionInfo.ExpiresAt.String()))
	if sessionInfo.ExpiresAt.Before(time.Now()) {
		return authCore.TokensModel{}, authCore.ErrRefreshExpired
	}

	log.Info("comparing refresh token")
	ok := hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash)
	if !ok {
		log.Error("error verifying refresh token", zap.Error(err))
		return authCore.TokensModel{}, authCore.Unauthorized
	}

	log.Info("generating new tokens")
	newTokens, err := a.tokenManager.Generate(user)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	tokenHash := hash.GenerateHashSHA256(refreshToken)

	newSessionInfo := authCore.SessionModel{
		UserID:           user.ID,
		DeviceID:         deviceID,
		RefreshTokenHash: tokenHash,
		ExpiresAt:        newTokens.RefreshExpiresAt,
	}

	log.Info("session info update")
	err = a.sessionStorage.Update(ctx, newSessionInfo)
	if err != nil {
		log.Error("error updating session", zap.Error(err))
		return authCore.TokensModel{}, err
	}

	return newTokens, nil
}
