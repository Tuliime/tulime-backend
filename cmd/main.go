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
	"github.com/Tuliime/tulime-backend/internal/handlers/device"
	"github.com/Tuliime/tulime-backend/internal/handlers/farminputs"
	"github.com/Tuliime/tulime-backend/internal/handlers/farmmanager"
	"github.com/Tuliime/tulime-backend/internal/handlers/news"
	"github.com/Tuliime/tulime-backend/internal/handlers/notification"
	"github.com/Tuliime/tulime-backend/internal/handlers/status"
	"github.com/Tuliime/tulime-backend/internal/handlers/vetdoctor"
	"github.com/Tuliime/tulime-backend/internal/middlewares"
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
		AllowOrigins:  "*",
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization",
		ExposeHeaders: "Content-Length",
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
	userGroup.Post("/auth/rt-signin", auth.SignInWithRefreshToken)
	userGroup.Post("/auth/forgot-password", auth.ForgotPassword)
	userGroup.Patch("/auth/verify-otp", auth.VerifyOTP)
	userGroup.Patch("/auth/reset-password/:otp", auth.ResetPassword)
	userGroup.Patch("/:id/auth/change-password", middlewares.Auth, auth.ChangePassword)
	userGroup.Patch("/:id/image", middlewares.Auth, auth.UpdateUserImage)
	userGroup.Patch("/:id", middlewares.Auth, auth.UpdateUser)
	userGroup.Get("/", middlewares.Auth, auth.GetAllUsers)

	// Agroproduct
	agroproductsGroup := app.Group("/api/v0.01/agroproducts", func(c *fiber.Ctx) error {
		return c.Next()
	})
	agroproductsGroup.Get("/", agroproducts.GetAllAgroProducts)
	agroproductsGroup.Post("/", middlewares.Auth, agroproducts.PostAgroProduct)
	agroproductsGroup.Get("/:id", agroproducts.GetAgroProduct)
	agroproductsGroup.Patch("/:id", middlewares.Auth, agroproducts.UpdateAgroProduct)
	agroproductsGroup.Delete("/:id", middlewares.Auth, agroproducts.DeleteAgroProduct)
	agroproductsGroup.Patch("/:id/image", middlewares.Auth, agroproducts.UpdateAgroProductImage)
	// AgroproductPrices
	agroproductsGroup.Get("/prices", agroproducts.GetAllAgroProductPrices)
	agroproductsGroup.Post("/:id/price", middlewares.Auth, agroproducts.PostAgroProductPrice)
	agroproductsGroup.Get("/:id/price", agroproducts.GetPricesByAgroProduct)
	agroproductsGroup.Patch("/:id/price/:priceID", middlewares.Auth, agroproducts.UpdateAgroProductPrice)
	agroproductsGroup.Delete("/:id/price/:priceID", middlewares.Auth, agroproducts.DeleteAgroProductPrice)

	// News
	newsGroup := app.Group("/api/v0.01/news", func(c *fiber.Ctx) error {
		return c.Next()
	})
	newsGroup.Get("/", news.GetAllNews)
	newsGroup.Get("/:id", news.GetNews)
	newsGroup.Post("/", middlewares.Auth, news.PostNews)
	newsGroup.Patch("/:id", middlewares.Auth, news.UpdateNews)
	newsGroup.Patch("/:id/image", middlewares.Auth, news.UpdateNewsImage)
	newsGroup.Delete("/:id", middlewares.Auth, news.DeleteNews)

	// Farm inputs
	farmInputGroup := app.Group("/api/v0.01/farminputs", func(c *fiber.Ctx) error {
		return c.Next()
	})
	farmInputGroup.Get("/", farminputs.GetAllFarmInputs)
	farmInputGroup.Get("/:id", farminputs.GetFarmInput)
	farmInputGroup.Post("/", middlewares.Auth, farminputs.PostFarmInputs)
	farmInputGroup.Patch("/:id", middlewares.Auth, farminputs.UpdateFarmInput)
	farmInputGroup.Patch("/:id/image", middlewares.Auth, farminputs.UpdateFarmInputImage)
	farmInputGroup.Delete("/:id", middlewares.Auth, farminputs.DeleteFarmInput)

	// Farm manager
	farmManagerGroup := app.Group("/api/v0.01/farmmanager", func(c *fiber.Ctx) error {
		return c.Next()
	})
	farmManagerGroup.Get("/", farmmanager.GetAllFarmManagers)
	farmManagerGroup.Get("/:id", farmmanager.GetFarmManager)
	farmManagerGroup.Get("/user/:userID", farmmanager.GetFarmManagerByUser)
	farmManagerGroup.Post("/user/:userID", middlewares.Auth, farmmanager.PostFarmManager)
	farmManagerGroup.Patch("/:id", middlewares.Auth, farmmanager.UpdateFarmManager)
	farmManagerGroup.Delete("/:id", middlewares.Auth, farmmanager.DeleteFarmManager)

	// Vet Doctor
	vetDoctorGroup := app.Group("/api/v0.01/vetdoctor", func(c *fiber.Ctx) error {
		return c.Next()
	})
	vetDoctorGroup.Get("/", vetdoctor.GetAllVetDoctors)
	vetDoctorGroup.Get("/:id", vetdoctor.GetVetDoctor)
	vetDoctorGroup.Get("/user/:userID", vetdoctor.GetVetDoctorByUser)
	vetDoctorGroup.Post("/user/:userID", middlewares.Auth, vetdoctor.PostVetDoctorManager)
	vetDoctorGroup.Patch("/:id", middlewares.Auth, vetdoctor.UpdateVetDoctor)
	vetDoctorGroup.Delete("/:id", middlewares.Auth, vetdoctor.DeleteVetDoctor)

	// ChatRoom
	chatRoomGroup := app.Group("/api/v0.01/chatroom", func(c *fiber.Ctx) error {
		return c.Next()
	})
	chatRoomGroup.Get("/", middlewares.Auth, chatroom.GetChat)
	chatRoomGroup.Post("/", middlewares.Auth, chatroom.PostChat)
	chatRoomGroup.Get("/live", middlewares.Auth, chatroom.GetLiveChat)
	chatRoomGroup.Get("/onlinestatus", middlewares.Auth, chatroom.GetOnlineStatus)
	chatRoomGroup.Patch("/onlinestatus", middlewares.Auth, chatroom.UpdateOnlineStatus)
	chatRoomGroup.Patch("/typingstatus", middlewares.Auth, chatroom.UpdateTypingStatus)

	// ChatBoot
	chatBotGroup := app.Group("/api/v0.01/chatbot", func(c *fiber.Ctx) error {
		return c.Next()
	})
	chatBotGroup.Get("/user/:userID", middlewares.Auth, chatbot.GetChatByUser)
	chatBotGroup.Post("/user/:userID", middlewares.Auth, chatbot.PostChat)
	chatBotGroup.Delete("/:id", middlewares.Auth, chatbot.DeleteChat)

	// Device
	deviceGroup := app.Group("/api/v0.01/device", func(c *fiber.Ctx) error {
		return c.Next()
	})
	deviceGroup.Post("/", middlewares.Auth, device.PostDevice)
	deviceGroup.Get("/user/:userID", middlewares.Auth, device.GetDeviceByUser)
	deviceGroup.Patch("/disable/:id", middlewares.Auth, device.DisableDevice)
	deviceGroup.Patch("/enable/:id", middlewares.Auth, device.EnableDevice)

	// Notification
	notificationGroup := app.Group("/api/v0.01/notification", func(c *fiber.Ctx) error {
		return c.Next()
	})
	notificationGroup.Get("/user/:userID", middlewares.Auth, notification.GetNotificationByUser)
	notificationGroup.Get("/live", middlewares.Auth, notification.GetLiveNotification)
	notificationGroup.Patch("/:id", middlewares.Auth, notification.UpdateNotificationAsRead)

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
