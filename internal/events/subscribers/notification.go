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
