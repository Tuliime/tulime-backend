package chatroom

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"strings"
	"time"

	"github.com/Tuliime/tulime-backend/internal/events"
	"github.com/Tuliime/tulime-backend/internal/models"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
)

func formatSSEMessage(eventType string, data any) (string, error) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)

	m := map[string]any{
		"data": data,
		"type": eventType,
	}

	err := enc.Encode(m)
	if err != nil {
		return "", err
	}

	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("data: %v\n\n", buf.String()))

	return sb.String(), nil
}

var GetLiveChat = func(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	keepAliveTicker := time.NewTicker(30 * time.Second)
	keepAliveMsg := "keepalive"

	ctx, cancel := context.WithCancel(context.Background())

	type ChatRoom = models.Chatroom
	chatroomMessageChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("chatroomMessage", chatroomMessageChan)
	log.Printf("Client '%s' connecting...", c.Locals("userID"))

	type OnlineStatus = models.OnlineStatus
	onlineStatusChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("onlineStatus", onlineStatusChan)

	typingStatusChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("typingStatus", typingStatusChan)

	c.Context().HijackSetNoResponse(false)
	c.Context().Hijack(func(conn net.Conn) {
		defer conn.Close()

		buf := make([]byte, 1)
		for {
			if _, err := conn.Read(buf); err != nil {
				log.Printf("Connection monitoring detected closure: %v", err)
				cancel()
				log.Println("Connection closed and cleanup completed")
				return
			}
		}
	})
	log.Println("Hijacked: ", c.Context().Hijacked())

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping stream")
				keepAliveTicker.Stop()
				events.EB.Unsubscribe("chatroomMessage", chatroomMessageChan)
				return

			case chatroomMessageEvent := <-chatroomMessageChan:
				chatroomMessage, ok := chatroomMessageEvent.Data.(ChatRoom)
				if !ok {
					log.Printf("Invalid message type received: %T", chatroomMessageEvent.Data)
					return
				}

				chatroomMsgStr, err := formatSSEMessage("chatroom-message", chatroomMessage)
				if err != nil {
					log.Printf("Error formatting SSE message: %v\n", err)
					return
				}

				if _, err := fmt.Fprintf(w, "%s", chatroomMsgStr); err != nil {
					log.Printf("Error writing chatroom message: %v\n", err)
					return
				}

				if err := w.Flush(); err != nil {
					log.Printf("Error flushing chatroom message: %v\n", err)
					return
				}
				log.Printf("Message sent: %v", chatroomMessage)

			case onlineStatusEvent := <-onlineStatusChan:
				onlineStatus, ok := onlineStatusEvent.Data.(OnlineStatus)
				if !ok {
					log.Printf("Invalid message type received: %T", onlineStatusEvent.Data)
					return
				}

				onlineStatusStr, err := formatSSEMessage("online-status", onlineStatus)
				if err != nil {
					log.Printf("Error formatting SSE message: %v\n", err)
					return
				}

				if _, err := fmt.Fprintf(w, "%s", onlineStatusStr); err != nil {
					log.Printf("Error writing online status: %v\n", err)
					return
				}

				if err := w.Flush(); err != nil {
					log.Printf("Error flushing online status: %v\n", err)
					return
				}
				log.Printf("Online status sent: %v", onlineStatus)

			case typingStatusEvent := <-typingStatusChan:
				typingStatus, ok := typingStatusEvent.Data.(TypingStatus)
				if !ok {
					log.Printf("Invalid message type received: %T", typingStatusEvent.Data)
					return
				}

				typingStatusStr, err := formatSSEMessage("typing-status", typingStatus)
				if err != nil {
					log.Printf("Error formatting SSE message: %v\n", err)
					return
				}

				if _, err := fmt.Fprintf(w, "%s", typingStatusStr); err != nil {
					log.Printf("Error writing typing status: %v\n", err)
					return
				}

				if err := w.Flush(); err != nil {
					log.Printf("Error flushing typing status: %v\n", err)
					return
				}
				log.Printf("typing status sent: %v", typingStatus)

			case <-keepAliveTicker.C:
				keepAliveMsg, err := formatSSEMessage("keep-alive", keepAliveMsg)
				if err != nil {
					log.Printf("Error formatting keep-alive: %v", err)
					continue
				}

				if _, err := fmt.Fprintf(w, "%s", keepAliveMsg); err != nil {
					log.Printf("Error writing keep-alive message: %v\n", err)
					return
				}

				if err := w.Flush(); err != nil {
					log.Printf("Error flushing keep-alive message: %v\n", err)
					return
				}
				log.Println("Keep-alive sent")
			}
		}
	}))
	return nil
}
