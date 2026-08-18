package postgres

import (
	userCore "ITK_Code/m/v2/internal/core/user"
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		pool: pool,
	}
}

func (r *UserRepository) SaveUser(ctx context.Context,
	user userCore.User,
) (
	string,
	error,
) {
	var id string

	query := ` INSERT INTO users
		(
			email,
			name,
			password_hash,
			role,
			created_at,
			updated_at
		)
		VALUES($1,$2,$3,$4,$5,$6)
		ON CONFLICT (email) DO UPDATE SET
			deleted_at = NULL,
    		name = EXCLUDED.name,
    		password_hash = EXCLUDED.password_hash,
    		role = EXCLUDED.role,
    		updated_at = EXCLUDED.updated_at
		WHERE users.deleted_at IS NOT NULL
		RETURNING id
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
		user.CreateTime,
		user.UpdateTime,
	).Scan(&id)

	if errors.Is(err, pgx.ErrNoRows) {
		return "", userCore.ErrEmailIsExist
	}

	if err != nil {
		return "", err
	}

	return id, nil
}

func (r *UserRepository) Get(ctx context.Context,
	uid string,
) (
	userCore.User,
	error,
) {

	query := `
	SELECT
		id,
		email,
		name,
		password_hash,
		role,
		created_at,
		updated_at
	FROM users
	WHERE id=$1 
	AND deleted_at IS NULL; 
	`

	var userModel userCore.User

	err := r.pool.QueryRow(
		ctx,
		query,
		uid,
	).Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.Name,
		&userModel.PasswordHash,
		&userModel.Role,
		&userModel.CreateTime,
		&userModel.UpdateTime,
	)

	if err != nil {
		return userCore.User{}, err
	}

	return userModel, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context,
	email string,
) (
	userCore.User,
	error,
) {

	query := `
	SELECT
		id,
		email,
		name,
		password_hash,
		role,
		created_at,
		updated_at
	FROM users
	WHERE email=$1
	AND deleted_at IS NULL; 
	`

	var userModel userCore.User

	err := r.pool.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&userModel.ID,
		&userModel.Email,
		&userModel.Name,
		&userModel.PasswordHash,
		&userModel.Role,
		&userModel.CreateTime,
		&userModel.UpdateTime,
	)

	if err != nil {
		return userCore.User{}, err
	}

	return userModel, nil
}

func (r *UserRepository) Update(ctx context.Context, userID string, update userCore.UpdateUser) error {

	query := `
		UPDATE users
		SET
			name = COALESCE($1, name),
			updated_at = NOW()
		WHERE id = $2
		AND deleted_at IS NULL;
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		update.Name,
		userID,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return userCore.ErrEmailIsExist
		}
	}

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return userCore.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, current userCore.User, newPass string) error {

	query := `
		UPDATE users
		SET
			password_hash=$1,
			updated_at=NOW()
		WHERE id=$2
		AND deleted_at IS NULL;
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		newPass,
		current.ID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return userCore.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) Delete(ctx context.Context,
	id string,
) error {

	query := `
		UPDATE users
		SET
			deleted_at=$1
		WHERE id=$2
		AND deleted_at IS NULL;
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		time.Now(),
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return userCore.ErrUserNotFound
	}

	return nil
}

func (r *UserRepository) IsAdmin(ctx context.Context,
	id string,
) (
	bool,
	error,
) {

	query := `
		SELECT role
		FROM users
		WHERE id=$1
		AND deleted_at IS NULL; 
	`

	var role userCore.Role

	err := r.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(&role)

	if err != nil {
		return false, err
	}

	return role == userCore.AdminRole, nil
}
