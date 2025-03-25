package messenger

import (
	"strconv"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

// TODO 1: To enhance the fetching algorithm to get
// messages starting from first unread
// TODO 2: To include the advert details of a give tag
var GetMessagesByRoom = func(c *fiber.Ctx) error {
	messenger := models.Messenger{}
	messengerRoom := models.MessengerRoom{}
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")
	inCludeCursorParam := c.Query("includeCursor", "false")
	direction := c.Query("direction")
	messengerRoomID := c.Query("messengerRoomID")
	userOneID := c.Query("userOneID")
	userTwoID := c.Query("userTwoID")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	includeCursor, err := strconv.ParseBool(inCludeCursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	if direction == "" {
		direction = "BACKWARD"
	}

	if messengerRoomID == "" {
		if userOneID == "" || userTwoID == "" {
			return fiber.NewError(fiber.StatusBadRequest, "Missing userOneID/userTwoID")
		}
		messengerRoom, err = messengerRoom.FindByUsers(userOneID, userTwoID)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		if messengerRoom.ID == "" {
			return fiber.NewError(fiber.StatusInternalServerError,
				"Provided userOneID and userTwoID don't have messengerRoomID yet!")
		}
		messengerRoomID = messengerRoom.ID
	}

	messengerMsgs, err := messenger.FindByRoom(messengerRoomID, limit+1,
		cursorParam, includeCursor, direction)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var repliedMessages []models.Messenger

	for _, messengerMsg := range messengerMsgs {
		if messengerMsg.Reply == "" {
			continue
		}
		reply, err := messenger.FindReply(messengerMsg.Reply)
		if err != nil {
			if err.Error() == constants.RECORD_NOT_FOUND_ERROR {
				continue
			}
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		repliedMessages = append(repliedMessages, reply)
	}

	var prevCursor, nextCursor string
	var hasPrevItems, hasNextItems bool

	if len(messengerMsgs) > 0 && direction == "BACKWARD" {
		if len(messengerMsgs) > int(limit) {
			messengerMsgs = messengerMsgs[1:] // Remove first element
			prevCursor = messengerMsgs[0].ID
			hasPrevItems = true
			if cursorParam != "" {
				nextCursor = messengerMsgs[len(messengerMsgs)-1].ID
				hasNextItems = true
			}
		} else {
			prevCursor = ""
			hasPrevItems = false
			if cursorParam != "" {
				nextCursor = messengerMsgs[len(messengerMsgs)-1].ID
				hasNextItems = true
			}
		}
	}

	if len(messengerMsgs) > 0 && direction == "FORWARD" {
		if len(messengerMsgs) > int(limit) {
			messengerMsgs = messengerMsgs[:len(messengerMsgs)-1] // Remove last element
			nextCursor = messengerMsgs[len(messengerMsgs)-1].ID
			hasNextItems = true
			if cursorParam != "" {
				prevCursor = messengerMsgs[0].ID
				hasPrevItems = true
			}
		} else {
			nextCursor = ""
			hasNextItems = false
			if cursorParam != "" {
				prevCursor = messengerMsgs[0].ID
				hasPrevItems = true
			}
		}
	}

	pagination := map[string]interface{}{
		"limit":         limit,
		"prevCursor":    prevCursor,
		"nextCursor":    nextCursor,
		"includeCursor": includeCursor,
		"hasNextItems":  hasNextItems,
		"hasPrevItems":  hasPrevItems,
		"direction":     direction,
	}

	messengerMap := fiber.Map{
		"messages": messengerMsgs,
		"replies":  repliedMessages,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       messengerMap,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}
