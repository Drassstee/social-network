// Package websocket provides real-time communication capabilities using WebSockets,
// including chat, notifications, and user status tracking.
package websocket

import (
	"context"
	"encoding/json"
	"forum/internal/models"
	"forum/internal/testutil"
	"testing"
	"time"
)

//--------------------------------------------------------------------------------------|

// TestHub_Run verifies the core Hub lifecycle, managing client connections (register/unregister)
// and broadcasting messages properly.
func TestHub_Run(t *testing.T) {
	hub := NewHub(&testutil.MockChatRepo{}, &testutil.MockNotificationsRepo{}, nil)
	go hub.Run(context.Background())

	t.Run("Register", func(t *testing.T) {
		client := &Client{UserID: testutil.TestUserID, send: make(chan []byte, 1)}
		hub.register <- client

		testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
			online := hub.GetOnlineUsers()
			return len(online) == 1 && online[0] == testutil.TestUserID
		})
	})

	t.Run("Private Message", func(t *testing.T) {
		sender := &Client{UserID: testutil.TestUserID, send: make(chan []byte, 10)}
		receiver := &Client{UserID: testutil.TestUserID2, send: make(chan []byte, 10)}

		hub.register <- sender
		hub.register <- receiver

		// Wait for registrations to settle
		testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
			return len(hub.GetOnlineUsers()) == 2
		})

		// Drain status updates
		testutil.DrainChannel(sender.send)
		testutil.DrainChannel(receiver.send)

		msgData := map[string]interface{}{
			"receiver_id": testutil.TestUserID2,
			"body":        "hello",
		}
		raw, _ := json.Marshal(msgData)

		wsMsg := wsMessage{
			Type:   "private_message",
			Sender: testutil.TestUserID,
			Data:   raw,
		}
		wsRaw, _ := json.Marshal(wsMsg)

		hub.broadcast <- wsRaw

		select {
		case msg := <-receiver.send:
			var wrap struct {
				Type string         `json:"type"`
				Data models.Message `json:"data"`
			}
			if err := json.Unmarshal(msg, &wrap); err != nil {
				t.Fatalf("Failed to unmarshal received message: %v", err)
			}
			if wrap.Type != "private_message" || wrap.Data.Body != "hello" {
				t.Errorf("Incorrect message data: %+v", wrap)
			}
		case <-time.After(1 * time.Second):
			t.Fatal("Timeout waiting for private message")
		}
	})

	t.Run("Unregister", func(t *testing.T) {
		client := &Client{UserID: testutil.TestUserID, send: make(chan []byte, 1)}
		hub.register <- client

		testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
			online := hub.GetOnlineUsers()
			for _, id := range online {
				if id == testutil.TestUserID {
					return true
				}
			}
			return false
		})

		hub.unregister <- client

		testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
			online := hub.GetOnlineUsers()
			for _, id := range online {
				if id == testutil.TestUserID {
					return false
				}
			}
			return true
		})
	})
}

//--------------------------------------------------------------------------------------|

// TestHub_BroadcastStatusUpdate_WithUsername verifies that when users connect or disconnect,
// status updates containing their username are broadcast to other clients.
func TestHub_BroadcastStatusUpdate_WithUsername(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepo{
		GetByIDFunc: func(ctx context.Context, id int) (*models.User, error) {
			return &models.User{ID: id, Username: "testuser123"}, nil
		},
	}
	hub := NewHub(&testutil.MockChatRepo{}, &testutil.MockNotificationsRepo{}, mockUserRepo)
	go hub.Run(context.Background())

	// Create two clients
	client1 := &Client{UserID: testutil.TestUserID, send: make(chan []byte, 10)}
	client2 := &Client{UserID: testutil.TestUserID2, send: make(chan []byte, 10)}

	hub.register <- client1
	testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
		return len(hub.GetOnlineUsers()) >= 1
	})

	// Drain status update for client1's registration
	testutil.DrainChannel(client1.send)

	// Register client2 - this should broadcast a status update to client1
	hub.register <- client2
	testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
		return len(hub.GetOnlineUsers()) >= 2
	})

	// Client1 should receive a status update about client2
	select {
	case msg := <-client1.send:
		var update struct {
			Type string `json:"type"`
			Data struct {
				UserID   int    `json:"user_id"`
				Username string `json:"username"`
				Online   bool   `json:"online"`
			} `json:"data"`
		}
		if err := json.Unmarshal(msg, &update); err != nil {
			t.Fatalf("Failed to unmarshal status update: %v", err)
		}
		if update.Type != "status_update" {
			t.Errorf("Expected type 'status_update', got '%s'", update.Type)
		}
		if update.Data.UserID != testutil.TestUserID2 {
			t.Errorf("Expected user_id %d, got %d", testutil.TestUserID2, update.Data.UserID)
		}
		if update.Data.Username != "testuser123" {
			t.Errorf("Expected username 'testuser123', got '%s'", update.Data.Username)
		}
		if !update.Data.Online {
			t.Error("Expected online to be true")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("Timeout waiting for status update")
	}
}

//--------------------------------------------------------------------------------------|

// TestHub_ErrorPaths ensures that the Hub gracefully handles errors such as invalid JSON,
// malformed message data, repository failures, and marshaling issues.
func TestHub_ErrorPaths(t *testing.T) {
	mockChat := &testutil.MockChatRepo{}
	hub := NewHub(mockChat, &testutil.MockNotificationsRepo{}, nil)
	go hub.Run(context.Background())

	t.Run("Invalid Inbound JSON", func(t *testing.T) {
		hub.broadcast <- []byte("invalid json")
		// No easy way to verify as it just logs, but we ensure it doesn't crash
	})

	t.Run("Invalid Private Message Data", func(t *testing.T) {
		wsMsg := wsMessage{
			Type:   "private_message",
			Sender: testutil.TestUserID,
			Data:   []byte("invalid data json"),
		}
		raw, _ := json.Marshal(wsMsg)
		hub.broadcast <- raw
	})

	t.Run("SaveMessage Error", func(t *testing.T) {
		mockChat.SaveMessageFunc = func(ctx context.Context, senderID, receiverID int, body string, imageURL *string) (*models.Message, error) {
			return nil, context.DeadlineExceeded
		}
		msgData := map[string]interface{}{"receiver_id": testutil.TestUserID2, "body": "fail"}
		raw, _ := json.Marshal(msgData)
		wsMsg := wsMessage{Type: "private_message", Sender: testutil.TestUserID, Data: raw}
		wsRaw, _ := json.Marshal(wsMsg)
		hub.broadcast <- wsRaw
	})

	t.Run("JSON Marshal Error", func(t *testing.T) {
		// Passing a channel will cause json.Marshal to fail
		hub.broadcastJSON(make(chan int))
		hub.SendToUserJSON(testutil.TestUserID, make(chan int))
	})
}

//--------------------------------------------------------------------------------------|

// TestHub_Shutdown verifies that the Hub correctly cleans up resources and closes
// client connections when the context is canceled.
func TestHub_Shutdown(t *testing.T) {
	hub := NewHub(&testutil.MockChatRepo{}, &testutil.MockNotificationsRepo{}, nil)
	ctx, cancel := context.WithCancel(context.Background())

	client := &Client{UserID: 1, send: make(chan []byte, 1)}
	hub.clients[1] = client

	go hub.Run(ctx)
	cancel()

	// Wait for loop to exit and cleanup
	testutil.WaitForCondition(t, testutil.DefaultTestTimeout, func() bool {
		hub.mu.Lock()
		defer hub.mu.Unlock()
		return len(hub.clients) == 0
	})
}
