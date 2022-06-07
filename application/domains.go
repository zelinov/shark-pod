package application

import (
	"github.com/google/uuid"
	"time"
)

type Nonce struct {
	PublicKey string
	Nonce     uuid.UUID
}

type Credentials struct {
	Signature string
	PublicKey string
}

type TokenFilter struct {
	Owner    string
	Creators *[]string
}

type JwtToken struct {
	Access  Token
	Refresh Token
}

type Token struct {
	Id        uuid.UUID
	Key       string
	PublicKey string
	Expired   time.Time
}

type TokenMetadata struct {
	PublicKey string
	ImageUrl  string
	Creators  []string
	IsNft     bool
}
