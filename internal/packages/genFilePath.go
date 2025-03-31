package packages

import (
	"fmt"
	"log"
	"os"
)

// GenFilePath function generates file path for assets
// to be stored in firebase storage bucket
func GenFilePath(filename string) string {
	env := os.Getenv("GO_ENV")
	var appEnv string

	switch env {
	case "development":
		appEnv = "dev"
	case "production":
		appEnv = "prod"
	default:
		log.Fatal("Unrecognized GO_ENV:", env)
	}

	filePath := fmt.Sprintf("tulime/%s/%s", appEnv, GenerateFilename())

	return filePath
}
