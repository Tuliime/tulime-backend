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
	To       string `json:"to,omitempty"`
	Title    string `json:"title,omitempty"`
	Body     string `json:"body,omitempty"`
	Priority string `json:"priority,omitempty"`
	// Data     map[string]interface{} `json:"data,omitempty"`
	Data string `json:"data,omitempty"`
}

func SendExpoNotification(sendNotification models.SendNotification) error {
	notification := models.Notification{}
	expoNotification := ExpoNotification{
		To:       sendNotification.DeviceToken,
		Title:    sendNotification.Notification.Title,
		Body:     sendNotification.Notification.Body,
		Priority: "high",
		Data:     sendNotification.Notification.Data,
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

// ExponentPushToken[MDeazfOMaoHvBMqbrwItvf]  //Test push notification

// curl -X POST https://exp.host/--/api/v2/push/send \
//      -H "Content-Type: application/json" \
//      -d '{
//            "to": "ExponentPushToken[MDeazfOMaoHvBMqbrwItvf]",
//            "title": "Test Expo Dev",
//            "body": "Hello from Expo, and it'\''s Dankan sending it",
//            "data": {
//              "url": "tulimeapp://notification",
//              "extraInfo": "Notification Screen"
//            },
//            "sound": "default",
//            "priority": "high",
//            "icon": "https://example.com/icon.png",
//            "attachments": {
//              "image": "https://firebasestorage.googleapis.com/v0/b/reserve-now-677ca.appspot.com/o/tulime%2Frice.png?alt=media&token=d9fb8814-0d9a-40e0-8bca-ebec0782fe4a"
//            }
//          }'
