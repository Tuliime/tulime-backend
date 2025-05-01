package subscribers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
)

func processMessengerNotifications(messengerMsg models.Messenger) {
	user := models.User{}
	device := models.Device{}
	messenger := models.Messenger{}

	devices, err := device.FindByUser(messengerMsg.RecipientID)
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		return
	}

	user, err = user.FindOne(messengerMsg.SenderID)
	if err != nil {
		log.Printf("Error fetching user info %v: ", err)
	}

	var wg sync.WaitGroup
	for _, device := range devices {
		if device.NotificationDisabled {
			continue
		}

		var notificationBody string
		var isTag bool
		var isReply bool = messengerMsg.Reply != ""
		if len(messengerMsg.Tag) > 0 {
			isTag = messengerMsg.Tag[0].ID != ""
		}

		var messageWithPhoto bool = messengerMsg.Text != "" && messengerMsg.File.ID != ""
		var photoWithoutMessage bool = messengerMsg.File.ID != "" && messengerMsg.Text == ""

		if messageWithPhoto {
			notificationBody = fmt.Sprintf("📷 Photo - %s", messengerMsg.Text)
		} else if photoWithoutMessage {
			notificationBody = "📷 Photo"
		} else {
			notificationBody = messengerMsg.Text
		}

		if isReply {
			notificationBody = fmt.Sprintf("↩️ Replied:  %s", notificationBody)
		} else if isTag {
			notificationBody = fmt.Sprintf("🏷️ Tagged:  %s", notificationBody)
		} else {
			notificationBody = fmt.Sprintf("🆕 New:  %s", notificationBody)
		}

		jsonNotificationData, err := json.Marshal(struct {
			MessengerID string `json:"messengerID"`
			FileURL     string `json:"fileURL"`
			Type        string `json:"type"`
			ClientPath  string `json:"clientPath"`
		}{
			MessengerID: messengerMsg.ID,
			FileURL:     messenger.File.URL,
			Type:        "messenger",
			ClientPath:  fmt.Sprintf("/ecommerce/messenger/%s", messengerMsg.ID),
		})

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		notificationData := string(jsonNotificationData)

		notification := models.Notification{UserID: device.UserID,
			Title: user.Name, Body: notificationBody, Data: notificationData,
			Type: "messenger"}

		sendNotification := models.SendNotification{Notification: notification,
			DeviceToken: device.Token}

		wg.Add(1)
		go func(sendNotification models.SendNotification) {
			defer wg.Done()
			events.EB.Publish("sendNotification", sendNotification)
		}(sendNotification)
	}
	wg.Wait()

}

func MessengerNotificationEventListener() {
	type Messenger = models.Messenger
	messengerNotificationChan := make(chan events.DataEvent)
	events.EB.Subscribe("messengerNotification", messengerNotificationChan)

	for {
		messengerNotificationEvent := <-messengerNotificationChan
		messengerMsg, ok := messengerNotificationEvent.Data.(Messenger)
		if !ok {
			log.Printf("Invalid messenger msg type received: %T", messengerNotificationEvent.Data)
			return
		}
		go processMessengerNotifications(messengerMsg)
	}
}
