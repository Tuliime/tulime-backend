package subscribers

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
)

func processChatNotifications(chatroomMessage models.Chatroom) {
	user := models.User{}
	device := models.Device{}
	chatroom := models.Chatroom{}

	devices, err := device.FindAll()
	if err != nil {
		log.Printf("Error fetching devices: %v", err)
		return
	}

	var repliedMessage models.Chatroom
	if chatroomMessage.Reply != "" {
		repliedMessage, err = chatroom.FindReply(chatroomMessage.Reply)
		if err != nil {
			log.Printf("Error fetching replied message %v: ", err)
		}

	}

	user, err = user.FindOne(chatroomMessage.UserID)
	if err != nil {
		log.Printf("Error fetching user info %v: ", err)
	}

	var wg sync.WaitGroup
	for _, device := range devices {
		// Prevent sending notification to the message sender
		if chatroomMessage.UserID == device.UserID {
			continue
		}
		if device.NotificationDisabled {
			continue
		}

		var notificationBody string
		var isReply, isMention bool

		isReply = repliedMessage.UserID == device.UserID

		for _, mention := range chatroomMessage.Mention {
			if mention.UserID == device.UserID {
				isMention = true
				break
			}
		}

		if isReply {
			notificationBody = fmt.Sprintf("%s replied to your message.", user.Name)
		} else if isMention {
			notificationBody = fmt.Sprintf("%s mentioned you in a message.", user.Name)
		} else {
			notificationBody = fmt.Sprintf("%s posted a new message.", user.Name)
		}

		jsonNotificationData, err := json.Marshal(struct {
			ChatroomID string `json:"chatroomID"`
			FileURL    string `json:"fileURL"`
		}{
			ChatroomID: chatroomMessage.ID,
			FileURL:    chatroom.File.URL,
		})

		if err != nil {
			fmt.Println("Error:", err)
			return
		}

		notificationData := string(jsonNotificationData)

		notification := models.Notification{UserID: device.UserID,
			Title: "Tulime ChatFarm", Body: notificationBody, Data: notificationData,
			Type: "chat"}

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

func ChatNotificationEventListener() {
	type Chatroom = models.Chatroom
	chatNotificationChan := make(chan events.DataEvent)
	events.EB.Subscribe("chatNotification", chatNotificationChan)

	for {
		chatNotificationEvent := <-chatNotificationChan
		chatroomMessage, ok := chatNotificationEvent.Data.(Chatroom)
		if !ok {
			log.Printf("Invalid Chatroom msg type received: %T", chatNotificationEvent.Data)
			return
		}
		go processChatNotifications(chatroomMessage)
	}
}
