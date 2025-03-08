package notifications

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
)

type ExpoNotification struct {
	To    string `json:"to"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func SendExpoNotification(sendNotification models.SendNotification) error {
	notification := models.Notification{}
	expoNotification := ExpoNotification{
		To:    sendNotification.DeviceToken,
		Title: sendNotification.Notification.Title,
		Body:  sendNotification.Notification.Body,
	}

	jsonData, err := json.Marshal(expoNotification)
	if err != nil {
		return err
	}

	resp, err := http.Post("https://exp.host/--/api/v2/push/send", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	log.Printf("resp.StatusCode: %v", resp.StatusCode)

	sendNotification.Notification.SendStatusCode = resp.StatusCode

	newNotification, err := notification.Create(sendNotification.Notification)
	if err != nil {
		return err
	}

	events.EB.Publish("liveNotification", newNotification)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to send notification: %s", resp.Status)
	}
	return nil
}
