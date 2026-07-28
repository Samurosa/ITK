package application

import (
	context2 "ITK_Code/m/v2/internal/adapters/outbound/context"
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/hash"
	"ITK_Code/m/v2/internal/core/dto"
	errors2 "ITK_Code/m/v2/internal/core/errors"
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

	ip, err := context2.GetClientIPFromContext(ctx)
	if err != nil {
		log.Error("Failed to get user ip from context", zap.Error(err))
		return "", time.Time{}, errors2.ErrInvalidContext
	}
	log.Info("got user IP from context")

	allowed, err := a.rateLimiting.Allow(ctx, ip)
	if err != nil {
		log.Error("Failed to check rate limiting", zap.Error(err))
		return "", time.Time{}, errors2.ErrTooManyRequests
	}
	if !allowed {
		log.Error("Rate limiting is not allowed", zap.String("ip", ip))
		return "", time.Time{}, errors2.ErrTooManyRequests
	}
	log.Info("validate rate limiting")

	passHash, err := hash.GeneratePasswordHash(password)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, errors2.ErrPassGenHash
	}
	log.Info("generate password hash", zap.String("id", ip))

	newUser := dto.User{
		Email:        email,
		Name:         name,
		PasswordHash: passHash,
		Role:         dto.UserRole,
		CreateTime:   now,
		UpdateTime:   now,
	}

	uid, err := a.userSaver.SaveUser(ctx, newUser)
	if errors.Is(err, errors2.ErrEmailIsExist) {
		log.Info("email is exist", zap.String("email", email))
		return "", time.Time{}, errors2.ErrEmailIsExist
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
	dto.TokensModel,
	error,
) {
	log := a.log.Named("login")

	gotUser, err := a.userProvider.GetByEmail(ctx, email)
	if err != nil {
		log.Error("error getting user by email", zap.String("email", email), zap.Error(err))
		return dto.TokensModel{}, errors2.ErrInvalidLoginCredentials
	}
	log.Info("got user by email", zap.String("email", email), zap.String("id", gotUser.ID))

	err = hash.VerifyPasswordHash(password, gotUser.PasswordHash)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrInvalidLoginCredentials
	}
	log.Info("verify password passed", zap.String("id", gotUser.ID))

	tokens, accessToken, _, err := a.tokenManager.Generate(gotUser, deviceId)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return dto.TokensModel{}, err
	}
	log.Info("generate tokens for user", zap.String("id", gotUser.ID))

	tokenHash := hash.GenerateHashSHA256(tokens.RefreshToken)
	session := dto.SessionModel{
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
		return dto.TokensModel{}, err
	}
	log.Info("session created", zap.String("id", gotUser.ID))

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

	jtiFromContext, err := context2.GetJTIFromContext(ctx)
	if err != nil {
		log.Error("error getting jti from context", zap.Error(err))
		return false, time.Time{}, errors2.ErrInvalidContext
	}
	log.Info("got token jti from context")

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jtiFromContext)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, errors2.ErrSessionNotFound
	}
	log.Info("got session", zap.String("id", sessionInfo.UserID))

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token",
			zap.Error(err),
			zap.String("refreshToken", hash.GenerateHashSHA256(refreshToken)),
			zap.String("storedHash", sessionInfo.RefreshTokenHash),
		)
		return false, time.Time{}, errors2.ErrNoAccess
	}
	log.Info("verify password passed", zap.String("id", sessionInfo.UserID))

	err = a.sessionStorage.DeleteByJTI(ctx, jtiFromContext, sessionInfo.DeviceID)
	if err != nil {
		log.Error("error deleting session", zap.Error(err))
		return false, time.Time{}, errors2.ErrSessionNotFound
	}
	log.Info("session deleted", zap.String("id", sessionInfo.UserID))

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

	jti, err := context2.GetJTIFromContext(ctx)
	if err != nil {
		log.Error("error getting user id by context", zap.Error(err))
		return false, time.Time{}, errors2.ErrInvalidContext
	}
	log.Info("got token jti from context")

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jti)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return false, time.Time{}, errors2.Unauthorized
	}
	log.Info("got session", zap.String("id", sessionInfo.UserID))

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token",
			zap.Error(err),
			zap.String("refreshToken", hash.GenerateHashSHA256(refreshToken)),
			zap.String("storedHash", sessionInfo.RefreshTokenHash),
		)
		return false, time.Time{}, errors2.ErrNoAccess
	}
	log.Info("verify password passed", zap.String("id", sessionInfo.UserID))

	err = a.sessionStorage.DeleteByUser(ctx, sessionInfo.UserID)
	if err != nil {
		log.Error("error deleting sessions", zap.Error(err))
		return false, time.Time{}, err
	}
	log.Info("user sessions deleted", zap.String("id", sessionInfo.UserID))

	return true, time.Now(), nil
}

func (a *Auth) RefreshToken(ctx context.Context,
	refreshToken string,
) (
	dto.TokensModel,
	error,
) {
	log := a.log.Named("Refresh tokens")

	log.Info("parsing token")
	claims, err := a.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		log.Error("error parsing refresh token", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrInvalidToken
	}
	log.Info("parsed refresh token is successful")

	storedJTI := claims.AccessTokenJTI
	//синхронизация
	ok, err := a.syncPrimitiveForRedis.AcquireRefreshLock(ctx, storedJTI)
	if err != nil {
		log.Error("error acquiring refresh lock", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrSyncRedis
	}
	if !ok {
		log.Error("generate tokens processing", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrGenerateTokenProcessing
	}
	log.Info("acquiring refresh lock success")

	defer func() {
		releaseCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Second,
		)
		defer cancel()
		defer log.Info("released refresh lock success")

		err = a.syncPrimitiveForRedis.ReleaseRefreshLock(
			releaseCtx,
			storedJTI,
		)
		if err != nil {
			log.Error("error releasing refresh lock", zap.Error(err))
			return
		}

	}()

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, storedJTI)
	if err != nil {
		log.Error("error getting session", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrSessionNotFound
	}
	log.Info("got session", zap.String("id", sessionInfo.UserID))

	userID := sessionInfo.UserID
	gotUser, err := a.userProvider.Get(ctx, userID)
	if err != nil {
		log.Error("error getting user by id", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrUserNotFound
	}
	log.Info("got user from session info", zap.String("id", userID))

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrNoAccess
	}
	log.Info("verify password passed", zap.String("id", sessionInfo.UserID))

	newTokens, accessToken, _, err := a.tokenManager.Generate(gotUser, sessionInfo.DeviceID)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrGenerateToken
	}
	log.Info("generated new tokens", zap.String("id", sessionInfo.UserID))

	tokenHash := hash.GenerateHashSHA256(newTokens.RefreshToken)
	newSessionInfo := dto.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         sessionInfo.DeviceID,
		RefreshTokenHash: tokenHash,
		TTL:              newTokens.RefreshTTL,
		ExpiresAt:        newTokens.RefreshExpiresAt,
		CreatedAt:        newTokens.RefreshCreatedAt,
	}

	err = a.sessionStorage.Update(ctx, storedJTI, accessToken.Jti, newSessionInfo)
	if err != nil {
		log.Error("error updating session", zap.Error(err))
		return dto.TokensModel{}, errors2.ErrGenerateToken
	}
	log.Info("session updated", zap.String("id", sessionInfo.UserID))

	return newTokens, nil
}
