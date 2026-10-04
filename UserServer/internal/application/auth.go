package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/corerrors"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"errors"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.uber.org/zap"
)

type Auth struct {
	tokenManager          ports.TokenManager
	sessionStorage        ports.SessionRepository
	syncPrimitiveForRedis ports.RefreshLock
	rateLimiting          ports.RateLimiter
	passwordHasher        ports.PasswordHasher
	tokenHasher           ports.TokenHasher

	userRepository ports.UserRepository
}

func NewAuthService(
	tokenManager ports.TokenManager,
	sessionStorage ports.SessionRepository,
	syncPrimitiveForRedis ports.RefreshLock,
	rateLimiting ports.RateLimiter,
	userRepository ports.UserRepository,
	passwordHasher ports.PasswordHasher,
	tokenHasher ports.TokenHasher,
) *Auth {
	return &Auth{
		tokenManager:          tokenManager,
		sessionStorage:        sessionStorage,
		syncPrimitiveForRedis: syncPrimitiveForRedis,
		rateLimiting:          rateLimiting,
		userRepository:        userRepository,
		passwordHasher:        passwordHasher,
		tokenHasher:           tokenHasher,
	}
}

func (a *Auth) Registration(ctx context.Context,
	email string,
	password string,
	name string,
	id string,
	deviceID string,
) (
	string,
	time.Time,
	error,
) {
	now := time.Now()

	log := logging.FromContext(ctx).Named("Registration")

	allowed, err := a.rateLimiting.Allow(ctx, id, deviceID)
	if err != nil {
		logOperationError(log, "failed to check registration rate limit", err)
		return "", time.Time{}, corerrors.ErrTooManyRequests
	}
	if !allowed {
		log.Warn("registration rate limit exceeded")
		return "", time.Time{}, corerrors.ErrTooManyRequests
	}
	log.Debug("registration rate limit check passed")

	passHash, err := a.passwordHasher.GeneratePasswordHash(password)
	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, corerrors.ErrPassGenHash
	}
	log.Debug("password hash generated")

	newUser := user.User{
		Email:        email,
		Name:         name,
		PasswordHash: passHash,
		Role:         user.AdminRole,
		CreateTime:   now,
		UpdateTime:   now,
	}

	uid, err := a.userRepository.SaveUser(ctx, newUser)
	if err != nil {
		logOperationError(log, "failed to register user", err)
		return "", time.Time{}, err
	}
	// SaveUser can also restore a previously deleted account.
	log.Info("user registered", zap.String("user_id", uid))

	return uid, now, nil
}

func (a *Auth) Login(ctx context.Context,
	email string,
	password string,
	ip string,
	deviceID string,
) (
	dto.TokensModel,
	error,
) {
	log := logging.FromContext(ctx).Named("Login")
	if deviceID == "" {
		log.Warn("login rejected: device id is missing")
		return dto.TokensModel{}, corerrors.ErrDeviceIDEmpty
	}

	allowed, err := a.rateLimiting.Allow(ctx, ip, deviceID)
	if err != nil {
		logOperationError(log, "failed to check login rate limit", err)
		return dto.TokensModel{}, corerrors.ErrTooManyRequests
	}
	if !allowed {
		log.Warn("login rate limit exceeded")
		return dto.TokensModel{}, corerrors.ErrTooManyRequests
	}
	log.Debug("login rate limit check passed")

	gotUser, err := a.userRepository.GetByEmail(ctx, email)
	if errors.Is(err, user.ErrUserNotFound) {
		// Keep the missing-user path close to the existing-user path by
		// spending the same bcrypt work before returning a generic error.
		_ = a.passwordHasher.VerifyPasswordHash(password, []byte("$2a$14$fidR2tQBZMd5vck77HC6TeeEcC4oXWjR4jZqxP76Jpl1biQEaQmpa"))
		log.Warn("login rejected: incorrect credentials")
		return dto.TokensModel{}, auth.ErrIncorrectCredentials
	}
	if err != nil {
		logOperationError(log, "failed to get user for login", err)
		return dto.TokensModel{}, err
	}
	log = log.With(zap.String("user_id", gotUser.ID))
	log.Debug("user retrieved for login")

	err = a.passwordHasher.VerifyPasswordHash(password, gotUser.PasswordHash)
	if err != nil {
		logPasswordVerificationError(log, err)
		return dto.TokensModel{}, auth.ErrIncorrectCredentials
	}
	log.Debug("password verified")

	tokens, accessToken, _, err := a.tokenManager.Generate(gotUser, deviceID)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return dto.TokensModel{}, err
	}
	log.Debug("tokens generated")

	tokenHash := a.tokenHasher.GenerateHashSHA256(tokens.RefreshToken)
	session := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         deviceID,
		RefreshTokenHash: tokenHash,
		TTL:              tokens.RefreshTTL,
		ExpiresAt:        tokens.RefreshExpiresAt,
		CreatedAt:        tokens.RefreshIssuedAt,
	}

	err = a.sessionStorage.Create(ctx, accessToken.JTI, session)
	if err != nil {
		logOperationError(log, "failed to create session", err)
		return dto.TokensModel{}, err
	}
	log.Info("login completed; session created")

	return tokens, nil
}

func (a *Auth) Logout(ctx context.Context,
	jti string,
	refreshToken string,
) (
	err error,
) {
	log := logging.FromContext(ctx).Named("Logout")

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jti)
	if err != nil {
		logOperationError(log, "failed to get session for logout", err)
		return err
	}
	log = log.With(zap.String("user_id", sessionInfo.UserID))
	log.Debug("session retrieved")

	if err = a.tokenHasher.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Warn("refresh token does not match session")
		return auth.ErrNoAccess
	}
	log.Debug("refresh token verified")

	err = a.sessionStorage.DeleteByJTI(ctx, jti, sessionInfo.UserID)
	if err != nil {
		logOperationError(log, "failed to revoke session", err)
		return err
	}
	log.Info("logout completed; session revoked")

	return nil
}

func (a *Auth) LogoutAllDevices(ctx context.Context,
	jti string,
) error {
	log := logging.FromContext(ctx).Named("LogoutAllDevices")

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jti)
	if err != nil {
		logOperationError(log, "failed to get session for logout", err)
		return err
	}
	log = log.With(zap.String("user_id", sessionInfo.UserID))
	log.Debug("session retrieved")

	err = a.sessionStorage.DeleteByUser(ctx, sessionInfo.UserID)
	if err != nil {
		logOperationError(log, "failed to revoke user sessions", err)
		return err
	}
	log.Info("logout completed; all user sessions revoked")

	return nil
}

func (a *Auth) RefreshToken(ctx context.Context,
	refreshToken string,
) (
	dto.TokensModel,
	error,
) {
	log := logging.FromContext(ctx).Named("RefreshToken")

	log.Debug("parsing refresh token")
	claims, err := a.tokenManager.ParseRefreshToken(refreshToken)
	if err != nil {
		log.Warn("refresh token rejected")
		return dto.TokensModel{}, corerrors.ErrInvalidToken
	}
	log.Debug("refresh token parsed")

	storedJTI := claims.AccessTokenJTI

	lockOwner, err := a.syncPrimitiveForRedis.AcquireRefreshLock(ctx, storedJTI)
	if err != nil {
		logOperationError(log, "failed to acquire refresh lock", err)
		return dto.TokensModel{}, corerrors.ErrSyncRedis
	}
	if lockOwner == "" {
		log.Debug("token refresh already in progress")
		return dto.TokensModel{}, corerrors.ErrGenerateTokenProcessing
	}
	log.Debug("refresh lock acquired")

	defer func() {
		releaseCtx, cancel := context.WithTimeout(
			context.WithoutCancel(ctx),
			time.Second,
		)
		defer cancel()

		releaseErr := a.syncPrimitiveForRedis.ReleaseRefreshLock(
			releaseCtx,
			storedJTI,
			lockOwner,
		)
		if releaseErr != nil {
			log.Warn("failed to release refresh lock; waiting for lock expiry", zap.Error(releaseErr))
			return
		}
		log.Debug("refresh lock release completed")
	}()

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, storedJTI)
	if err != nil {
		logOperationError(log, "failed to get session for token refresh", err)
		return dto.TokensModel{}, err
	}
	log = log.With(zap.String("user_id", sessionInfo.UserID))
	log.Debug("session retrieved")

	userID := sessionInfo.UserID
	gotUser, err := a.userRepository.Get(ctx, userID)
	if err != nil {
		logOperationError(log, "failed to get user for token refresh", err)
		return dto.TokensModel{}, err
	}
	log.Debug("session user retrieved")

	if err = a.tokenHasher.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash); err != nil {
		log.Warn("refresh token does not match session")
		return dto.TokensModel{}, auth.ErrNoAccess
	}
	log.Debug("refresh token verified")

	newTokens, accessToken, _, err := a.tokenManager.Generate(gotUser, sessionInfo.DeviceID)
	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return dto.TokensModel{}, corerrors.ErrGenerateToken
	}
	log.Debug("new tokens generated")

	tokenHash := a.tokenHasher.GenerateHashSHA256(newTokens.RefreshToken)
	newSessionInfo := auth.SessionModel{
		UserID:           gotUser.ID,
		DeviceID:         sessionInfo.DeviceID,
		RefreshTokenHash: tokenHash,
		TTL:              newTokens.RefreshTTL,
		ExpiresAt:        newTokens.RefreshExpiresAt,
		CreatedAt:        newTokens.RefreshIssuedAt,
	}

	err = a.sessionStorage.Update(ctx, storedJTI, accessToken.JTI, newSessionInfo)
	if err != nil {
		logOperationError(log, "failed to rotate session", err)
		return dto.TokensModel{}, corerrors.ErrGenerateToken
	}
	log.Info("session rotated; tokens refreshed")

	return newTokens, nil
}
