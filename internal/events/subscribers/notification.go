package subscribers

import (
	"log"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/handlers/notifications"
	"github.com/Tuliime/tulime-backend/internal/models"
)

func NotificationEventListener() {
	type SendNotification = models.SendNotification
	sendNotificationChan := make(chan events.DataEvent)
	events.EB.Subscribe("sendNotification", sendNotificationChan)

	for {
		notificationEvent := <-sendNotificationChan
		notification, ok := notificationEvent.Data.(SendNotification)
		if !ok {
			log.Printf("Invalid message type received: %T", notificationEvent.Data)
			return
		}
		log.Printf("sending notification...")
		go notifications.SendExpoNotification(notification)
	}
}
