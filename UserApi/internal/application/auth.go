package application

import (
	"ITK_Code/m/v2/internal/adapters/hash"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"errors"
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
		Deleted:      false,
	}

	uid, err := a.userSaver.SaveUser(ctx, newUser)
	if errors.Is(err, user.ErrEmailIsExist) {
		log.Info("email is exist", zap.String("email", email))
		return "", time.Time{}, user.ErrEmailIsExist
	}
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

	log.Info("create token")
	tokens, accessToken, refreshToken, err := a.tokenManager.Generate(gotUser, deviceId)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return auth.TokensModel{}, err
	}
	_ = refreshToken

	tokenHash := hash.GenerateHashSHA256(tokens.RefreshToken)

	log.Info("session info save")
	session := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         deviceId,
		RefreshTokenHash: tokenHash,
		TTL:              tokens.RefreshTTL,
		ExpiresAt:        tokens.RefreshExpiresAt,
		CreatedAt:        tokens.RefreshCreatedAt,
	}

	err = a.sessionStorage.Create(ctx, accessToken.Jti, session)
	if err != nil {
		log.Error("error creating session", zap.Error(err))
		return auth.TokensModel{}, err
	}

	return tokens, nil
}

func (a *Auth) Logout(ctx context.Context,
	refreshToken string,
) (
	success bool,
	loggedOutAt time.Time,
	err error,
) {
	log := a.log.Named("Logout")

	jtiFromContext, err := auth.GetJTIFromContext(ctx)
	if err != nil {
		log.Error("error getting jti from context", zap.Error(err))
		return false, time.Time{}, auth.ErrInvalidContext
	}

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jtiFromContext)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, auth.ErrSessionNotFound
	}

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token",
			zap.Error(err),
			zap.String("refreshToken", hash.GenerateHashSHA256(refreshToken)),
			zap.String("storedHash", sessionInfo.RefreshTokenHash),
		)
		return false, time.Time{}, auth.ErrNoAccess
	}

	err = a.sessionStorage.DeleteByJTI(ctx, jtiFromContext, sessionInfo.DeviceID)
	if err != nil {
		log.Error("error deleting session", zap.Error(err))
		return false, time.Time{}, auth.ErrSessionNotFound
	}

	return true, time.Now(), nil
}

func (a *Auth) LogoutAllDevices(ctx context.Context,
	refreshToken string,
) (
	bool,
	time.Time,
	error,
) {
	log := a.log.Named("Logout all devices")

	jti, err := auth.GetJTIFromContext(ctx)
	if err != nil {
		log.Error("error getting user id by context", zap.Error(err))
		return false, time.Time{}, auth.ErrInvalidContext
	}

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jti)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, auth.Unauthorized
	}

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token",
			zap.Error(err),
			zap.String("refreshToken", hash.GenerateHashSHA256(refreshToken)),
			zap.String("storedHash", sessionInfo.RefreshTokenHash),
		)
		return false, time.Time{}, auth.ErrNoAccess
	}

	err = a.sessionStorage.DeleteByUser(ctx, sessionInfo.UserID)
	if err != nil {
		log.Error("error deleting sessions", zap.Error(err))
		return false, time.Time{}, err
	}

	return true, time.Now(), nil
}

func (a *Auth) RefreshToken(ctx context.Context,
	refreshToken string,
) (
	auth.TokensModel,
	error,
) {
	log := a.log.Named("Refresh tokens")

	log.Info("parsing token")
	claims, err := a.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		log.Error("error parsing refresh token", zap.Error(err))
		return auth.TokensModel{}, auth.ErrInvalidToken
	}

	if claims.ExpiresAt.Before(time.Now()) {
		log.Error("token expired")
		return auth.TokensModel{}, auth.ErrRefreshExpired
	}

	storedJTI := claims.AccessTokenJti

	//синхронизация
	ok, err := a.syncPrimitiveForRedis.AcquireRefreshLock(ctx, storedJTI)
	if err != nil {
		log.Error("error acquiring refresh lock", zap.Error(err))
		return auth.TokensModel{}, auth.ErrSyncRedis
	}

	if !ok {
		log.Error("generate tokens processing", zap.Error(err))
		return auth.TokensModel{}, auth.ErrGenerateTokenProcessing
	}

	defer func() {

		releaseCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Second,
		)

		defer cancel()

		err = a.syncPrimitiveForRedis.ReleaseRefreshLock(
			releaseCtx,
			storedJTI,
		)
		if err != nil {
			log.Error("error releasing refresh lock", zap.Error(err))
			return
		}

	}()

	//поиск сессии
	log.Info("searching session with jti", zap.String("jti", storedJTI))
	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, storedJTI)
	if err != nil {
		log.Error("error getting session", zap.Error(err))
		return auth.TokensModel{}, auth.ErrSessionNotFound
	}

	userID := sessionInfo.UserID

	log.Info("searching user with db", zap.String("userId", userID))
	gotUser, err := a.userProvider.Get(ctx, userID)
	if err != nil {
		log.Error("error getting user by id", zap.Error(err))
		return auth.TokensModel{}, user.ErrUserNotFound
	}

	log.Info("comparing refresh token")
	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token",
			zap.Error(err),
			zap.String("refreshToken", hash.GenerateHashSHA256(refreshToken)),
			zap.String("storedHash", sessionInfo.RefreshTokenHash),
		)
		return auth.TokensModel{}, auth.ErrNoAccess
	}

	log.Info("generating new tokens")
	newTokens, accessToken, _, err := a.tokenManager.Generate(gotUser, sessionInfo.DeviceID)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return auth.TokensModel{}, auth.ErrGenerateToken
	}

	tokenHash := hash.GenerateHashSHA256(newTokens.RefreshToken)

	newSessionInfo := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         sessionInfo.DeviceID,
		RefreshTokenHash: tokenHash,
		TTL:              newTokens.RefreshTTL,
		ExpiresAt:        newTokens.RefreshExpiresAt,
		CreatedAt:        newTokens.RefreshCreatedAt,
	}

	log.Info("session info update")
	err = a.sessionStorage.Update(ctx, storedJTI, accessToken.Jti, newSessionInfo)
	if err != nil {
		log.Error("error updating session", zap.Error(err))
		return auth.TokensModel{}, auth.ErrGenerateToken
	}

	return newTokens, nil
}
