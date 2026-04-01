// Package websocket provides real-time communication capabilities using WebSockets,
// including chat, notifications, and user status tracking.
package websocket

import (
	"bytes"
	"context"
	"encoding/json"
	"forum/internal/testutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

//--------------------------------------------------------------------------------------|

// TestServeWs verifies the full WebSocket lifecycle: connecting, registering with the Hub,
// sending messages, and receiving broadcasts.
func TestServeWs(t *testing.T) {
	hub := NewHub(&testutil.MockChatRepo{}, &testutil.MockNotificationsRepo{}, nil)
	go hub.Run(context.Background())

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ServeWs(hub, w, r, testutil.TestUserID)
	}))
	defer server.Close()

	url := "ws" + strings.TrimPrefix(server.URL, "http")

	t.Run("Full loop: join and message", func(t *testing.T) {
		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			t.Fatalf("Failed to dial: %v", err)
		}
		defer conn.Close()

		// Wait for registration with a timeout-aware check
		testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
			online := hub.GetOnlineUsers()
			for _, id := range online {
				if id == testutil.TestUserID {
					return true
				}
			}
			return false
		})

		// Send a message from client (private message to self)
		msg := map[string]interface{}{
			"type": "private_message",
			"data": map[string]interface{}{
				"receiver_id": testutil.TestUserID,
				"body":        "hello",
			},
		}
		b, _ := json.Marshal(msg)
		if err := conn.WriteMessage(websocket.TextMessage, b); err != nil {
			t.Fatalf("Failed to write: %v", err)
		}

		// We might receive the status update message first.
		// We loop until we find our private message.
		deadline := time.Now().Add(testutil.LongTestTimeout)
		for time.Now().Before(deadline) {
			_, p, err := conn.ReadMessage()
			if err != nil {
				t.Fatalf("Failed to read: %v", err)
			}

			// Use a decoder as we might have multiple JSON objects in one frame
			dec := json.NewDecoder(bytes.NewReader(p))
			for dec.More() {
				var received struct {
					Type string `json:"type"`
				}
				if err := dec.Decode(&received); err != nil {
					t.Fatalf("Failed to decode: %v", err)
				}

				if received.Type == "private_message" {
					return // Success
				}
			}
			// Otherwise keep waiting for next frame
		}
		t.Fatal("Timed out waiting for private_message")
	})
}
