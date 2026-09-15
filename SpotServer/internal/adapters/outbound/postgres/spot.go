package postgres

import (
	"ITK_Code/m/v2/internal/adapters/outbound/encoding/cursor"
	errorsCore "ITK_Code/m/v2/internal/core/coreErrors"
	"ITK_Code/m/v2/internal/core/dto"
	"ITK_Code/m/v2/internal/core/spot/models"
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
)

func (s *Storage) Save(ctx context.Context, spot models.CreateSpot) (string, error) {

	var spotID string

	query := `
		INSERT INTO spot
		(
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
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)
		ON CONFLICT (base_asset, quote_asset) 
		DO NOTHING
		RETURNING id
	`

	err := s.pool.QueryRow(ctx,
		query,
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
			SELECT id FROM spot WHERE base_asset = $1 AND quote_asset = $2
		`,
			spot.BaseAsset,
			spot.QuoteAsset,
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

func (s *Storage) List(ctx context.Context, searchReq models.ListSpotsRequest) ([]dto.SpotListItem, string, bool, error) {

	baseQuery := `
		SELECT
		    id,
			base_asset,
			quote_asset,
			name,
			description,
			status,
			created_at
		FROM spot
`
	args := make([]any, 0, 5)
	argsPos := 1
	conditions := make([]string, 0, 4)

	if searchReq.Cursor != "" {
		gotCursor, err := cursor.DecodeCursor(searchReq.Cursor)
		if err != nil {
			return nil, "", false, err
		}

		conditions = append(
			conditions,
			fmt.Sprintf(
				"(created_at, id) < ($%d, $%d)",
				argsPos,
				argsPos+1,
			),
		)

		args = append(
			args,
			gotCursor.CreatedAt,
			gotCursor.ID,
		)

		argsPos += 2
	}

	if searchReq.Status != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("status = $%d", argsPos),
		)

		args = append(args, searchReq.Status)

		argsPos++
	}

	if searchReq.BaseAsset != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("base_asset = $%d", argsPos),
		)

		args = append(args, searchReq.BaseAsset)

		argsPos++
	}
	if searchReq.QuoteAsset != "" {
		conditions = append(
			conditions,
			fmt.Sprintf("quote_asset = $%d", argsPos),
		)

		args = append(args, searchReq.QuoteAsset)

		argsPos++
	}

	query := baseQuery

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY created_at DESC, id DESC " + fmt.Sprintf(" LIMIT $%d ", argsPos)

	limit := int(searchReq.PageSize) + 1
	if limit <= 1 {
		limit = 2
	}

	args = append(args, limit)

	rows, err := s.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, "", false, err
	}
	defer rows.Close()

	spots := make([]dto.SpotListItem, 0, limit)

	for rows.Next() {

		var spot dto.SpotListItem

		err = rows.Scan(
			&spot.ID,
			&spot.BaseAsset,
			&spot.QuoteAsset,
			&spot.Name,
			&spot.Description,
			&spot.Status,
			&spot.CreatedAt,
		)
		if err != nil {
			return nil, "", false, err
		}

		spots = append(spots, spot)
	}

	if err = rows.Err(); err != nil {
		return nil, "", false, err
	}

	if len(spots) == 0 {
		return spots, "", false, nil
	}

	hasMore := len(spots) > int(searchReq.PageSize)

	if hasMore {
		spots = spots[:searchReq.PageSize]

		newCursor, err := cursor.EncodeCursor(
			cursor.SpotCursor{
				CreatedAt: spots[len(spots)-1].CreatedAt,
				ID:        spots[len(spots)-1].ID,
			},
		)
		if err != nil {
			return nil, "", false, err
		}

		return spots, newCursor, true, nil
	}

	return spots, "", false, nil
}
