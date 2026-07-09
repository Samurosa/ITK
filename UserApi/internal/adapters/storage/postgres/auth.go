package postgres

import (
	authCore "ITK_Code/m/v2/internal/core/auth"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type SessionRepository struct {
	pool *pgxpool.Pool
}

func NewSessionStorage(pool *pgxpool.Pool) *SessionRepository {
	return &SessionRepository{
		pool: pool,
	}
}

func (s *SessionRepository) Create(
	ctx context.Context,
	session authCore.SessionModel,
) error {

	query := `
		INSERT INTO sessions
		(
			user_id,
			device_id,
			refresh_token_hash,
			expires_at
		)
		VALUES ($1,$2,$3,$4)
	`

	_, err := s.pool.Exec(
		ctx,
		query,
		session.UserID,
		session.DeviceID,
		session.RefreshTokenHash,
		session.ExpiresAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (s *SessionRepository) GetByUserAndDevice(ctx context.Context,
	userID string,
	deviceID string,
) (
	authCore.SessionModel,
	error,
) {
	query := `
		SELECT id,
		       user_id,
		       device_id,
		       refresh_token_hash,
		       expires_at,
		       created_at
FROM sessions
WHERE user_id=$1
AND device_id=$2;
	`
	var session authCore.SessionModel

	err := s.pool.QueryRow(
		ctx,
		query,
		userID,
		deviceID,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.RefreshTokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return authCore.SessionModel{}, authCore.ErrSessionNotFound
	}

	return session, nil
}

func (s *SessionRepository) GetByRefreshToken(ctx context.Context,
	refreshTokenHash []byte,
) (
	authCore.SessionModel,
	error,
) {
	query := `
		SELECT id,
		       user_id,
		       device_id,
		       refresh_token_hash,
		       expires_at,
		       created_at
FROM sessions
WHERE refresh_token_hash=$1;
	`
	var session authCore.SessionModel

	err := s.pool.QueryRow(
		ctx,
		query,
		refreshTokenHash,
	).Scan(
		&session.ID,
		&session.UserID,
		&session.DeviceID,
		&session.RefreshTokenHash,
		&session.ExpiresAt,
		&session.CreatedAt,
	)

	if err != nil {
		return authCore.SessionModel{}, err
	}

	return session, nil
}

func (s *SessionRepository) GetAllByUser(
	ctx context.Context,
	userID string,
) (
	[]authCore.SessionModel,
	error,
) {

	query := `
		SELECT
			id,
			user_id,
			device_id,
			refresh_token_hash,
			expires_at,
			created_at
		FROM sessions
		WHERE user_id=$1
	`

	rows, err := s.pool.Query(
		ctx,
		query,
		userID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	sessions := make([]authCore.SessionModel, 0)

	for rows.Next() {

		var session authCore.SessionModel

		err := rows.Scan(
			&session.ID,
			&session.UserID,
			&session.DeviceID,
			&session.RefreshTokenHash,
			&session.ExpiresAt,
			&session.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		sessions = append(
			sessions,
			session,
		)
	}

	if len(sessions) == 0 {
		return sessions, err
	}

	return sessions, nil
}

func (s *SessionRepository) Update(
	ctx context.Context,
	session authCore.SessionModel,
) error {

	query := `
		UPDATE sessions
		SET
			refresh_token_hash=$1,
			expires_at=$2
		WHERE user_id=$3
		AND device_id=$4
	`

	result, err := s.pool.Exec(
		ctx,
		query,
		session.RefreshTokenHash,
		session.ExpiresAt,
		session.UserID,
		session.DeviceID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return err
	}

	return nil
}

func (s *SessionRepository) DeleteByUserAndDevice(
	ctx context.Context,
	userID string,
	deviceID string,
) error {

	query := `
		DELETE FROM sessions
		WHERE user_id=$1
		AND device_id=$2
	`

	result, err := s.pool.Exec(
		ctx,
		query,
		userID,
		deviceID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return err
	}

	return nil
}

func (s *SessionRepository) DeleteByUser(
	ctx context.Context,
	userID string,
) error {

	query := `
		DELETE FROM sessions
		WHERE user_id=$1
	`

	_, err := s.pool.Exec(
		ctx,
		query,
		userID,
	)

	return err
}

func (s *SessionRepository) DeleteExpiredSessions(
	ctx context.Context,
) error {

	query := `
		DELETE FROM sessions
		WHERE expires_at < NOW()
	`

	_, err := s.pool.Exec(
		ctx,
		query,
	)

	return err
}
