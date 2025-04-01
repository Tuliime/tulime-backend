package chatroom

import (
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
	"github.com/h2non/bimg"
)

var PostChat = func(c *fiber.Ctx) error {
	chatroom := models.Chatroom{}
	chatroomMention := models.ChatroomMention{}
	user := models.User{}

	chatroom.UserID = c.FormValue("userID")
	chatroom.Text = c.FormValue("text")
	chatroom.Reply = c.FormValue("reply")
	sentAt := c.FormValue("sentAt")
	mention := c.FormValue("mention")

	var fileUploaded bool = true

	if chatroom.UserID == "" || sentAt == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing userID/SentAt!")
	}

	parsedSentAt, err := time.Parse(time.RFC3339, sentAt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid sentAt format! Must be an ISO 8601 string.")
	}
	chatroom.SentAt = parsedSentAt

	var mentions []string
	if mention != "" {
		err = json.Unmarshal([]byte(mention), &mentions)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid mention format! Must be a JSON stringified array of strings.")
		}
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		if err.Error() == constants.NO_FILE_UPLOADED_ERROR {
			fmt.Println("No file uploaded")
			// Prevent empty Text field when there is no file uploaded
			if chatroom.Text == "" {
				return fiber.NewError(fiber.StatusBadRequest, "Missing Text field!")
			}
			fileUploaded = false

		} else {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	var imageUrl, filePath string
	var chatroomFile models.ChatroomFile
	var chatroomMentions []models.ChatroomMention
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

	chatroom.ArrivedAt = time.Now()

	newChatroom, err := chatroom.Create(chatroom)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Save chatroomFile
	if imageUrl != "" {
		dimensionJson, err := json.Marshal(dimensions)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		chatroomFile = models.ChatroomFile{ChatroomID: newChatroom.ID,
			URL: imageUrl, Path: filePath, Dimensions: models.JSONB(dimensionJson)}
		newChatRoomFile, err := chatroomFile.Create(chatroomFile)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newChatroom.File = newChatRoomFile
	}

	// Save all mentions of the chat message
	if mention != "" {
		for _, mention := range mentions {
			if mention == "" {
				continue
			}
			chatroomMentions = append(chatroomMentions,
				models.ChatroomMention{ChatroomID: newChatroom.ID, UserID: mention})
		}
		newChatroomMentions, err := chatroomMention.CreateMany(chatroomMentions)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newChatroom.Mention = newChatroomMentions
	}

	// Get replied message if it exists
	if newChatroom.Reply != "" {
		reply, err := chatroom.FindReply(newChatroom.Reply)
		if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if reply.File.ID != "" {
			var dimensions models.ImageDimensions
			if err := json.Unmarshal(reply.File.Dimensions, &dimensions); err != nil {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
			repliedMessageFile = File{
				ID:         reply.File.ID,
				ChatroomID: reply.File.ChatroomID,
				URL:        reply.File.URL,
				Path:       reply.File.Path,
				Dimensions: dimensions,
				CreatedAt:  reply.File.CreatedAt,
				UpdatedAt:  reply.File.UpdatedAt,
				DeletedAt:  reply.File.DeletedAt,
			}
		}
		repliedMessage = Message{
			ID:             reply.ID,
			UserID:         reply.UserID,
			Text:           reply.Text,
			Reply:          reply.Reply,
			RepliedMessage: nil,
			File:           repliedMessageFile,
			Mention:        reply.Mention,
			SentAt:         reply.SentAt,
			ArrivedAt:      reply.ArrivedAt,
			CreatedAt:      reply.CreatedAt,
			UpdatedAt:      reply.UpdatedAt,
			DeletedAt:      reply.DeletedAt,
			User: User{
				ID:             reply.User.ID,
				Name:           reply.User.Name,
				TelNumber:      reply.User.TelNumber,
				Role:           reply.User.Role,
				ImageUrl:       reply.User.ImageUrl,
				ImagePath:      reply.User.ImagePath,
				ProfileBgColor: reply.User.ProfileBgColor,
				ChatroomColor:  reply.User.ChatroomColor,
				CreatedAt:      reply.User.CreatedAt,
				UpdatedAt:      reply.User.UpdatedAt,
			},
		}
	}

	user, err = user.FindOne(newChatroom.UserID)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// structure the message
	var chatMessageFile any = nil
	if newChatroom.File.ID != "" {
		chatMessageFile = File{
			ID:         newChatroom.File.ID,
			ChatroomID: newChatroom.File.ChatroomID,
			URL:        newChatroom.File.URL,
			Path:       newChatroom.File.Path,
			Dimensions: dimensions,
			CreatedAt:  newChatroom.File.CreatedAt,
			UpdatedAt:  newChatroom.File.UpdatedAt,
			DeletedAt:  newChatroom.File.DeletedAt,
		}
	}

	message := Message{
		ID:             newChatroom.ID,
		UserID:         newChatroom.UserID,
		Text:           newChatroom.Text,
		Reply:          newChatroom.Reply,
		RepliedMessage: repliedMessage,
		File:           chatMessageFile,
		Mention:        newChatroom.Mention,
		SentAt:         newChatroom.SentAt,
		ArrivedAt:      newChatroom.ArrivedAt,
		CreatedAt:      newChatroom.CreatedAt,
		UpdatedAt:      newChatroom.UpdatedAt,
		DeletedAt:      newChatroom.DeletedAt,
		User: User{
			ID:             user.ID,
			Name:           user.Name,
			TelNumber:      user.TelNumber,
			Role:           user.Role,
			ImageUrl:       user.ImageUrl,
			ImagePath:      user.ImagePath,
			ProfileBgColor: user.ProfileBgColor,
			ChatroomColor:  user.ChatroomColor,
			CreatedAt:      user.CreatedAt,
			UpdatedAt:      user.UpdatedAt,
		},
	}

	events.EB.Publish("chatroomMessage", message)
	events.EB.Publish("chatNotification", newChatroom)

	response := fiber.Map{
		"status":  "success",
		"message": "chat posted successfully!",
		"data":    message,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
