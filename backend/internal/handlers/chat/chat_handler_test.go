// Package websocket provides real-time communication capabilities using WebSockets,
// including chat, notifications, and user status tracking.
package websocket

import (
	"context"
	"encoding/json"
	"forum/internal/models"
	"forum/internal/testutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

//--------------------------------------------------------------------------------------|

// TestChatHandler_GetOnlineUsers verifies that the endpoint correctly returns the list of
// online users from the Hub, including fetching user details from the repository.
func TestChatHandler_GetOnlineUsers(t *testing.T) {
	mockUserRepo := &testutil.MockUserRepo{
		GetByIDFunc: func(ctx context.Context, id int) (*models.User, error) {
			return &models.User{ID: id, Username: "user" + string(rune('0'+id))}, nil
		},
	}
	hub := NewHub(&testutil.MockChatRepo{}, &testutil.MockNotificationsRepo{}, mockUserRepo)
	handler := &ChatHandler{Hub: hub, UserRepo: mockUserRepo}

	hub.clients[testutil.TestUserID] = &Client{UserID: testutil.TestUserID}
	hub.clients[testutil.TestUserID2] = &Client{UserID: testutil.TestUserID2}

	req := httptest.NewRequest("GET", "/api/online-users", nil)
	w := httptest.NewRecorder()

	if err := handler.GetOnlineUsers(w, req, &models.UserIdentity{ID: testutil.TestUserID}); err != nil {
		t.Fatalf("GetOnlineUsers failed: %v", err)
	}

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var online []struct {
		ID       int    `json:"id"`
		Username string `json:"username"`
	}
	json.NewDecoder(w.Body).Decode(&online)
	if len(online) != 2 {
		t.Errorf("Expected 2 online users, got %d", len(online))
	}
}

//--------------------------------------------------------------------------------------|

// TestChatHandler_GetMessages ensures correct retrieval of chat history between two users,
// validating request parameters and response format.
func TestChatHandler_GetMessages(t *testing.T) {
	mockRepo := &testutil.MockChatRepo{
		GetMessagesFunc: func(ctx context.Context, user1ID, user2ID, limit, offset int) ([]models.Message, error) {
			return []models.Message{{ID: 1, Body: "test"}}, nil
		},
	}
	handler := &ChatHandler{Repo: mockRepo}

	t.Run("Valid request", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/messages?user_id=2", nil)
		// Set identity in context
		ctx := context.WithValue(req.Context(), models.UserKey, &models.UserIdentity{ID: testutil.TestUserID})
		req = req.WithContext(ctx)

		w := httptest.NewRecorder()

		if err := handler.GetMessages(w, req, &models.UserIdentity{ID: testutil.TestUserID}); err != nil {
			t.Fatalf("GetMessages failed: %v", err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		var msgs []models.Message
		json.NewDecoder(w.Body).Decode(&msgs)
		if len(msgs) != 1 || msgs[0].Body != "test" {
			t.Errorf("Incorrect messages returned: %+v", msgs)
		}
	})

	t.Run("Invalid user ID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/api/messages?user_id=abc", nil)
		ctx := context.WithValue(req.Context(), models.UserKey, &models.UserIdentity{ID: testutil.TestUserID})
		req = req.WithContext(ctx)
		w := httptest.NewRecorder()

		err := handler.GetMessages(w, req, &models.UserIdentity{ID: testutil.TestUserID})
		if err == nil {
			t.Fatal("Expected error for invalid user_id, got nil")
		}
	})
}
