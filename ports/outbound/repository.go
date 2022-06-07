package outbound

import (
	"github.com/google/uuid"
	"sharks/application"
)

type JwtRepository interface {
	Get(id uuid.UUID) (*application.Token, error)
	Save(token *application.Token) error
	Delete(id uuid.UUID) error
	DeleteAllByPublicKey(pk string) error
}

type TokenRepository interface {
	FindByPublicKey(id string) *application.TokenMetadata
	Save(metadata *application.TokenMetadata) error
	SaveAll(tokensMetadata []*application.TokenMetadata) error
}

type NonceRepository interface {
	Get(id string) (*application.Nonce, error)
	Save(nonce *application.Nonce) error
}
