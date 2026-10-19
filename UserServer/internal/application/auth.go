package application

import (
	"ITK_Code/m/v2/internal/application/ports"
	"ITK_Code/m/v2/internal/core/auth"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/errorcore"
	"ITK_Code/m/v2/internal/core/user"
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

var authTracer = otel.Tracer("user-service/application/auth")

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
	userID string,
	createdAt time.Time,
	err error,
) {
	ctx, span := authTracer.Start(ctx, "Auth.Registration")
	defer func() {
		recordSpanError(span, err)
		span.End()
	}()

	log := logging.FromContext(ctx).Named("Registration")

	now := time.Now()

	allowed, err := a.rateLimiting.Allow(ctx, id, deviceID)
	if err != nil {
		logOperationError(log, "failed to check registration rate limit", err)
		return "", time.Time{}, fmt.Errorf("check registration rate limit: %w", err)
	}

	if !allowed {
		log.Warn("registration rate limit exceeded")
		return "", time.Time{}, errorcore.ErrTooManyRequests
	}

	log.Debug("registration rate limit check passed")

	_, hashSpan := authTracer.Start(ctx, "Auth.HashPassword")

	passHash, err := a.passwordHasher.GeneratePasswordHash(password)

	if err != nil {
		hashSpan.RecordError(err)
		hashSpan.SetStatus(codes.Error, "password hashing failed")
	}
	hashSpan.End()

	if err != nil {
		log.Error("error generating password hash", zap.Error(err))
		return "", time.Time{}, fmt.Errorf("%w: %w", errorcore.ErrPassGenHash, err)
	}

	log.Debug("password hash generated")

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
		logOperationError(log, "failed to register user", err)
		return "", time.Time{}, err
	}

	log.Info("user registered", zap.String("user_id", uid))

	return uid, now, nil
}

func (a *Auth) Login(ctx context.Context,
	email string,
	password string,
	ip string,
	deviceID string,
) (
	tokens dto.TokensModel,
	err error,
) {
	ctx, span := authTracer.Start(ctx, "Auth.Login")
	defer func() {
		recordSpanError(span, err)
		span.End()
	}()

	log := logging.FromContext(ctx).Named("Login")
	if deviceID == "" {
		log.Warn("login rejected: device id is missing")
		return dto.TokensModel{}, errorcore.ErrDeviceIDEmpty
	}

	allowed, err := a.rateLimiting.Allow(ctx, ip, deviceID)
	if err != nil {
		logOperationError(log, "failed to check login rate limit", err)
		return dto.TokensModel{}, fmt.Errorf("check login rate limit: %w", err)
	}

	if !allowed {
		log.Warn("login rate limit exceeded")
		return dto.TokensModel{}, errorcore.ErrTooManyRequests
	}

	log.Debug("login rate limit check passed")

	gotUser, err := a.userRepository.GetByEmail(ctx, email)

	if errors.Is(err, user.ErrUserNotFound) {
		_, verifySpan := authTracer.Start(ctx, "Auth.VerifyPassword")

		_ = a.passwordHasher.VerifyPasswordHash(
			password,
			[]byte("$2a$14$fidR2tQBZMd5vck77HC6TeeEcC4oXWjR4jZqxP76Jpl1biQEaQmpa"),
		)

		verifySpan.End()

		log.Warn("login rejected: incorrect credentials")
		return dto.TokensModel{}, auth.ErrIncorrectCredentials
	}

	if err != nil {
		logOperationError(log, "failed to get user for login", err)
		return dto.TokensModel{}, err
	}

	log = log.With(zap.String("user_id", gotUser.ID))
	log.Debug("user retrieved for login")

	_, verifySpan := authTracer.Start(ctx, "Auth.VerifyPassword")

	err = a.passwordHasher.VerifyPasswordHash(
		password,
		gotUser.PasswordHash,
	)

	if err != nil && !errors.Is(err, auth.ErrIncorrectPassword) {
		verifySpan.RecordError(err)
		verifySpan.SetStatus(codes.Error, "password verification failed")
	}

	verifySpan.End()

	if err != nil {
		logPasswordVerificationError(log, err)
		if errors.Is(err, auth.ErrIncorrectPassword) {
			return dto.TokensModel{}, auth.ErrIncorrectCredentials
		}
		return dto.TokensModel{}, fmt.Errorf("%w: %w", errorcore.ErrPassVerifyHash, err)
	}
	log.Debug("password verified")

	_, tokensGenSpan := authTracer.Start(ctx, "Auth.GenerateTokens")

	tokens, accessToken, _, err := a.tokenManager.Generate(gotUser, deviceID)
	if err != nil {
		tokensGenSpan.RecordError(err)
		tokensGenSpan.SetStatus(codes.Error, "failed to generate tokens")
	}
	tokensGenSpan.End()

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
	ctx, span := authTracer.Start(ctx, "Auth.Logout")
	defer func() {
		recordSpanError(span, err)
		span.End()
	}()

	log := logging.FromContext(ctx).Named("Logout")

	sessionInfo, err := a.sessionStorage.GetByJTI(ctx, jti)
	if err != nil {
		logOperationError(log, "failed to get session for logout", err)
		return err
	}

	log = log.With(zap.String("user_id", sessionInfo.UserID))
	log.Debug("session retrieved")

	_, tokenCompareSpan := authTracer.Start(ctx, "Auth.CompareToken")

	err = a.tokenHasher.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash)

	tokenCompareSpan.End()

	if err != nil {
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
) (err error) {
	ctx, span := authTracer.Start(ctx, "Auth.LogoutAllDevices")
	defer func() {
		recordSpanError(span, err)
		span.End()
	}()

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
	tokens dto.TokensModel,
	err error,
) {
	ctx, span := authTracer.Start(ctx, "Auth.RefreshToken")
	defer func() {
		recordSpanError(span, err)
		span.End()
	}()

	log := logging.FromContext(ctx).Named("RefreshToken")

	log.Debug("parsing refresh token")

	_, parseRefreshTokenSpan := authTracer.Start(ctx, "Auth.ParseRefreshToken")

	claims, err := a.tokenManager.ParseRefreshToken(refreshToken)

	parseRefreshTokenSpan.End()

	if err != nil {
		log.Warn("refresh token rejected")
		return dto.TokensModel{}, errorcore.ErrInvalidToken
	}

	log.Debug("refresh token parsed")

	storedJTI := claims.AccessTokenJTI

	lockOwner, err := a.syncPrimitiveForRedis.AcquireRefreshLock(ctx, storedJTI)
	if err != nil {
		logOperationError(log, "failed to acquire refresh lock", err)
		return dto.TokensModel{}, fmt.Errorf("%w: %w", errorcore.ErrSyncRedis, err)
	}

	if lockOwner == "" {
		log.Debug("token refresh already in progress")
		return dto.TokensModel{}, errorcore.ErrGenerateTokenProcessing
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

	_, tokenCompareSpan := authTracer.Start(ctx, "Auth.CompareToken")

	err = a.tokenHasher.CompareHashSHA256(refreshToken, sessionInfo.RefreshTokenHash)

	tokenCompareSpan.End()

	if err != nil {
		log.Warn("refresh token does not match session")
		return dto.TokensModel{}, auth.ErrNoAccess
	}

	log.Debug("refresh token verified")

	_, tokensGenSpan := authTracer.Start(ctx, "Auth.GenerateTokens")

	newTokens, accessToken, _, err := a.tokenManager.Generate(gotUser, sessionInfo.DeviceID)
	if err != nil {
		tokensGenSpan.RecordError(err)
		tokensGenSpan.SetStatus(codes.Error, "failed to generate tokens")
	}

	tokensGenSpan.End()

	if err != nil {
		log.Error("error generating tokens", zap.Error(err))
		return dto.TokensModel{}, fmt.Errorf("%w: %w", errorcore.ErrGenerateToken, err)
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
		return dto.TokensModel{}, fmt.Errorf("rotate session: %w", err)
	}

	log.Info("session rotated; tokens refreshed")

	return newTokens, nil
}

func recordSpanError(span trace.Span, err error) {
	if err == nil ||
		errors.Is(err, context.Canceled) ||
		errors.Is(err, auth.ErrIncorrectCredentials) ||
		errors.Is(err, auth.ErrNoAccess) ||
		errors.Is(err, errorcore.ErrTooManyRequests) {
		return
	}

	span.RecordError(err)
	span.SetStatus(codes.Error, "application operation failed")
}
