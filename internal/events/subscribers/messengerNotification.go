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

		if isReply {
			notificationBody = fmt.Sprintf("%s replied to your message.", user.Name)
		} else if isTag {
			notificationBody = fmt.Sprintf("%s tagged a product in a new message.", user.Name)
		} else {
			notificationBody = fmt.Sprintf("%s sent you a new message.", user.Name)
		}

		jsonNotificationData, err := json.Marshal(struct {
			MessengerID string `json:"messengerID"`
			FileURL     string `json:"fileURL"`
		}{
			MessengerID: messengerMsg.ID,
			FileURL:     messenger.File.URL,
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
