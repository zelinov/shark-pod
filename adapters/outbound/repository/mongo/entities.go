package mongo

import (
	"time"
)

type TokenDocument struct {
	Id        string    `bson:"id"`
	Pk        string    `bson:"pk"`
	CreatedAt time.Time `bson:"created_at"`
}

type NonceDocument struct {
	PublicKey string    `bson:"public_key"`
	Nonce     string    `bson:"nonce"`
	CreatedAt time.Time `bson:"created_at"`
}

type MetadataDocument struct {
	PublicKey string    `bson:"public_key"`
	ImageUrl  string    `bson:"image_url"`
	Creators  []string  `bson:"creators"`
	IsNft     bool      `bson:"is_nft"`
	CreatedAt time.Time `bson:"created_at"`
}
