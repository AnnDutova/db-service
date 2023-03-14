package model

import "github.com/google/uuid"

type Token struct {
	IDToken
	RefreshToken
}

type IDToken struct {
	SS string `json:"idToken"`
}

type RefreshToken struct {
	ID  uuid.UUID `json:"-"`
	UID uuid.UUID `json:"-"`
	SS  string    `json:"refreshToken"`
}