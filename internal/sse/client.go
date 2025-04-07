package sse

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
)

type ClientManager struct {
	clients map[string]http.ResponseWriter
	sync.RWMutex
}

func NewClientManager() *ClientManager {
	return &ClientManager{
		clients: make(map[string]http.ResponseWriter),
	}
}

func (cm *ClientManager) AddClient(userId string, w http.ResponseWriter) {
	cm.Lock()
	defer cm.Unlock()
	cm.clients[userId] = w
}

func (cm *ClientManager) RemoveClient(userId string) {
	cm.Lock()
	defer cm.Unlock()
	delete(cm.clients, userId)
}

func (cm *ClientManager) GetClient(userId string) (http.ResponseWriter, bool) {
	cm.RLock()
	defer cm.RUnlock()
	client, ok := cm.clients[userId]
	return client, ok
}

func (cm *ClientManager) GetAllClients() map[string]http.ResponseWriter {
	cm.RLock()
	defer cm.RUnlock()

	clientsCopy := make(map[string]http.ResponseWriter)
	for userID, writer := range cm.clients {
		clientsCopy[userID] = writer
	}
	return clientsCopy
}

func (cm *ClientManager) SendEvent(eventType string, data any, userID string) error {
	w, ok := cm.GetClient(userID)
	if !ok {
		return nil
	}

	f, ok := w.(http.Flusher)
	if !ok {
		log.Println("Response writer does not implement http.Flusher")
		return nil
	}

	sseDataStr, err := cm.formatEventData(eventType, data)
	if err != nil {
		log.Printf("Error formatting SSE message: %v\n", err)
		return nil
	}

	if _, err = w.Write([]byte(sseDataStr)); err != nil {
		log.Println("Error writing to response writer:", err)
		return err
	}
	f.Flush()
	return nil
}

func (cm *ClientManager) formatEventData(eventType string, data any) (string, error) {
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
