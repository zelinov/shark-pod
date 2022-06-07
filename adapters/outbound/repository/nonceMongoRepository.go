package repository

import (
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mng "sharks/adapters/outbound/repository/mongo"
	"sharks/application"
	"time"
)

type NonceMongoRepository struct {
	db *mongo.Database
}

func NewNonceMongoRepository(db *mongo.Database) *NonceMongoRepository {
	return &NonceMongoRepository{db}
}

func (r *NonceMongoRepository) Get(id string) (nonce *application.Nonce, err error) {
	filter := bson.M{"public_key": id}
	d := mng.NonceDocument{}

	if cur := r.db.Collection(mng.NonceCollection).FindOneAndDelete(ctx(), filter); cur.Err() == nil {
		if err = cur.Decode(&d); err == nil {
			nonce = &application.Nonce{
				PublicKey: d.PublicKey,
				Nonce:     uuid.MustParse(d.Nonce),
			}
		}
	} else {
		err = cur.Err()
	}

	return
}

func (r *NonceMongoRepository) Save(nonce *application.Nonce) error {
	d := mng.NonceDocument{
		PublicKey: nonce.PublicKey,
		Nonce:     nonce.Nonce.String(),
		CreatedAt: time.Now(),
	}

	_, err := r.db.Collection(mng.NonceCollection).InsertOne(ctx(), d)

	return err
}
