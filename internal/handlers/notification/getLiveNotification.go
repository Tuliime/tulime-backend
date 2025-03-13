package notification

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

// TODO: To send notification to the individual users
var GetLiveNotification = func(c *fiber.Ctx) error {
	c.Set("Content-Type", "text/event-stream")
	c.Set("Cache-Control", "no-cache")
	c.Set("Connection", "keep-alive")
	c.Set("Transfer-Encoding", "chunked")

	keepAliveTicker := time.NewTicker(30 * time.Second)
	keepAliveMsg := "keepalive"

	ctx, cancel := context.WithCancel(context.Background())

	type Notification = models.Notification
	notificationChan := make(chan events.DataEvent)
	events.EB.Subscribe("liveNotification", notificationChan)
	log.Printf("Client '%s' connecting...", c.Locals("userID"))

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
				events.EB.Unsubscribe("chatroomMessage", notificationChan)
				return

			case notificationEvent := <-notificationChan:
				notification, ok := notificationEvent.Data.(Notification)
				if !ok {
					log.Printf("Invalid message type received: %T", notificationEvent.Data)
					return
				}

				notificationStr, err := formatSSEMessage("notification", notification)
				if err != nil {
					log.Printf("Error formatting SSE message: %v\n", err)
					return
				}

				if _, err := fmt.Fprintf(w, "%s", notificationStr); err != nil {
					log.Printf("Error writing chatroom message: %v\n", err)
					return
				}

				if err := w.Flush(); err != nil {
					log.Printf("Error flushing chatroom message: %v\n", err)
					return
				}
				log.Printf("Message sent: %v", notification)

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
