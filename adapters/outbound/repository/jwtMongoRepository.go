package repository

import (
	"fmt"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	mng "sharks/adapters/outbound/repository/mongo"
	"sharks/application"
	"time"
)

type JwtMongoRepository struct {
	db *mongo.Database
}

func NewJwtMongoRepository(db *mongo.Database) *JwtMongoRepository {
	return &JwtMongoRepository{db}
}

func (r *JwtMongoRepository) Get(id uuid.UUID) (token *application.Token, err error) {
	filter := bson.M{"id": id.String()}
	d := mng.TokenDocument{}

	if cur := r.db.Collection(mng.JwtCollection).FindOne(ctx(), filter); cur.Err() == nil {
		if err = cur.Decode(&d); err == nil {
			token = &application.Token{
				Id:        uuid.MustParse(d.Id),
				PublicKey: d.Pk,
			}
		}
	} else {
		err = cur.Err()
	}

	return
}

func (r *JwtMongoRepository) Save(token *application.Token) error {
	d := mng.TokenDocument{
		Id:        token.Id.String(),
		Pk:        token.PublicKey,
		CreatedAt: time.Now(),
	}

	_, err := r.db.Collection(mng.JwtCollection).InsertOne(ctx(), d)

	return err
}

func (r *JwtMongoRepository) Delete(id uuid.UUID) error {
	filter := bson.M{"id": id.String()}

	if res, err := r.db.Collection(mng.JwtCollection).DeleteOne(ctx(), filter); err != nil {
		return err
	} else if res.DeletedCount == 0 {
		return fmt.Errorf("unauthorized")
	}

	return nil
}

func (r *JwtMongoRepository) DeleteAllByPublicKey(pk string) error {
	filter := bson.M{"pk": pk}

	_, err := r.db.Collection(mng.JwtCollection).DeleteMany(ctx(), filter)

	return err
}
