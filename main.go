package main

import (
	"fmt"
	"net/http"
	"sharks/adapters/inbound/controller"
	"sharks/adapters/inbound/controller/rest"
	"sharks/adapters/outbound/client"
	"sharks/adapters/outbound/logger"
	"sharks/adapters/outbound/repository"
	"sharks/adapters/outbound/repository/mongo"
	"sharks/application/service"
	"sharks/application/service/solana"
	"sharks/config"
)

func main() {
	db := mongo.NewMongoClient()

	jwtRepository := repository.NewJwtMongoRepository(db)
	tokenRepository := repository.NewTokenMongoRepository(db)
	nonceRepository := repository.NewNonceMongoRepository(db)

	solanaClient := client.NewSimpleSolanaClient()

	solanaService := solana.NewSolanaService(solanaClient, tokenRepository)
	tokenService := service.NewNftTokenService(tokenRepository, solanaService)
	authService := service.NewJwtAuthService(jwtRepository, nonceRepository, solanaService, tokenService)

	handler := controller.NewHandler(rest.NewHttpHandler(authService, tokenService), authService)

	if err := http.ListenAndServe(fmt.Sprintf(":%s", config.GetConfig().HttpPort), handler); err != nil {
		logger.Log.Fatal(err.Error())

		mongo.CloseMongoClient()

		logger.Log.Info("MongoClient closed")
	}
}
