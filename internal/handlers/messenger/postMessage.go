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
	messengerTag := models.MessengerTag{}

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

	fileHeader, err := c.FormFile("file")
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

	var filePath, imageUrl string
	var messengerTags []models.MessengerTag
	dimensions := models.ImageDimensions{Height: 0, Width: 0}
	var repliedMessage, repliedMessageFile any = nil, nil

	if fileUploaded {
		// Validate file size (10 MB limit)
		const maxFileSize = 10 << 20 // 10 MB in bytes
		if fileHeader.Size > maxFileSize {
			return fiber.NewError(fiber.StatusBadRequest, "File size exceeds the 10 MB limit")
		}
		file, err := fileHeader.Open()
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

		// Reset file cursor
		file.Seek(0, io.SeekStart)
		// Compress image file
		imageProcessor := packages.ImageProcessor{}
		compressedFileBuf, err := imageProcessor.CompressMultipartFile(file, 75)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to compress image")
		}

		filePath = packages.GenFilePath(fileHeader.Filename)
		firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

		imageUrl, err = firebaseStorage.AddFromBuffer(compressedFileBuf)
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
		messengerFile := models.MessengerFile{MessengerID: newMessage.ID,
			URL: imageUrl, Path: filePath, Dimensions: models.JSONB(dimensionJson)}
		newChatRoomFile, err := messengerFile.Create(messengerFile)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newMessage.File = newChatRoomFile
	}

	// Save all tags of the messenger
	if len(tags) > 0 {
		for _, tag := range tags {
			messengerTags = append(messengerTags,
				models.MessengerTag{MessengerID: newMessage.ID, AdvertID: tag})
		}
		messengerTags, err = messengerTag.CreateMany(messengerTags)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		log.Printf("messengerTags: %v\n", messengerTags)
	}
	newMessage.Tag = messengerTags

	// Get replied message if it exists
	if newMessage.Reply != "" {
		reply, err := messenger.FindReply(newMessage.Reply)
		if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if reply.File.ID != "" {
			var dimensions models.ImageDimensions
			if err := json.Unmarshal(reply.File.Dimensions, &dimensions); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			repliedMessageFile = File{
				ID:          reply.File.ID,
				MessengerID: reply.File.MessengerID,
				URL:         reply.File.URL,
				Path:        reply.File.Path,
				Dimensions:  dimensions,
				CreatedAt:   reply.File.CreatedAt,
				UpdatedAt:   reply.File.UpdatedAt,
			}
		}
		repliedMessage = Message{
			ID:              reply.ID,
			MessengerRoomID: messenger.MessengerRoomID,
			SenderID:        reply.SenderID,
			RecipientID:     reply.RecipientID,
			Text:            reply.Text,
			Reply:           reply.Reply,
			RepliedMessage:  nil,
			File:            repliedMessageFile,
			Tag:             reply.Tag,
			SentAt:          reply.SentAt,
			ArrivedAt:       reply.ArrivedAt,
			CreatedAt:       reply.CreatedAt,
			UpdatedAt:       reply.UpdatedAt,
			Sender: User{
				ID:             reply.Sender.ID,
				Name:           reply.Sender.Name,
				TelNumber:      reply.Sender.TelNumber,
				Role:           reply.Sender.Role,
				ImageUrl:       reply.Sender.ImageUrl,
				ImagePath:      reply.Sender.ImagePath,
				ProfileBgColor: reply.Sender.ProfileBgColor,
				ChatroomColor:  reply.Sender.ChatroomColor,
				CreatedAt:      reply.Sender.CreatedAt,
				UpdatedAt:      reply.Sender.UpdatedAt,
			},
			Recipient: User{
				ID:             reply.Recipient.ID,
				Name:           reply.Recipient.Name,
				TelNumber:      reply.Recipient.TelNumber,
				Role:           reply.Recipient.Role,
				ImageUrl:       reply.Recipient.ImageUrl,
				ImagePath:      reply.Recipient.ImagePath,
				ProfileBgColor: reply.Recipient.ProfileBgColor,
				ChatroomColor:  reply.Recipient.ChatroomColor,
				CreatedAt:      reply.Recipient.CreatedAt,
				UpdatedAt:      reply.Recipient.UpdatedAt,
			},
		}
	}

	newMessage, err = messenger.FindOne(newMessage.ID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// structure the message
	var messengerFile any = nil
	if newMessage.File.ID != "" {
		messengerFile = File{
			ID:          newMessage.File.ID,
			MessengerID: newMessage.File.MessengerID,
			URL:         newMessage.File.URL,
			Path:        newMessage.File.Path,
			Dimensions:  dimensions,
			CreatedAt:   newMessage.File.CreatedAt,
			UpdatedAt:   newMessage.File.UpdatedAt,
		}
	}

	message := Message{
		ID:              newMessage.ID,
		MessengerRoomID: newMessage.MessengerRoomID,
		SenderID:        newMessage.SenderID,
		RecipientID:     newMessage.RecipientID,
		Text:            newMessage.Text,
		Reply:           newMessage.Reply,
		RepliedMessage:  repliedMessage,
		File:            messengerFile,
		Tag:             newMessage.Tag,
		SentAt:          newMessage.SentAt,
		ArrivedAt:       newMessage.ArrivedAt,
		CreatedAt:       newMessage.CreatedAt,
		UpdatedAt:       newMessage.UpdatedAt,
		Sender: User{
			ID:             newMessage.Sender.ID,
			Name:           newMessage.Sender.Name,
			TelNumber:      newMessage.Sender.TelNumber,
			Role:           newMessage.Sender.Role,
			ImageUrl:       newMessage.Sender.ImageUrl,
			ImagePath:      newMessage.Sender.ImagePath,
			ProfileBgColor: newMessage.Sender.ProfileBgColor,
			ChatroomColor:  newMessage.Sender.ChatroomColor,
			CreatedAt:      newMessage.Sender.CreatedAt,
			UpdatedAt:      newMessage.Sender.UpdatedAt,
		},
		Recipient: User{
			ID:             newMessage.Recipient.ID,
			Name:           newMessage.Recipient.Name,
			TelNumber:      newMessage.Recipient.TelNumber,
			Role:           newMessage.Recipient.Role,
			ImageUrl:       newMessage.Recipient.ImageUrl,
			ImagePath:      newMessage.Recipient.ImagePath,
			ProfileBgColor: newMessage.Recipient.ProfileBgColor,
			ChatroomColor:  newMessage.Recipient.ChatroomColor,
			CreatedAt:      newMessage.Recipient.CreatedAt,
			UpdatedAt:      newMessage.Recipient.UpdatedAt,
		},
	}

	events.EB.Publish("messenger", message)
	events.EB.Publish("messengerNotification", newMessage)

	response := fiber.Map{
		"status":  "success",
		"message": "message posted successfully!",
		"data":    message,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
