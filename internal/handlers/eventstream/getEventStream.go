package eventstream

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/handlers/chatroom"
	"github.com/Tuliime/tulime-backend/internal/handlers/messenger"
	"github.com/Tuliime/tulime-backend/internal/middlewares"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/Tuliime/tulime-backend/internal/sse"
)

func GetEventStream(w http.ResponseWriter, r *http.Request) {
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

	type ChatroomMessage = chatroom.Message
	chatroomMessageChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("chatroomMessage", chatroomMessageChan)

	type MessengerMessage = messenger.Message
	messengerChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("messenger", messengerChan)

	type OnlineStatus = models.OnlineStatus
	onlineStatusChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("onlineStatus", onlineStatusChan)

	type TypingStatus = chatroom.TypingStatus
	typingStatusChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("typingStatus", typingStatusChan)

	type Notification = models.Notification
	notificationChan := make(chan events.DataEvent)
	events.EB.Subscribe("liveNotification", notificationChan)

	for {
		select {
		case chatroomMessageEvent := <-chatroomMessageChan:
			chatroomMessage, ok := chatroomMessageEvent.Data.(ChatroomMessage)
			if !ok {
				log.Printf("Invalid message type received: %T", chatroomMessageEvent.Data)
				return
			}
			if err := cm.SendEvent("chatroom-message", chatroomMessage, userID); err != nil {
				log.Printf("Error sending chatroom-message event: %v\n", err)
				return
			}
		case messengerEvent := <-messengerChan:
			messengerMsg, ok := messengerEvent.Data.(MessengerMessage)
			if !ok {
				log.Printf("Invalid messenger type received: %T", messengerEvent.Data)
				return
			}
			if err := cm.SendEvent("messenger", messengerMsg, messengerMsg.RecipientID); err != nil {
				log.Printf("Error sending messenger msg event: %v\n", err)
				return
			}
		case onlineStatusEvent := <-onlineStatusChan:
			onlineStatus, ok := onlineStatusEvent.Data.(OnlineStatus)
			if !ok {
				log.Printf("Invalid message type received: %T", onlineStatusEvent.Data)
				return
			}
			log.Printf("sending online-status event : %+v", onlineStatus)

			if err := cm.SendEvent("online-status", onlineStatus, userID); err != nil {
				log.Printf("Error sending online-status event: %v\n", err)
				return
			}
		case typingStatusEvent := <-typingStatusChan:
			typingStatus, ok := typingStatusEvent.Data.(TypingStatus)
			if !ok {
				log.Printf("Invalid typing-status received: %T", typingStatusEvent.Data)
				return
			}
			if typingStatus.Type == "messenger" {
				if err := cm.SendEvent("typing-status-messenger", typingStatus,
					typingStatus.RecipientID); err != nil {
					log.Printf("Error sending typing-status event: %v\n", err)
					return
				}
			}
			if typingStatus.Type == "chatroom" {
				if err := cm.SendEvent("typing-status-chatroom", typingStatus, userID); err != nil {
					log.Printf("Error sending typing-status event: %v\n", err)
					return
				}
			}
		case notificationEvent := <-notificationChan:
			notification, ok := notificationEvent.Data.(Notification)
			if !ok {
				log.Printf("Invalid notification type received: %T", notificationEvent.Data)
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
			log.Printf("All Clients :%+v", cm.GetAllClients())
			cm.RemoveClient(userID)
			keepAliveTicker.Stop()
			cancel()
			log.Println("Client disconnected:", userID)
			return
		}
	}
}
