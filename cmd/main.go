package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Tuliime/tulime-backend/internal/handlers/agroproducts"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"
)

func main() {
	app := fiber.New(fiber.Config{
		ErrorHandler: packages.DefaultErrorHandler,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	app.Use(logger.New())

	// Load dev .env file
	env := os.Getenv("GO_ENV")
	if env == "development" {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file")
		}
		log.Println("Loaded .env var file")
	}

	agroProducts := app.Group("/api/v0.01/agroproducts", func(c *fiber.Ctx) error {
		return c.Next()
	})

	// Agroproduct
	agroProducts.Get("/", agroproducts.GetAllProducts)
	agroProducts.Post("/", agroproducts.PostAgroProduct)
	// AgroproductPrices
	agroProducts.Get("/prices", agroproducts.GetAllAgroProductPrices)
	agroProducts.Post("/:id/price", agroproducts.PostAgroProductPrice)
	agroProducts.Get("/:id/price", agroproducts.GetPricesByAgroProduct)
	agroProducts.Patch("/:id/price/:priceId", agroproducts.UpdateAgroProductPrice)
	agroProducts.Delete("/:id/price/:priceId", agroproducts.DeleteAgroProductPrice)

	app.Use("*", func(c *fiber.Ctx) error {
		message := fmt.Sprintf("api route '%s' doesn't exist!", c.Path())
		return fiber.NewError(fiber.StatusNotFound, message)
	})

	log.Fatal(app.Listen(":5000"))
}
