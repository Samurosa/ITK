package postgres

import (
	"ITK_Code/m/v2/internal/core/dto"
	userCore "ITK_Code/m/v2/internal/core/errors"
	"context"
	"errors"
	"time"

	"github.com/Samurosa/exchange-contract/protobuf/gen/go/user"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type SpotRepository struct {
	pool *pgxpool.Pool
}

func NewUserStorage(pool *pgxpool.Pool) *SpotRepository {
	return &SpotRepository{
		pool: pool,
	}
}

func (r *SpotRepository) Save(ctx context.Context, spot dto.CreateSpot) (string, error) {

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

/*func (r *SpotRepository) Get(ctx context.Context) error {
	panic("implement me")
}

func (r *SpotRepository) Update(ctx context.Context) error {
	panic("implement me")
}*/

func (r *SpotRepository) Enable(ctx context.Context, spotID string) error {

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

func (r *SpotRepository) Disable(ctx context.Context, spotID string) error {

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
