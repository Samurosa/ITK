package authorization

import (
	"ITK_Code/m/v2/internal/domain/models"
	"errors"
)

type Secret struct {
	secret string
}

func NewSecret(secret string) *Secret {
	return &Secret{
		secret: secret,
	}
}
func (s *Secret) GetSecret() (models.App, error) {
	if s.secret == "" {
		return models.App{}, errors.New("no secret:" + s.secret)
	}

	return models.App{
		Secret: []byte(s.secret),
	}, nil
}
