package subscribers

import "log"

func InitEventSubscribers() {
	log.Println("Initiating global event subscribers...")

	go ChatNotificationEventListener()
	go MessengerNotificationEventListener()
	go NotificationEventListener()
}

// func init() {
// 	InitEventSubscribers()
// }
