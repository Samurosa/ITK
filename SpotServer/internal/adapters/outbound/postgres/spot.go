package postgres

import (
	"ITK_Code/m/v2/internal/core/dto"
	errorsCore "ITK_Code/m/v2/internal/core/errors"
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) Save(ctx context.Context, spot dto.CreateSpot) (string, error) {

	var spotID string

	query := `
		INSERT INTO spot
		(
			 symbol,
			 base_asset,
			 quote_asset,
			 price_precision,
			 quantity_precision,
			 min_order_size,
			 max_order_size,
			 allowed_roles,
			 name,
			 description,
			 status
		 )
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)
		ON CONFLICT (base_asset, quote_asset) 
		DO NOTHING
		RETURNING id
	`

	err := s.pool.QueryRow(ctx,
		query,
		spot.Symbol,
		spot.BaseAsset,
		spot.QuoteAsset,
		spot.PricePrecision,
		spot.QuantityPrecision,
		spot.MinOrderSize,
		spot.MaxOrderSize,
		spot.AllowedRoles,
		spot.Name,
		spot.Description,
		string(dto.ActiveStatus),
	).Scan(
		&spotID,
	)
	if errors.Is(err, pgx.ErrNoRows) {
		err = s.pool.QueryRow(ctx,
			`
			SELECT id FROM spot WHERE symbol = $1
		`,
			spot.Symbol,
		).Scan(
			&spotID,
		)
		if err != nil {
			return "", err
		}
		return spotID, nil
	}
	if err != nil {
		return "", err
	}
	return spotID, nil
}

func (s *Storage) Get(ctx context.Context, spotID string) (dto.Spot, error) {

	var spot dto.Spot

	err := s.pool.QueryRow(ctx,
		`
		SELECT
		    id,
			symbol,
			base_asset,
			quote_asset,
			price_precision,
			quantity_precision,
			min_order_size,
			max_order_size,
			allowed_roles,
			name,
			description,
			status,
			created_at,
			updated_at,
			disabled_at
		FROM spot
		WHERE id = $1
	`,
		spotID,
	).Scan(
		&spot.ID,
		&spot.Symbol,
		&spot.BaseAsset,
		&spot.QuoteAsset,
		&spot.PricePrecision,
		&spot.QuantityPrecision,
		&spot.MinOrderSize,
		&spot.MaxOrderSize,
		&spot.AllowedRoles,
		&spot.Name,
		&spot.Description,
		&spot.Status,
		&spot.CreatedAt,
		&spot.UpdatedAt,
		&spot.DisabledAt,
	)
	if err != nil {
		return dto.Spot{}, err
	}
	return spot, nil
}

func (s *Storage) Enable(ctx context.Context, spotID string) error {
	query :=
		`
	UPDATE spot
	SET 
	    status = $1,
	    updated_at = now(),
	    disabled_at = null
	WHERE id = $2
`

	result, err := s.pool.Exec(ctx,
		query,
		string(dto.ActiveStatus),
		spotID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errorsCore.ErrSpotNotFound
	}

	return nil
}
func (s *Storage) Disable(ctx context.Context, spotID string) error {
	query :=
		`
	UPDATE spot
	SET 
	    status = $1,
	    updated_at = now(),
	    disabled_at = now()
	WHERE id = $2
`

	result, err := s.pool.Exec(ctx,
		query,
		string(dto.DisabledStatus),
		spotID,
	)

	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errorsCore.ErrSpotNotFound
	}

	return nil
}
