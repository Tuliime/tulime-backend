package chatroom

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/handlers/messenger"
	"github.com/Tuliime/tulime-backend/internal/middlewares"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/Tuliime/tulime-backend/internal/packages"
	"github.com/Tuliime/tulime-backend/internal/sse"
)

// TODO: To move all the sse operations to one single endpoint("api/v0.01/live")
func GetLiveChat(w http.ResponseWriter, r *http.Request) {
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

	chatroomMessageChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("chatroomMessage", chatroomMessageChan)

	type Messenger = messenger.Message
	messengerChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("messenger", messengerChan)

	type OnlineStatus = models.OnlineStatus
	onlineStatusChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("onlineStatus", onlineStatusChan)

	typingStatusChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("typingStatus", typingStatusChan)

	for {
		select {
		case chatroomMessageEvent := <-chatroomMessageChan:
			chatroomMessage, ok := chatroomMessageEvent.Data.(Message)
			if !ok {
				log.Printf("Invalid message type received: %T", chatroomMessageEvent.Data)
				return
			}
			if err := cm.SendEvent("chatroom-message", chatroomMessage, userID); err != nil {
				log.Printf("Error sending chatroom-message event: %v\n", err)
				return
			}
		case messengerEvent := <-messengerChan:
			messengerMsg, ok := messengerEvent.Data.(Messenger)
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
			if err := cm.SendEvent("online-status", onlineStatus, userID); err != nil {
				log.Printf("Error sending online-status event: %v\n", err)
				return
			}
		case typingStatusEvent := <-typingStatusChan:
			typingStatus, ok := typingStatusEvent.Data.(TypingStatus)
			if !ok {
				log.Printf("Invalid message type received: %T", typingStatusEvent.Data)
				return
			}
			if err := cm.SendEvent("typing-status", typingStatus, userID); err != nil {
				log.Printf("Error sending keep-alive event: %v\n", err)
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
