// package chatroom

// import (
// 	"bufio"
// 	"bytes"
// 	"context"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net"
// 	"strings"
// 	"time"

// 	"github.com/Tuliime/tulime-backend/internal/events"
// 	"github.com/Tuliime/tulime-backend/internal/models"
// 	"github.com/gofiber/fiber/v2"
// 	"github.com/valyala/fasthttp"
// )

// func formatSSEMessage(eventType string, data any) (string, error) {
// 	var buf bytes.Buffer
// 	enc := json.NewEncoder(&buf)

// 	m := map[string]any{
// 		"data": data,
// 		"type": eventType,
// 	}

// 	err := enc.Encode(m)
// 	if err != nil {
// 		return "", err
// 	}

// 	sb := strings.Builder{}
// 	sb.WriteString(fmt.Sprintf("data: %v\n\n", buf.String()))

// 	return sb.String(), nil
// }

// var GetLiveChat = func(c *fiber.Ctx) error {
// 	c.Set("Content-Type", "text/event-stream")
// 	c.Set("Cache-Control", "no-cache")
// 	c.Set("Connection", "keep-alive")
// 	c.Set("Transfer-Encoding", "chunked")

// 	keepAliveTicker := time.NewTicker(30 * time.Second)
// 	keepAliveMsg := "keepalive"

// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()

// 	type ChatRoom = models.Chatroom
// 	chatroomMessageChan := make(chan events.DataEvent, 100)
// 	events.EB.Subscribe("chatroomMessage", chatroomMessageChan)

// 	log.Println("Client connecting...")

// 	c.Context().Hijack(func(conn net.Conn) {
// 		defer conn.Close()
// 		log.Printf("Inside hijacker")

// 		buf := make([]byte, 1)
// 		for {
// 			// log.Printf("Inside hijacker loop...")
// 			if _, err := conn.Read(buf); err != nil {
// 				log.Printf("Connection monitoring detected closure: %v", err)
// 				keepAliveTicker.Stop()
// 				events.EB.Unsubscribe("chatroomMessage", chatroomMessageChan)
// 				cancel()
// 				log.Println("Connection closed and cleanup completed")
// 				return
// 			}
// 		}
// 	})

// 	streamStarted := make(chan struct{})

// 	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
// 		close(streamStarted)

// 		for {
// 			select {
// 			case <-ctx.Done():
// 				log.Println("Context cancelled, stopping stream")
// 				return

// 			case chatroomMessageEvent := <-chatroomMessageChan:
// 				chatroomMessage, ok := chatroomMessageEvent.Data.(ChatRoom)
// 				if !ok {
// 					log.Printf("Invalid message type received: %T", chatroomMessageEvent.Data)
// 					return
// 				}

// 				chatroomMsgStr, err := formatSSEMessage("chatroom-message", chatroomMessage)
// 				if err != nil {
// 					log.Printf("Error formatting SSE message: %v\n", err)
// 					return
// 				}

// 				if _, err := fmt.Fprintf(w, "%s", chatroomMsgStr); err != nil {
// 					log.Printf("Error writing chatroom message: %v\n", err)
// 					return
// 				}

// 				if err := w.Flush(); err != nil {
// 					log.Printf("Error flushing chatroom message: %v\n", err)
// 					return
// 				}
// 				log.Printf("Message sent: %v", chatroomMessage)

// 			case <-keepAliveTicker.C:
// 				keepAliveMsg, err := formatSSEMessage("keep-alive", keepAliveMsg)
// 				if err != nil {
// 					log.Printf("Error formatting keep-alive: %v", err)
// 					continue
// 				}

// 				if _, err := fmt.Fprintf(w, "%s", keepAliveMsg); err != nil {
// 					log.Printf("Error writing keep-alive message: %v\n", err)
// 					return
// 				}

// 				if err := w.Flush(); err != nil {
// 					log.Printf("Error flushing keep-alive message: %v\n", err)
// 					return
// 				}
// 				log.Println("Keep-alive sent")
// 			}
// 		}
// 	}))

// 	<-streamStarted
// 	return nil
// }

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
	defer keepAliveTicker.Stop() // Ensure ticker is stopped

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	type ChatRoom = models.Chatroom
	chatroomMessageChan := make(chan events.DataEvent, 100)
	events.EB.Subscribe("chatroomMessage", chatroomMessageChan)
	defer events.EB.Unsubscribe("chatroomMessage", chatroomMessageChan) // Ensure unsubscribe

	log.Println("Client connecting...")

	streamStarted := make(chan struct{}) // Channel to signal stream start

	// *** CRITICAL: SetBodyStreamWriter *BEFORE* Hijack ***
	c.Context().SetBodyStreamWriter(fasthttp.StreamWriter(func(w *bufio.Writer) {
		close(streamStarted) // Signal stream has started

		for {
			select {
			case <-ctx.Done():
				log.Println("Context cancelled, stopping stream")
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

	c.Context().Hijack(func(conn net.Conn) {
		defer conn.Close()
		log.Printf("Inside hijacker")

		buf := make([]byte, 1)
		for {
			if _, err := conn.Read(buf); err != nil {
				log.Printf("Connection monitoring detected closure: %v", err)
				cancel() // Signal SSE stream to stop
				log.Println("Connection closed and cleanup completed")
				return
			}
		}
	})

	<-streamStarted // Wait for stream to initialize
	return nil      // Important: Return nil after Hijack
}
