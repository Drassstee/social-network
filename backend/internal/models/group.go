package models

import (
	"context"
	"time"
)

//--------------------------------------------------------------------------------------|

// Group represents a social group that users can create and join.
type Group struct {
	ID          int       `json:"id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

//--------------------------------------------------------------------------------------|

// GroupMember represents a user's membership in a group.
type GroupMember struct {
	GroupID  int       `json:"group_id"`
	UserID   int       `json:"user_id"`
	Role     string    `json:"role"` // "creator" or "member"
	JoinedAt time.Time `json:"joined_at"`

	// Joined fields for display purposes.
	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

//--------------------------------------------------------------------------------------|

// GroupInvitation represents an invitation from a group member to another user.
type GroupInvitation struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	InviterID int       `json:"inviter_id"`
	InviteeID int       `json:"invitee_id"`
	Status    string    `json:"status"` // "pending", "accepted", "declined"
	CreatedAt time.Time `json:"created_at"`

	// Joined fields for display purposes.
	GroupTitle  string `json:"group_title,omitempty"`
	InviterName string `json:"inviter_name,omitempty"`
}

//--------------------------------------------------------------------------------------|

// GroupJoinRequest represents a request from a user to join a group.
type GroupJoinRequest struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	UserID    int       `json:"user_id"`
	Status    string    `json:"status"` // "pending", "accepted", "declined"
	CreatedAt time.Time `json:"created_at"`

	// Joined fields for display purposes.
	Username string `json:"username,omitempty"`
}

//--------------------------------------------------------------------------------------|

// GroupEvent represents an event within a group.
type GroupEvent struct {
	ID          int       `json:"id"`
	GroupID     int       `json:"group_id"`
	CreatorID   int       `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`

	// Aggregated response counts (populated at query time).
	GoingCount    int `json:"going_count"`
	NotGoingCount int `json:"not_going_count"`
}

//--------------------------------------------------------------------------------------|

// GroupEventResponse represents a member's RSVP to a group event.
type GroupEventResponse struct {
	EventID  int    `json:"event_id"`
	UserID   int    `json:"user_id"`
	Response string `json:"response"` // "going" or "not_going"

	Username string `json:"username,omitempty"`
}

//--------------------------------------------------------------------------------------|

// GroupMessage represents a chat message within a group.
type GroupMessage struct {
	ID        int       `json:"id"`
	GroupID   int       `json:"group_id"`
	SenderID  int       `json:"sender_id"`
	Body      string    `json:"body"`
	ImageURL  *string   `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	Username string `json:"username,omitempty"`
}

//--------------------------------------------------------------------------------------|
// Interfaces
//--------------------------------------------------------------------------------------|

// GroupRepo defines the contract for all group-related persistence.
type GroupRepo interface {
	// Transactions
	WithTx(tx any) GroupRepo

	// Group CRUD
	CreateGroup(ctx context.Context, group *Group) error
	GetGroupByID(ctx context.Context, id int) (*Group, error)
	ListGroups(ctx context.Context, limit, offset int) ([]Group, error)

	// Membership
	AddMember(ctx context.Context, groupID, userID int, role string) error
	RemoveMember(ctx context.Context, groupID, userID int) error
	GetMembers(ctx context.Context, groupID int) ([]GroupMember, error)
	IsMember(ctx context.Context, groupID, userID int) (bool, error)
	GetMemberGroupIDs(ctx context.Context, userID int) ([]int, error)

	// Invitations
	CreateInvitation(ctx context.Context, inv *GroupInvitation) error
	GetInvitationByID(ctx context.Context, id int) (*GroupInvitation, error)
	GetPendingInvitations(ctx context.Context, userID int) ([]GroupInvitation, error)
	UpdateInvitationStatus(ctx context.Context, id int, status string) error

	// Join Requests
	CreateJoinRequest(ctx context.Context, req *GroupJoinRequest) error
	GetJoinRequestByID(ctx context.Context, id int) (*GroupJoinRequest, error)
	GetPendingJoinRequests(ctx context.Context, groupID int) ([]GroupJoinRequest, error)
	UpdateJoinRequestStatus(ctx context.Context, id int, status string) error

	// Events
	CreateEvent(ctx context.Context, event *GroupEvent) error
	GetEventByID(ctx context.Context, id int) (*GroupEvent, error)
	GetGroupEvents(ctx context.Context, groupID int) ([]GroupEvent, error)
	RespondToEvent(ctx context.Context, resp *GroupEventResponse) error
	GetEventResponses(ctx context.Context, eventID int) ([]GroupEventResponse, error)

	// Group Messages
	SaveGroupMessage(ctx context.Context, msg *GroupMessage) error
	GetGroupMessages(ctx context.Context, groupID, limit, offset int) ([]GroupMessage, error)
}

// GroupService separates groups business logic from transport.
type GroupService interface {
	CreateGroup(ctx context.Context, creatorID int, title, description string) (*Group, error)
	GetGroup(ctx context.Context, id int) (*Group, error)
	ListGroups(ctx context.Context, limit, offset int) ([]Group, error)
	GetMembers(ctx context.Context, groupID int) ([]GroupMember, error)
	LeaveGroup(ctx context.Context, groupID, userID int) error
	InviteUser(ctx context.Context, groupID, inviterID, inviteeID int) error
	GetPendingInvitations(ctx context.Context, userID int) ([]GroupInvitation, error)
	RespondToInvitation(ctx context.Context, invitationID, userID int, accept bool) error
	RequestJoin(ctx context.Context, groupID, userID int) error
	GetPendingJoinRequests(ctx context.Context, groupID, requestingUserID int) ([]GroupJoinRequest, error)
	RespondToJoinRequest(ctx context.Context, requestID, creatorID int, accept bool) error
	CreateEvent(ctx context.Context, event *GroupEvent) error
	GetGroupEvents(ctx context.Context, groupID int) ([]GroupEvent, error)
	RespondToEvent(ctx context.Context, eventID, userID int, response string) error
	GetEventResponses(ctx context.Context, eventID int) ([]GroupEventResponse, error)
	SendGroupMessage(ctx context.Context, groupID, senderID int, body string, imageURL *string) (*GroupMessage, error)
	GetGroupMessages(ctx context.Context, groupID, limit, offset int) ([]GroupMessage, error)
}
