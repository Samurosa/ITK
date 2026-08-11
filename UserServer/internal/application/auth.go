package application

import (
	"ITK_Code/m/v2/internal/adapters/outbound/crypto/hash"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errors"
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

	allowed, err := a.rateLimiting.Allow(ctx)
	if err != nil {
		log.Error("Failed to check rate limiting", zap.Error(err))
		return "", time.Time{}, errors.ErrTooManyRequests
	}
	if !allowed {
		log.Error("Rate limiting is not allowed")
		return "", time.Time{}, errors.ErrTooManyRequests
	}
	log.Debug("validate rate limiting")

	passHash, err := hash.GeneratePasswordHash(password)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, errors.ErrPassGenHash
	}
	log.Debug("generate password hash")

	newUser := user.User{
		Email:        email,
		Name:         name,
		PasswordHash: passHash,
		Role:         user.UserRole,
		CreateTime:   now,
		UpdateTime:   now,
	}

	uid, err := a.userRepository.SaveUser(ctx, newUser)
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

	allowed, err := a.rateLimiting.Allow(ctx)
	if err != nil {
		log.Error("Failed to check rate limiting", zap.Error(err))
		return dto.TokensModel{}, errors.ErrTooManyRequests
	}
	if !allowed {
		log.Error("Rate limiting is not allowed")
		return dto.TokensModel{}, errors.ErrTooManyRequests
	}
	log.Debug("validate rate limiting")

	gotUser, err := a.userRepository.GetByEmail(ctx, email)
	if err != nil {
		log.Error("error getting user by email", zap.String("email", email), zap.Error(err))
		gotUser.PasswordHash = []byte("$2a$14$fidR2tQBZMd5vck77HC6TeeEcC4oXWjR4jZqxP76Jpl1biQEaQmpa")
	}
	log.Debug("got user by email", zap.String("email", email), zap.String("id", gotUser.ID))

	err = hash.VerifyPasswordHash(password, gotUser.PasswordHash)
	if err != nil {
		log.Error("error verifying user by password", zap.Error(err))
		return dto.TokensModel{}, auth.ErrIncorrectCredentials
	}
	log.Debug("verify password passed", zap.String("id", gotUser.ID))

	tokens, accessToken, _, err := a.tokenManager.Generate(gotUser, deviceId)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return dto.TokensModel{}, err
	}
	log.Debug("generate tokens for user", zap.String("id", gotUser.ID))

	tokenHash := hash.GenerateHashSHA256(tokens.RefreshToken)
	session := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         deviceId,
		RefreshTokenHash: tokenHash,
		TTL:              tokens.RefreshTTL,
		ExpiresAt:        tokens.RefreshExpiresAt,
		CreatedAt:        tokens.RefreshIssuedAt,
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
	jti string,
	refreshToken string,
) (
	err error,
) {
	log := a.log.Named("Logout")

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jti)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return auth.ErrSessionNotFound
	}
	log.Debug("got session", zap.String("id", sessionInfo.UserID))

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token",
			zap.Error(err),
			zap.String("refreshToken", hash.GenerateHashSHA256(refreshToken)),
			zap.String("storedHash", sessionInfo.RefreshTokenHash),
		)
		return auth.ErrNoAccess
	}
	log.Debug("verify password passed", zap.String("id", sessionInfo.UserID))

	err = a.sessionStorage.DeleteByJTI(ctx, jti, sessionInfo.DeviceID)
	if err != nil {
		log.Error("error deleting session", zap.Error(err))
		return auth.ErrSessionNotFound
	}
	log.Info("session deleted", zap.String("id", sessionInfo.UserID))

	return nil
}

func (a *Auth) LogoutAllDevices(ctx context.Context,
	jti string,
	refreshToken string,
) error {
	log := a.log.Named("Logout all devices")

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jti)
	if err != nil {
		log.Error("error getting session info", zap.Error(err))
		return auth.ErrUnauthorized
	}
	log.Debug("got session", zap.String("id", sessionInfo.UserID))

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token",
			zap.Error(err),
			zap.String("refreshToken", hash.GenerateHashSHA256(refreshToken)),
			zap.String("storedHash", sessionInfo.RefreshTokenHash),
		)
		return auth.ErrNoAccess
	}
	log.Debug("verify password passed", zap.String("id", sessionInfo.UserID))

	err = a.sessionStorage.DeleteByUser(ctx, sessionInfo.UserID)
	if err != nil {
		log.Error("error deleting sessions", zap.Error(err))
		return auth.ErrSessionNotFound
	}
	log.Info("user sessions deleted", zap.String("id", sessionInfo.UserID))

	return nil
}

func (a *Auth) RefreshToken(ctx context.Context,
	refreshToken string,
) (
	dto.TokensModel,
	error,
) {
	log := a.log.Named("Refresh tokens")

	log.Debug("parsing token")
	claims, err := a.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		log.Error("error parsing refresh token", zap.Error(err))
		return dto.TokensModel{}, errors.ErrInvalidToken
	}
	log.Debug("parsed refresh token is successful")

	storedJTI := claims.AccessTokenJTI
	//синхронизация
	ok, err := a.syncPrimitiveForRedis.AcquireRefreshLock(ctx, storedJTI)
	if err != nil {
		log.Error("error acquiring refresh lock", zap.Error(err))
		return dto.TokensModel{}, errors.ErrSyncRedis
	}
	if !ok {
		log.Error("generate tokens processing", zap.Error(err))
		return dto.TokensModel{}, errors.ErrGenerateTokenProcessing
	}
	log.Debug("acquiring refresh lock success")

	defer func() {
		releaseCtx, cancel := context.WithTimeout(
			context.Background(),
			time.Second,
		)
		defer cancel()
		defer log.Debug("released refresh lock success")

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
		return dto.TokensModel{}, auth.ErrSessionNotFound
	}
	log.Debug("got session", zap.String("id", sessionInfo.UserID))

	userID := sessionInfo.UserID
	gotUser, err := a.userRepository.Get(ctx, userID)
	if err != nil {
		log.Error("error getting user by id", zap.Error(err))
		return dto.TokensModel{}, user.ErrUserNotFound
	}
	log.Debug("got user from session info", zap.String("id", userID))

	if err = hash.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Error("error comparing refresh token", zap.Error(err))
		return dto.TokensModel{}, auth.ErrNoAccess
	}
	log.Debug("verify password passed", zap.String("id", sessionInfo.UserID))

	newTokens, accessToken, _, err := a.tokenManager.Generate(gotUser, sessionInfo.DeviceID)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return dto.TokensModel{}, errors.ErrGenerateToken
	}
	log.Debug("generated new tokens", zap.String("id", sessionInfo.UserID))

	tokenHash := hash.GenerateHashSHA256(newTokens.RefreshToken)
	newSessionInfo := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         sessionInfo.DeviceID,
		RefreshTokenHash: tokenHash,
		TTL:              newTokens.RefreshTTL,
		ExpiresAt:        newTokens.RefreshExpiresAt,
		CreatedAt:        newTokens.RefreshIssuedAt,
	}

	err = a.sessionStorage.Update(ctx, storedJTI, accessToken.Jti, newSessionInfo)
	if err != nil {
		log.Error("error updating session", zap.Error(err))
		return dto.TokensModel{}, errors.ErrGenerateToken
	}
	log.Info("session tokens refreshed", zap.String("id", sessionInfo.UserID))

	return newTokens, nil
}
