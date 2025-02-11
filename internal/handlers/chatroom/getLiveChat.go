package chatroom

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
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

	type ChatRoom = models.Chatroom
	chatroomMessageChan := make(chan events.DataEvent)
	events.EB.Subscribe("chatroomMessage", chatroomMessageChan)

	ctx, cancel := context.WithCancel(c.Context())
	// disconnect := c.Context().Done()
	disconnect := ctx.Done()

	go func() {
		<-disconnect
		keepAliveTicker.Stop()
		events.EB.Unsubscribe("chatroomMessage", chatroomMessageChan)
		log.Println("Client disconnected")
		cancel()
	}()

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		go func() {
			for range keepAliveTicker.C {
				keepAliveMsgStr, err := formatSSEMessage("keep-alive", keepAliveMsg)
				if err != nil {
					log.Printf("Error formatting keep-alive message: %v\n", err)
					return
				}

				if _, err = fmt.Fprintf(w, "%s", keepAliveMsgStr); err != nil {
					log.Printf("Error writing keep-alive message: %v\n", err)
					return
				}

				if err = w.Flush(); err != nil {
					log.Printf("Error flushing keep-alive message: %v\n", err)
					return
				}
			}
		}()

		for chatroomMessageEvent := range chatroomMessageChan {
			chatroomMessage, ok := chatroomMessageEvent.Data.(ChatRoom)
			if !ok {
				log.Println("Invalid data type for ChatRoom")
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
		}
	}))

	return nil
}
