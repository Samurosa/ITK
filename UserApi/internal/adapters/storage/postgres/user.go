package postgres

import (
	userCore "ITK_Code/m/v2/internal/core/user"
	"context"
	"errors"

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

	query := `
		INSERT INTO users
		(
			email,
			name,
			password_hash,
			role,
			created_at,
			updated_at,
		 	deleted
		)
		VALUES($1,$2,$3,$4,$5,$6,$7)
		RETURNING id
	`

	var id string

	err := r.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.Name,
		user.PasswordHash,
		user.Role,
		user.CreateTime,
		user.UpdateTime,
		user.Deleted,
	).Scan(&id)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return "", userCore.ErrEmailIsExist
		}
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
		updated_at,
		deleted
	FROM users
	WHERE id=$1
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
		&userModel.Deleted,
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
		updated_at,
		deleted
	FROM users
	WHERE email=$1
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
		&userModel.Deleted,
	)

	if err != nil {
		return userCore.User{}, err
	}

	return userModel, nil
}

func (r *UserRepository) Update(ctx context.Context, userID string, update userCore.UpdateUser) (bool, error) {

	query := `
		UPDATE users
		SET
			name = COALESCE($1, name),
			email = COALESCE($2, email),
			updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		update.Name,
		update.Email,
		userID,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return false, userCore.ErrEmailIsExist
		}
	}

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, userCore.ErrUserNotFound
	}

	return true, nil
}

func (r *UserRepository) UpdatePassword(ctx context.Context, current userCore.User, newPass string) (bool, error) {

	if newPass == "" {
		return false, errors.New("password hash is empty")
	}

	query := `
		UPDATE users
		SET
			password_hash=$1,
			updated_at=NOW()
		WHERE id=$2
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		newPass,
		current.ID,
	)

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, err
	}

	return true, nil
}

func (r *UserRepository) Delete(ctx context.Context,
	id string,
) error {

	query := `
		UPDATE users
		SET
			deleted=$1
		WHERE id=$2
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		true,
		id,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return err
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
