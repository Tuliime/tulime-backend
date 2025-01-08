package models

import (
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var (
	gormDB *gorm.DB
	once   sync.Once
)

func Db() *gorm.DB {
	once.Do(func() {
		var dsn string
		var err error

		env := os.Getenv("GO_ENV")
		log.Println("GO_ENV:", env)

		// Load dev .env file
		// env := os.Getenv("GO_ENV")
		if env == "development" {
			err := godotenv.Load()
			if err != nil {
				log.Fatalf("Error loading .env file")
			}
			log.Println("Loaded .env var file")
		}
		log.Println("DB DSN: ", os.Getenv("TULIME_DEV_DSN"))
		log.Println("DB DSN: ", os.Getenv("TULIME_PROD_DSN"))

		switch env {
		case "development":
			dsn = os.Getenv("TULIME_DEV_DSN")
		case "production":
			dsn = os.Getenv("TULIME_PROD_DSN")
		default:
			log.Fatal("Unrecognized GO_ENV:", env)
		}

		gormDB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			SkipDefaultTransaction: true, PrepareStmt: true,
		})

		if err != nil {
			log.Fatal("Failed to connect to the database", err)
		}
		log.Println("Connected to postgres successfully")

		err = gormDB.AutoMigrate(&Agroproduct{}, &AgroproductPrice{})
		if err != nil {
			log.Fatal("Failed to make auto migration", err)
		}
		log.Println("Auto Migration successful")

	})

	return gormDB
}
