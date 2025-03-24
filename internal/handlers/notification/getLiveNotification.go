package notification

import (
	"context"
	"log"

	"net/http"
	"time"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/middlewares"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/Tuliime/tulime-backend/internal/sse"
)

func GetLiveNotification(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	keepAliveTicker := time.NewTicker(30 * time.Second)

	userID, ok := r.Context().Value(middlewares.UserIDKey).(string)
	if !ok {
		packages.AppError("UserID not found in context", 500, w)
		return
	}
	log.Println("Client connected:", userID)

	cm := sse.NewClientManager()
	cm.AddClient(userID, w)

	if err := cm.SendEvent("keep-alive", "keepalive", userID); err != nil {
		return
	}

	ctx, cancel := context.WithCancel(r.Context())
	disconnect := ctx.Done()
	defer cancel()

	type Notification = models.Notification
	notificationChan := make(chan events.DataEvent)
	events.EB.Subscribe("liveNotification", notificationChan)

	for {
		select {
		case notificationEvent := <-notificationChan:
			notification, ok := notificationEvent.Data.(Notification)
			if !ok {
				log.Printf("Invalid message type received: %T", notificationEvent.Data)
				return
			}
			if err := cm.SendEvent("notification", notification, userID); err != nil {
				log.Printf("Error sending notification event: %v\n", err)
				return
			}
		case <-keepAliveTicker.C:
			if err := cm.SendEvent("keep-alive", "keepalive", userID); err != nil {
				log.Printf("Error sending keep-alive event: %v\n", err)
				return
			}
		case <-disconnect:
			cm.RemoveClient(userID)
			keepAliveTicker.Stop()
			cancel()
			log.Println("Client disconnected:", userID)
			return
		}
	}
}
