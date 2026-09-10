package cursor

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SpotCursor struct {
	CreatedAt time.Time `json:"created_at"`
	ID        string    `json:"id"`
}

func EncodeCursor(cursor SpotCursor) (string, error) {
	data, err := json.Marshal(cursor)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(data), nil
}

func DecodeCursor(cursor string) (SpotCursor, error) {

	data, err := base64.RawURLEncoding.DecodeString(cursor)
	if err != nil {
		return SpotCursor{}, err
	}

	var spotCursor SpotCursor

	if err = json.Unmarshal(data, &spotCursor); err != nil {
		return SpotCursor{}, err
	}

	if spotCursor.ID == "" {
		return SpotCursor{}, ErrInvalidCursor
	}

	if spotCursor.CreatedAt.IsZero() {
		return SpotCursor{}, ErrInvalidCursor
	}

	if _, err := uuid.Parse(spotCursor.ID); err != nil {
		return SpotCursor{}, ErrInvalidCursor
	}

	return spotCursor, nil
}
