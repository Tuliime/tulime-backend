package packages

import (
	"log"
	"math/rand"
	"os"
	"strconv"
)

// GenFilePath function generates file path for assets
// to be stored in firebase storage bucket
func GenFilePath(filename string) string {
	randNumStr := strconv.Itoa(rand.Intn(9000) + 1000)
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

	filePath := "tulime/" + appEnv + "/" + randNumStr + "_" + filename

	return filePath
}
