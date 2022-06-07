package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

const (
	goMinProcs = 2

	httpServerPortEnv      = "SERVICE_PORT"
	accessSecretEnv        = "ACCESS_SECRET"
	accessExpirationEnv    = "ACCESS_EXPIRATION_MINUTES"
	refreshSecretEnv       = "REFRESH_SECRET"
	refreshExpirationEnv   = "REFRESH_EXPIRATION_MINUTES"
	creatorAccountEnv      = "CREATOR_ACCOUNT_PUBLIC_KEY"
	creatorNftCountEnv     = "CREATOR_NFT_COUNT"
	solanaRpcPoolUrlEnv    = "SOLANA_RPC_POOL_URL"
	mongoDBHostEnv         = "MONGO_HOST"
	mongoDBPortEnv         = "MONGO_PORT"
	mongoDBUserEnv         = "MONGO_USERNAME"
	mongoDBPasswordEnv     = "MONGO_PASSWORD"
	mongoDBDatabaseNameEnv = "MONGO_DATABASE"
)

type Configuration struct {
	HttpPort            string
	AccessSecret        []byte
	AccessExpiration    int
	RefreshSecret       []byte
	RefreshExpiration   int
	CreatorAccounts     []string
	CreatorNftCount     int
	SolanaRpcPoolUrl    string
	MongoDBHost         string
	MongoDBPort         uint64
	MongoDBUser         string
	MongoDBPassword     string
	MongoDBDatabaseName string
	LogLevel            string
}

var config *Configuration

func init() {
	if runtime.GOMAXPROCS(0) < goMinProcs {
		runtime.GOMAXPROCS(goMinProcs)
	}

	log.Println("GOMAXPROCS=" + strconv.Itoa(runtime.GOMAXPROCS(0)))
	log.Println("NUMCPU=" + strconv.Itoa(runtime.NumCPU()))

	if godotenv.Load(".env") != nil {
		log.Fatalf("Error loading .env file")
	}

	port := getStringEnv(httpServerPortEnv, "8080")
	accessSecret := []byte(os.Getenv(accessSecretEnv))
	accessExpiration := getUintEnv(accessExpirationEnv, 15)
	refreshSecret := []byte(os.Getenv(refreshSecretEnv))
	refreshExpiration := getUintEnv(refreshExpirationEnv, 10*60)
	creatorAccounts := getStringArrayEnv(creatorAccountEnv, []string{})
	creatorNftCount := getUintEnv(creatorNftCountEnv, 0)
	solanaRpcPoolUrl := getStringEnv(solanaRpcPoolUrlEnv, "")
	mongoDBHost := getStringEnv(mongoDBHostEnv, "localhost")
	mongoDBPort := getUintEnv(mongoDBPortEnv, 27017)
	mongoDBUser := getStringEnv(mongoDBUserEnv, "")
	mongoDBPassword := getStringEnv(mongoDBPasswordEnv, "")
	mongoDBDatabaseName := getStringEnv(mongoDBDatabaseNameEnv, "")
	logLevel := "info"

	config = &Configuration{
		HttpPort:            port,
		AccessSecret:        accessSecret,
		AccessExpiration:    int(accessExpiration),
		RefreshSecret:       refreshSecret,
		RefreshExpiration:   int(refreshExpiration),
		CreatorAccounts:     creatorAccounts,
		CreatorNftCount:     int(creatorNftCount),
		SolanaRpcPoolUrl:    solanaRpcPoolUrl,
		MongoDBHost:         mongoDBHost,
		MongoDBPort:         mongoDBPort,
		MongoDBUser:         mongoDBUser,
		MongoDBPassword:     mongoDBPassword,
		MongoDBDatabaseName: mongoDBDatabaseName,
		LogLevel:            logLevel,
	}
}

func GetConfig() *Configuration {
	return config
}

func getStringEnv(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)

	if !exists {
		value = defaultValue
	}

	log.Println(key + "=" + value)

	return value
}

func getUintEnv(key string, defaultValue uint64) uint64 {
	var uintValue uint64
	value, exists := os.LookupEnv(key)

	if !exists {
		uintValue = defaultValue
	} else {
		uintValue, _ = strconv.ParseUint(value, 10, 64)
	}

	log.Println(key + "=" + strconv.FormatUint(uintValue, 10))

	return uintValue
}

func getStringArrayEnv(key string, defaultValue []string) []string {
	var arrayValue []string
	value, exists := os.LookupEnv(key)

	if !exists {
		arrayValue = defaultValue
	} else {
		for _, s := range strings.Split(value, ",") {
			arrayValue = append(arrayValue, strings.TrimSpace(s))
		}
	}

	fmt.Printf("%s=%v\n", key, arrayValue)

	return arrayValue
}
