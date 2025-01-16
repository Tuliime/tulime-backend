package chatroom

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var PostChat = func(c *fiber.Ctx) error {
	chatRoom := models.Chatroom{}

	chatRoom.UserID = c.FormValue("userID")
	chatRoom.Text = c.FormValue("text")
	chatRoom.Reply = c.FormValue("reply")
	sentAt := c.FormValue("sentAt")
	mention := c.FormValue("mention")

	var fileUploaded bool = true

	// TODO: text to be validated alongside file
	if chatRoom.UserID == "" || chatRoom.Text == "" || sentAt == "" {
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

	file, err := c.FormFile("file")
	if err != nil {
		if err.Error() == constants.NO_FILE_UPLOADED_ERROR {
			fmt.Println("Error:", err.Error())
			fileUploaded = false

		} else {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
	}

	var filePath string
	var imageUrl string
	var chatRoomFile models.ChatroomFile
	var chatRoomMentions []models.ChatroomMention

	if fileUploaded {
		// Validate file size (10 MB limit)
		const maxFileSize = 10 << 20 // 10 MB in bytes
		if file.Size > maxFileSize {
			return fiber.NewError(fiber.StatusBadRequest, "File size exceeds the 10 MB limit")
		}
		fileReader, err := file.Open()
		if err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		defer fileReader.Close()

		filePath = packages.GenFilePath(file.Filename)
		firebaseStorage := packages.FirebaseStorage{FilePath: filePath}

		imageUrl, err = firebaseStorage.Add(fileReader, file)
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
		chatRoomFile = models.ChatroomFile{ChatroomID: newChatRoom.ID, URL: imageUrl, Path: filePath}
		newChatRoomFile, err := chatRoomFile.Create(chatRoomFile)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		newChatRoom.File = newChatRoomFile
	}

	// Save all mentions of the chat message
	for _, mention := range mentions {
		chatroomMention := models.ChatroomMention{ChatroomID: newChatRoom.ID, UserID: mention}
		newChatroomMention, err := chatroomMention.Create(chatroomMention)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		chatRoomMentions = append(chatRoomMentions, newChatroomMention)
	}

	newChatRoom.Mention = chatRoomMentions

	events.EB.Publish("chatroomMessage", newChatRoom)

	response := fiber.Map{
		"status":  "success",
		"message": "chat posted successfully!",
		"data":    newChatRoom,
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}
