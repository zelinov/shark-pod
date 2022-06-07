package repository

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sharks/adapters/outbound/logger"
	mng "sharks/adapters/outbound/repository/mongo"
	"sharks/application"
	"time"
)

type TokenMongoRepository struct {
	db *mongo.Database
}

func NewTokenMongoRepository(db *mongo.Database) *TokenMongoRepository {
	return &TokenMongoRepository{db}
}

func (r *TokenMongoRepository) FindByPublicKey(id string) (t *application.TokenMetadata) {
	filter := bson.M{"public_key": id}

	if cur := r.db.Collection(mng.TokenMetadataCollection).FindOne(ctx(), filter); cur.Err() == nil {
		entity := mng.MetadataDocument{}

		if err := cur.Decode(&entity); err == nil {
			t = mapEntityToTokenMetadata(entity)
		} else {
			logger.Log.Error(err.Error())
		}
	}

	return
}

func (r *TokenMongoRepository) FindAllByPublicKeys(pks []string) (res []*application.TokenMetadata) {
	filter := bson.M{
		"public_key": bson.M{"$in": pks},
	}

	if cur, err := r.db.Collection(mng.TokenMetadataCollection).Find(ctx(), filter); err == nil {
		defer cur.Close(ctx())

		for cur.Next(ctx()) {
			entity := mng.MetadataDocument{}

			if err = cur.Decode(entity); err == nil {
				logger.Log.Error(err.Error())
			}

			res = append(res, mapEntityToTokenMetadata(entity))
		}

		if err = cur.Err(); err != nil {
			logger.Log.Error(err.Error())
		}
	}

	return
}

func (r *TokenMongoRepository) Save(token *application.TokenMetadata) (err error) {
	d := mapTokenMetadataToEntity(token)
	d.CreatedAt = time.Now()

	if _, err = r.db.Collection(mng.TokenMetadataCollection).InsertOne(ctx(), d); err != nil {
		logger.Log.Error(err.Error())
	}

	return
}

func (r *TokenMongoRepository) SaveAll(tokensMetadata []*application.TokenMetadata) (err error) {
	var data []interface{}

	for _, t := range tokensMetadata {
		d := mapTokenMetadataToEntity(t)
		d.CreatedAt = time.Now()
		data = append(data, d)
	}

	if _, err = r.db.Collection(mng.TokenMetadataCollection).InsertMany(ctx(), data); err != nil {
		logger.Log.Error(err.Error())
	}

	return
}
