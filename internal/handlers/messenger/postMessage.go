package messenger

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/bimg"
)

var PostMessage = func(c *fiber.Ctx) error {
	messenger := models.Messenger{}
	messengerRoom := models.MessengerRoom{}

	messenger.MessengerRoomID = c.FormValue("messengerRoomID")
	messenger.SenderID = c.FormValue("senderID")
	messenger.RecipientID = c.FormValue("recipientID")
	messenger.Text = c.FormValue("text")
	messenger.Reply = c.FormValue("reply")
	sentAt := c.FormValue("sentAt")
	tag := c.FormValue("tag")

	var fileUploaded bool = true

	if messenger.SenderID == "" || messenger.RecipientID == "" || sentAt == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing senderID/recipientID/SentAt!")
	}

	parsedSentAt, err := time.Parse(time.RFC3339, sentAt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid sentAt format! Must be an ISO 8601 string.")
	}
	fmt.Printf("parsedSentAt: %v\n", parsedSentAt)
	messenger.SentAt = parsedSentAt

	var tags []string
	if tag != "" {
		err = json.Unmarshal([]byte(tag), &tags)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid tag format! Must be a JSON stringified array of strings.")
		}
	}
	fmt.Printf("tags: %v\n", tags)

	fileReader, err := c.FormFile("file")
	if err != nil {
		if err.Error() == constants.NO_FILE_UPLOADED_ERROR {
			fmt.Println("No file uploaded")
			// Prevent empty Text field when there is no file uploaded
			if messenger.Text == "" {
				return fiber.NewError(fiber.StatusBadRequest, "Missing Text field!")
			}
			fileUploaded = false

		} else {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	var filePath string
	var imageUrl string
	var messengerFile models.MessengerFile
	var messengerTags []models.MessengerTag
	dimensions := struct {
		Height int `json:"height"`
		Width  int `json:"Width"`
	}{Height: 0, Width: 0}

	if fileUploaded {
		// Validate file size (10 MB limit)
		const maxFileSize = 10 << 20 // 10 MB in bytes
		if fileReader.Size > maxFileSize {
			return fiber.NewError(fiber.StatusBadRequest, "File size exceeds the 10 MB limit")
		}
		file, err := fileReader.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		defer file.Close()

		// Get image dimensions
		buf, err := io.ReadAll(file)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		image := bimg.NewImage(buf)
		size, err := image.Size()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		dimensions.Height = size.Height
		dimensions.Width = size.Width

		filePath = packages.GenFilePath(fileReader.Filename)
		firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

		imageUrl, err = firebaseStorage.Add(file, fileReader)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	messenger.ArrivedAt = time.Now()

	//Get messengerRoom ID if not available
	if messenger.MessengerRoomID == "" {
		messengerRoom, err = messengerRoom.FindByUsers(messenger.SenderID, messenger.RecipientID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if messengerRoom.ID == "" {
			newMessengerRoom, err := messengerRoom.Create(models.MessengerRoom{
				UserOneID: messenger.SenderID,
				UserTwoID: messenger.RecipientID})
			if err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			messenger.MessengerRoomID = newMessengerRoom.ID
		} else {
			messenger.MessengerRoomID = messengerRoom.ID
		}
	}

	newMessage, err := messenger.Create(messenger)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Save messengerFile
	if imageUrl != "" {
		dimensionJson, err := json.Marshal(dimensions)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		messengerFile = models.MessengerFile{MessengerID: newMessage.ID,
			URL: imageUrl, Path: filePath, Dimensions: models.JSONB(dimensionJson)}
		newChatRoomFile, err := messengerFile.Create(messengerFile)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newMessage.File = newChatRoomFile
	}

	// Save all tags of the messenger
	for _, tag := range tags {
		if tag == "" {
			continue
		}
		messengerTag := models.MessengerTag{MessengerID: newMessage.ID, AdvertID: tag}
		newMessengerTag, err := messengerTag.Create(messengerTag)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		messengerTags = append(messengerTags, newMessengerTag)
	}

	log.Printf("messengerTags: %v\n", messengerTags)

	newMessage.Tag = messengerTags

	events.EB.Publish("messenger", newMessage)
	events.EB.Publish("messengerNotification", newMessage)

	response := fiber.Map{
		"status":  "success",
		"message": "message posted successfully!",
		"data":    newMessage,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
