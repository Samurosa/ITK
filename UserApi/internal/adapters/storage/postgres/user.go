package postgres

import (
	userCore "ITK_Code/m/v2/internal/core/user"
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailIsExist = errors.New("email is exist")
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
			updated_at
		)
		VALUES($1,$2,$3,$4,$5,$6)
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
	).Scan(&id)

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
		return userCore.User{}, userCore.ErrUserNotFound
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
		return userCore.User{}, userCore.ErrUserNotFound
	}

	return userModel, nil
}

func (r *UserRepository) IsExistsUserByEmail(ctx context.Context,
	email string,
) bool {

	query := `
		SELECT EXISTS(
			SELECT 1
			FROM users
			WHERE email=$1
		)
	`

	var exists bool

	err := r.pool.QueryRow(
		ctx,
		query,
		email,
	).Scan(&exists)

	if err != nil {
		return false
	}

	return exists
}

func (r *UserRepository) Update(ctx context.Context, userID string, update userCore.UpdateUser) (bool, error) {

	query := `
		UPDATE users
		SET
			name = COALESCE($1, name),
			email = COALESCE($2, email),
			role = COALESCE($3, role),
			updated_at = NOW()
		WHERE id = $4
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		update.Name,
		update.Email,
		update.Role,
		userID,
	)

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		if pgErr.Code == "23505" {
			return false, ErrEmailIsExist
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

func (r *UserRepository) UpdatePassword(ctx context.Context, current userCore.User, update userCore.UpdateUser) (bool, error) {

	if update.PassHash == nil {
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
		*update.PassHash,
		current.ID,
	)

	if err != nil {
		return false, err
	}

	if result.RowsAffected() == 0 {
		return false, userCore.ErrUserNotFound
	}

	return true, nil
}

func (r *UserRepository) Delete(ctx context.Context,
	uid string,
) error {

	query := `
		DELETE FROM users
		WHERE id=$1
	`

	result, err := r.pool.Exec(
		ctx,
		query,
		uid,
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
	uid string,
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
		uid,
	).Scan(&role)

	if err != nil {
		return false, userCore.ErrUserNotFound
	}

	return role == userCore.AdminRole, nil
}
