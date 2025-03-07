package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

func RedisClient() *redis.Client {
	var REDIS_URL string

	env := os.Getenv("GO_ENV")
	log.Println("GO_ENV:", env)

	// Load dev .env file
	if env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		log.Println("Loaded .env var file")
	}

	log.Println("DEV_REDIS_URL: ", os.Getenv("DEV_REDIS_URL"))
	log.Println("PROD_REDIS_URL: ", os.Getenv("PROD_REDIS_URL"))

	switch env {
	case "development":
		REDIS_URL = os.Getenv("DEV_REDIS_URL")
	case "production":
		REDIS_URL = os.Getenv("PROD_REDIS_URL")
	default:
		log.Fatal("Unrecognized GO_ENV:", env)
	}

	opt, err := redis.ParseURL(REDIS_URL)
	if err != nil {
		log.Fatal("Failed to connect to redis", err)
	}

	client := redis.NewClient(opt)
	log.Println("Connected to redis successfully")

	return client
}
