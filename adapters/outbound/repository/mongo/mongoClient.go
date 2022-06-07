package mongo

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.uber.org/zap"
	"sharks/adapters/outbound/logger"
	"sharks/config"
)

var Client *mongo.Client
var DB *mongo.Database

func NewMongoClient() *mongo.Database {
	conf := config.GetConfig()
	connectionString := fmt.Sprintf("mongodb://%s:%d", conf.MongoDBHost, conf.MongoDBPort)
	credential := options.Credential{
		AuthSource: "admin",
		Username:   conf.MongoDBUser,
		Password:   conf.MongoDBPassword,
	}
	clientOptions := options.
		Client().
		ApplyURI(connectionString).
		SetAuth(credential)

	var err error

	if Client, err = mongo.Connect(context.Background(), clientOptions); err != nil {
		logger.Log.Fatal("mongo connect error", zap.Error(err))
	}

	if err = Client.Ping(context.Background(), readpref.Primary()); err != nil {
		logger.Log.Fatal("mongo connect error", zap.Error(err))
	}

	DB = Client.Database(conf.MongoDBDatabaseName)

	setIndexes()

	return DB
}

func CloseMongoClient() {
	if Client != nil {
		if err := Client.Disconnect(context.Background()); err != nil {
			logger.Log.Fatal("mongo disconnect error", zap.Error(err))
		}
	}
}
