package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
	"go.uber.org/zap"
	"sharks/adapters/outbound/logger"
	"sharks/config"
	"time"
)

const JwtCollection = "jwtTokens"
const NonceCollection = "nonce"
const TokenMetadataCollection = "tokenMetadata"

func setIndexes() {
	conf := config.GetConfig()

	refreshTokensIndex := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{"id", bsonx.Int64(1)}},
		},
		{
			Keys:    bsonx.Doc{{"created_at", bsonx.Int64(1)}},
			Options: options.Index().SetExpireAfterSeconds(int32(conf.RefreshExpiration * 60)),
		},
	}

	opts := options.CreateIndexes().SetMaxTime(128 * time.Second)

	if _, err := DB.
		Collection(JwtCollection).
		Indexes().
		CreateMany(context.Background(), refreshTokensIndex, opts); err != nil {
		logger.Log.Fatal(
			fmt.Sprintf("create %s indexes error", JwtCollection),
			zap.Error(err),
		)
	}

	nonceIndex := []mongo.IndexModel{
		{
			Keys: bsonx.Doc{{"public_key", bsonx.Int64(1)}},
		},
		{
			Keys:    bsonx.Doc{{"created_at", bsonx.Int64(1)}},
			Options: options.Index().SetExpireAfterSeconds(120),
		},
	}

	if _, err := DB.
		Collection(NonceCollection).
		Indexes().
		CreateMany(context.Background(), nonceIndex, opts); err != nil {
		logger.Log.Fatal(
			fmt.Sprintf("create %s indexes error", NonceCollection),
			zap.Error(err),
		)
	}

	solanaTokenIndex := []mongo.IndexModel{
		{
			Keys:    bsonx.Doc{{"public_key", bsonx.Int64(1)}},
			Options: options.Index().SetUnique(true),
		},
	}

	if _, err := DB.
		Collection(TokenMetadataCollection).
		Indexes().
		CreateMany(context.Background(), solanaTokenIndex, opts); err != nil {
		logger.Log.Fatal(
			fmt.Sprintf("create %s indexes error", TokenMetadataCollection),
			zap.Error(err),
		)
	}
}
