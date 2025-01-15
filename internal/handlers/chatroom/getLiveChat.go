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
	}

	err := enc.Encode(m)
	if err != nil {
		return "", nil
	}
	sb := strings.Builder{}

	sb.WriteString(fmt.Sprintf("event: %s\n", eventType))
	sb.WriteString(fmt.Sprintf("retry: %d\n", 15000))
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
	warmUpMsg := "warmup"

	type ChatRoom = models.Chatroom
	chatroomChan := make(chan events.DataEvent)
	events.EB.Subscribe("chatroomMessage", chatroomChan)

	ctx, cancel := context.WithCancel(c.Context())
	disconnect := ctx.Done()
	// disconnect := c.Context().Done()

	go func() {
		<-disconnect
		keepAliveTicker.Stop()
		events.EB.Unsubscribe("chatroomMessage", chatroomChan)
		log.Println("Client disconnected")
		cancel()
	}()

	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {

		if _, err := fmt.Fprintf(w, "%s", warmUpMsg); err != nil {
			log.Printf("Error while writing Data: %v\n", err)
			return
		}

		if err := w.Flush(); err != nil {
			log.Printf("Error while flushing Data: %v\n", err)
			keepAliveTicker.Stop()
			return
		}
		for {
			select {
			case chatroomEvent := <-chatroomChan:
				chatroomMessage, ok := chatroomEvent.Data.(ChatRoom)
				if !ok {
					log.Println("Interface does not hold type ChatRoom")
					return
				}
				chatroomMsgStr, err := formatSSEMessage("current-value", chatroomMessage)
				if err != nil {
					log.Printf("Error formatting sse message: %v\n", err)
					return
				}

				if _, err := fmt.Fprintf(w, "%s", chatroomMsgStr); err != nil {
					log.Printf("Error while writing Data: %v\n", err)
					return
				}

				if err := w.Flush(); err != nil {
					log.Printf("Error while flushing Data: %v\n", err)
					keepAliveTicker.Stop()
					return
				}
			case <-keepAliveTicker.C:
				log.Println("keep alive message sent")
				keepAliveMsgStr, err := formatSSEMessage("current-value", keepAliveMsg)
				if err != nil {
					log.Printf("Error formatting sse message: %v\n", err)
					return
				}
				if _, err = fmt.Fprintf(w, "%s", keepAliveMsgStr); err != nil {
					log.Printf("Error while writing Data: %v\n", err)
					return
				}

				if err = w.Flush(); err != nil {
					log.Printf("Error while flushing Data: %v\n", err)
					keepAliveTicker.Stop()
					return
				}

				// case <-disconnect:
				// 	keepAliveTicker.Stop()
				// 	events.EB.Unsubscribe("chatroomMessage", chatroomChan)
				// 	log.Println("Client disconnected")
				// 	cancel()
				// 	return
			}
		}

	}))

	return nil
}
