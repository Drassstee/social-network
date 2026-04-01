// Package websocket provides real-time communication capabilities using WebSockets,
// including chat, notifications, and user status tracking.
package websocket

import (
	"context"
	"forum/internal/testutil"
	"testing"
)

//--------------------------------------------------------------------------------------|

// TestSqlChatRepository verifies the chat repository implementation, covering
// message persistence and retrieval for conversation history.
func TestSqlChatRepository(t *testing.T) {
	db, cleanup := testutil.SetupTestDBWithData(t)
	defer cleanup()

	repo := NewChatRepository(db)
	ctx := context.Background()
	aliceID := testutil.TestUserID
	bobID := testutil.TestUserID2

	t.Run("Save and Get Messages", func(t *testing.T) {
		msg, err := repo.SaveMessage(ctx, aliceID, bobID, "Hello Bob", nil)
		if err != nil {
			t.Errorf("SaveMessage failed: %v", err)
		}
		if msg != nil && msg.Username != testutil.ValidTestUsername {
			t.Errorf("Expected username '%s', got '%s'", testutil.ValidTestUsername, msg.Username)
		}

		// Save message with image
		imageURL := "/uploads/test.png"
		msgWithImage, err := repo.SaveMessage(ctx, bobID, aliceID, "Look at this", &imageURL)
		if err != nil {
			t.Errorf("SaveMessage with image failed: %v", err)
		}
		if msgWithImage.ImageURL == nil || *msgWithImage.ImageURL != imageURL {
			t.Errorf("Expected image URL '%s', got '%v'", imageURL, msgWithImage.ImageURL)
		}

		msgs, err := repo.GetMessages(ctx, aliceID, bobID, 10, 0)
		if err != nil {
			t.Errorf("GetMessages failed: %v", err)
		}
		if len(msgs) != 2 {
			t.Errorf("Expected 2 messages, got %d", len(msgs))
		}
		if msgs[0].ImageURL == nil || *msgs[0].ImageURL != imageURL {
			t.Errorf("Expected image URL '%s' for latest message, got '%v'", imageURL, msgs[0].ImageURL)
		}
	})
}
