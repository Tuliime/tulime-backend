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
	chatRoom := models.Chatroom{}

	chatRoom.UserID = c.FormValue("userID")
	chatRoom.Text = c.FormValue("text")
	chatRoom.Reply = c.FormValue("reply")
	sentAt := c.FormValue("sentAt")
	mention := c.FormValue("mention")

	var fileUploaded bool = true

	if chatRoom.UserID == "" || sentAt == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Missing userID/SentAt!")
	}

	parsedSentAt, err := time.Parse(time.RFC3339, sentAt)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid sentAt format! Must be an ISO 8601 string.")
	}
	fmt.Printf("parsedSentAt: %v\n", parsedSentAt)
	chatRoom.SentAt = parsedSentAt

	var mentions []string
	if mention != "" {
		err = json.Unmarshal([]byte(mention), &mentions)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, "Invalid mention format! Must be a JSON stringified array of strings.")
		}
	}
	fmt.Printf("Mentions: %v\n", mentions)

	fileHeader, err := c.FormFile("file")
	if err != nil {
		if err.Error() == constants.NO_FILE_UPLOADED_ERROR {
			fmt.Println("No file uploaded")
			// Prevent empty Text field when there is no file uploaded
			if chatRoom.Text == "" {
				return fiber.NewError(fiber.StatusBadRequest, "Missing Text field!")
			}
			fileUploaded = false

		} else {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	var filePath string
	var imageUrl string
	var chatRoomFile models.ChatroomFile
	var chatRoomMentions []models.ChatroomMention
	dimensions := struct {
		Height int `json:"height"`
		Width  int `json:"Width"`
	}{Height: 0, Width: 0}

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

		filePath = packages.GenFilePath(fileHeader.Filename)
		firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

		imageUrl, err = firebaseStorage.Add(file, fileHeader)
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	chatRoom.ArrivedAt = time.Now()

	newChatRoom, err := chatRoom.Create(chatRoom)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	// Save chatFile
	if imageUrl != "" {
		dimensionJson, err := json.Marshal(dimensions)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		chatRoomFile = models.ChatroomFile{ChatroomID: newChatRoom.ID,
			URL: imageUrl, Path: filePath, Dimensions: models.JSONB(dimensionJson)}
		newChatRoomFile, err := chatRoomFile.Create(chatRoomFile)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newChatRoom.File = newChatRoomFile
	}

	// Save all mentions of the chat message
	for _, mention := range mentions {
		if mention == "" {
			continue
		}
		chatroomMention := models.ChatroomMention{ChatroomID: newChatRoom.ID, UserID: mention}
		newChatroomMention, err := chatroomMention.Create(chatroomMention)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		chatRoomMentions = append(chatRoomMentions, newChatroomMention)
	}

	newChatRoom.Mention = chatRoomMentions

	events.EB.Publish("chatroomMessage", newChatRoom)
	events.EB.Publish("chatNotification", newChatRoom)

	response := fiber.Map{
		"status":  "success",
		"message": "chat posted successfully!",
		"data":    newChatRoom,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
