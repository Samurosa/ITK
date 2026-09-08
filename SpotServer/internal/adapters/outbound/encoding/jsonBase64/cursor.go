package jsonBase64

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
		return SpotCursor{}, errors.New("invalid cursor")
	}

	if spotCursor.CreatedAt.IsZero() {
		return SpotCursor{}, errors.New("invalid cursor")
	}

	if _, err := uuid.Parse(spotCursor.ID); err != nil {
		return SpotCursor{}, errors.New("invalid cursor")
	}

	return spotCursor, nil
}
