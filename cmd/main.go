package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Tuliime/tulime-backend/internal/events/subscribers"
	"github.com/Tuliime/tulime-backend/internal/handlers/agroproducts"
	"github.com/Tuliime/tulime-backend/internal/handlers/auth"
	"github.com/Tuliime/tulime-backend/internal/handlers/chatbot"
	"github.com/Tuliime/tulime-backend/internal/handlers/chatroom"
	"github.com/Tuliime/tulime-backend/internal/handlers/farminputs"
	"github.com/Tuliime/tulime-backend/internal/handlers/farmmanager"
	"github.com/Tuliime/tulime-backend/internal/handlers/news"
	"github.com/Tuliime/tulime-backend/internal/handlers/status"
	"github.com/Tuliime/tulime-backend/internal/handlers/vetdoctor"
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

	// auth
	userGroup := app.Group("/api/v0.01/user", func(c *fiber.Ctx) error {
		return c.Next()
	})
	userGroup.Post("/auth/signup", auth.SignUp)
	userGroup.Post("/auth/signin", auth.SignIn)
	userGroup.Post("/auth/forgot-password", auth.ForgotPassword)
	userGroup.Patch("/auth/verify-otp", auth.VerifyOTP)
	userGroup.Patch("/auth/reset-password/:otp", auth.ResetPassword)
	userGroup.Patch("/:id/auth/change-password", auth.ChangePassword)
	userGroup.Patch("/:id/image", auth.UpdateUserImage)
	userGroup.Patch("/:id", auth.UpdateUser)

	// Agroproduct
	agroproductsGroup := app.Group("/api/v0.01/agroproducts", func(c *fiber.Ctx) error {
		return c.Next()
	})
	agroproductsGroup.Get("/", agroproducts.GetAllAgroProducts)
	agroproductsGroup.Post("/", agroproducts.PostAgroProduct)
	agroproductsGroup.Get("/:id", agroproducts.GetAgroProduct)
	agroproductsGroup.Patch("/:id", agroproducts.UpdateAgroProduct)
	agroproductsGroup.Delete("/:id", agroproducts.DeleteAgroProduct)
	agroproductsGroup.Patch("/:id/image", agroproducts.UpdateAgroProductImage)
	// AgroproductPrices
	agroproductsGroup.Get("/prices", agroproducts.GetAllAgroProductPrices)
	agroproductsGroup.Post("/:id/price", agroproducts.PostAgroProductPrice)
	agroproductsGroup.Get("/:id/price", agroproducts.GetPricesByAgroProduct)
	agroproductsGroup.Patch("/:id/price/:priceId", agroproducts.UpdateAgroProductPrice)
	agroproductsGroup.Delete("/:id/price/:priceId", agroproducts.DeleteAgroProductPrice)

	// News
	newsGroup := app.Group("/api/v0.01/news", func(c *fiber.Ctx) error {
		return c.Next()
	})
	newsGroup.Get("/", news.GetAllNews)
	newsGroup.Get("/:id", news.GetNews)
	newsGroup.Post("/", news.PostNews)
	newsGroup.Patch("/:id", news.UpdateNews)
	newsGroup.Patch("/:id/image", news.UpdateNewsImage)
	newsGroup.Delete("/:id", news.DeleteNews)

	// Farm inputs
	farmInputGroup := app.Group("/api/v0.01/farminputs", func(c *fiber.Ctx) error {
		return c.Next()
	})
	farmInputGroup.Get("/", farminputs.GetAllFarmInputs)
	farmInputGroup.Get("/:id", farminputs.GetFarmInput)
	farmInputGroup.Post("/", farminputs.PostFarmInputs)
	farmInputGroup.Patch("/:id", farminputs.UpdateFarmInput)
	farmInputGroup.Patch("/:id/image", farminputs.UpdateFarmInputImage)
	farmInputGroup.Delete("/:id", farminputs.DeleteFarmInput)

	// Farm manager
	farmManagerGroup := app.Group("/api/v0.01/farmmanager", func(c *fiber.Ctx) error {
		return c.Next()
	})
	farmManagerGroup.Get("/", farmmanager.GetAllFarmManagers)
	farmManagerGroup.Get("/:id", farmmanager.GetFarmManager)
	farmManagerGroup.Get("/user/:userId", farmmanager.GetFarmManagerByUser)
	farmManagerGroup.Post("/user/:userId", farmmanager.PostFarmManager)
	farmManagerGroup.Patch("/:id", farmmanager.UpdateFarmManager)
	farmManagerGroup.Delete("/:id", farmmanager.DeleteFarmManager)

	// Vet Doctor
	vetDoctorGroup := app.Group("/api/v0.01/vetdoctor", func(c *fiber.Ctx) error {
		return c.Next()
	})
	vetDoctorGroup.Get("/", vetdoctor.GetAllVetDoctors)
	vetDoctorGroup.Get("/:id", vetdoctor.GetVetDoctor)
	vetDoctorGroup.Get("/user/:userId", vetdoctor.GetVetDoctorByUser)
	vetDoctorGroup.Post("/user/:userId", vetdoctor.PostVetDoctorManager)
	vetDoctorGroup.Patch("/:id", vetdoctor.UpdateVetDoctor)
	vetDoctorGroup.Delete("/:id", vetdoctor.DeleteVetDoctor)

	// ChatRoom
	chatRoomGroup := app.Group("/api/v0.01/chatroom", func(c *fiber.Ctx) error {
		return c.Next()
	})
	chatRoomGroup.Get("/", chatroom.GetChat)
	chatRoomGroup.Post("/", chatroom.PostChat)
	chatRoomGroup.Get("/live", chatroom.GetLiveChat)

	// ChatBoot
	chatBotGroup := app.Group("/api/v0.01/chatbot", func(c *fiber.Ctx) error {
		return c.Next()
	})
	chatBotGroup.Get("/user/:userId", chatbot.GetChatByUser)
	chatBotGroup.Post("/user/:userId", chatbot.PostChat)
	chatBotGroup.Delete("/:id", chatbot.DeleteChat)

	// Status
	app.Get("/status", status.GetAppStatus)

	app.Use("*", func(c *fiber.Ctx) error {
		message := fmt.Sprintf("api route '%s' doesn't exist!", c.Path())
		return fiber.NewError(fiber.StatusNotFound, message)
	})

	// Initialize all event subscribers in the app
	subscribers.InitEventSubscribers()

	log.Fatal(app.Listen(":5000"))

}
