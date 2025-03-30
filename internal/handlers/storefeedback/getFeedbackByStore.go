package storefeedback

import (
	"time"

	"github.com/Tuliime/tulime-backend/internal/constants"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/gofiber/fiber/v2"
)

var GetFeedbackByStore = func(c *fiber.Ctx) error {
	feedback := models.StoreFeedback{}
	storeID := c.Params("id")
	limitParam := c.Query("limit")
	cursorParam := c.Query("cursor")

	limit, err := packages.ValidateQueryLimit(limitParam)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	savedFeedback, err := feedback.FindByStore(storeID, limit+1, cursorParam)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}

	var feedbackResponse []FeedbackResponse

	for _, feedback := range savedFeedback {
		var replyFeedback []models.StoreFeedback

		if feedback.Reply != "" {
			replyFeedback, err = feedback.FindReply(feedback.Reply)
			if err != nil && err.Error() != constants.RECORD_NOT_FOUND_ERROR {
				return fiber.NewError(fiber.StatusInternalServerError, err.Error())
			}
		}

		feedbackResponse = append(feedbackResponse, FeedbackResponse{
			ID:            feedback.ID,
			StoreID:       feedback.StoreID,
			UserID:        feedback.UserID,
			Title:         feedback.Title,
			Description:   feedback.Description,
			Reply:         feedback.Reply,
			ReplyFeedback: replyFeedback,
			File:          feedback.File,
			CreatedAt:     feedback.CreatedAt,
			UpdatedAt:     feedback.UpdatedAt,
			Store:         feedback.Store,
			User: User{
				ID:             feedback.User.ID,
				Name:           feedback.User.Name,
				Role:           feedback.User.Role,
				TelNumber:      feedback.User.TelNumber,
				ImageUrl:       feedback.User.ImageUrl,
				ImagePath:      feedback.User.ImagePath,
				ProfileBgColor: feedback.User.ProfileBgColor,
				ChatroomColor:  feedback.User.ChatroomColor,
				CreatedAt:      feedback.User.CreatedAt,
				UpdatedAt:      feedback.User.UpdatedAt,
			},
		})
	}

	var nextCursor string
	var hasNextItems bool

	if len(savedFeedback) > 0 && len(savedFeedback) > int(limit) {
		savedFeedback = savedFeedback[:len(savedFeedback)-1] // Remove last element
		nextCursor = savedFeedback[0].ID
		hasNextItems = true
	} else {
		nextCursor = ""
		hasNextItems = false
	}

	pagination := fiber.Map{
		"limit":        limit,
		"nextCursor":   nextCursor,
		"hasNextItems": hasNextItems,
	}

	response := fiber.Map{
		"status":     "success",
		"data":       feedbackResponse,
		"pagination": pagination,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

type User struct {
	ID             string    `json:"id"`
	Name           string    `json:"name"`
	TelNumber      int       `json:"telNumber"`
	Role           string    `json:"role"`
	ImageUrl       string    `json:"imageUrl"`
	ImagePath      string    `json:"imagePath"`
	ProfileBgColor string    `json:"profileBgColor"`
	ChatroomColor  string    `json:"chatroomColor"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

type FeedbackResponse struct {
	ID            string                     `json:"id"`
	StoreID       string                     `json:"storeID"`
	UserID        string                     `json:"userID"`
	Title         string                     `json:"title"`
	Description   string                     `json:"description"`
	Reply         string                     `json:"reply"`
	ReplyFeedback []models.StoreFeedback     `json:"replyFeedback"`
	File          []models.StoreFeedbackFile `json:"files"`
	CreatedAt     time.Time                  `json:"createdAt"`
	UpdatedAt     time.Time                  `json:"updatedAt"`
	Store         *models.Store              `json:"store"`
	User          User                       `json:"user"`
}
