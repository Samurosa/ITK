package models

import "time"

type App struct {
	Secret      []byte
	TokensPairs Tokens
}

type Tokens struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresAt  time.Time
	RefreshExpiresAt time.Time
}
