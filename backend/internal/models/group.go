package models

import (
	"context"
	"time"
)

//--------------------------------------------------------------------------------------|

type Group struct {
	ID          int64     `json:"id"`
	CreatorID   int64     `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	MemberCount int       `json:"member_count"`
	CreatedAt   time.Time `json:"created_at"`
}

//--------------------------------------------------------------------------------------|

type GroupMember struct {
	GroupID  int64     `json:"group_id"`
	UserID   int64     `json:"user_id"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`

	Username  string `json:"username,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

//--------------------------------------------------------------------------------------|

type GroupInvitation struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	InviterID int64     `json:"inviter_id"`
	InviteeID int64     `json:"invitee_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`

	GroupTitle  string `json:"group_title,omitempty"`
	InviterName string `json:"inviter_name,omitempty"`
}

//--------------------------------------------------------------------------------------|

type GroupJoinRequest struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	UserID    int64     `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`

	Username string `json:"username,omitempty"`
}

// --------------------------------------------------------------------------------------|
type GroupEvent struct {
	ID          int64     `json:"id"`
	GroupID     int64     `json:"group_id"`
	CreatorID   int64     `json:"creator_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	EventTime   time.Time `json:"event_time"`
	CreatedAt   time.Time `json:"created_at"`

	GoingCount    int `json:"going_count"`
	NotGoingCount int `json:"not_going_count"`
}

//--------------------------------------------------------------------------------------|

type GroupEventResponse struct {
	EventID  int64  `json:"event_id"`
	UserID   int64  `json:"user_id"`
	Response string `json:"response"`

	Username string `json:"username,omitempty"`
}

//--------------------------------------------------------------------------------------|

type GroupMessage struct {
	ID        int64     `json:"id"`
	GroupID   int64     `json:"group_id"`
	SenderID  int64     `json:"sender_id"`
	Body      string    `json:"body"`
	ImageURL  *string   `json:"image_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`

	Username string `json:"username,omitempty"`
}

//--------------------------------------------------------------------------------------|

type GroupUnreadCount struct {
	GroupID     int64 `json:"group_id"`
	UnreadCount int   `json:"unread_count"`
}

//--------------------------------------------------------------------------------------|

type GroupRepo interface {
	// Transactions
	WithTx(tx any) GroupRepo

	// Group CRUD
	CreateGroup(ctx context.Context, group *Group) error
	GetGroupByID(ctx context.Context, id int64) (*Group, error)
	ListGroups(ctx context.Context, limit, offset int) ([]Group, error)

	// Membership
	AddMember(ctx context.Context, groupID, userID int64, role string) error
	RemoveMember(ctx context.Context, groupID, userID int64) error
	GetMembers(ctx context.Context, groupID int64) ([]GroupMember, error)
	IsMember(ctx context.Context, groupID, userID int64) (bool, error)
	GetMemberGroupIDs(ctx context.Context, userID int64) ([]int64, error)

	// Invitations
	CreateInvitation(ctx context.Context, inv *GroupInvitation) error
	GetInvitationByID(ctx context.Context, id int64) (*GroupInvitation, error)
	GetPendingInvitations(ctx context.Context, userID int64) ([]GroupInvitation, error)
	UpdateInvitationStatus(ctx context.Context, id int64, status string) error

	// Join Requests
	CreateJoinRequest(ctx context.Context, req *GroupJoinRequest) error
	GetJoinRequestByID(ctx context.Context, id int64) (*GroupJoinRequest, error)
	GetPendingJoinRequests(ctx context.Context, groupID int64) ([]GroupJoinRequest, error)
	UpdateJoinRequestStatus(ctx context.Context, id int64, status string) error

	// Events
	CreateEvent(ctx context.Context, event *GroupEvent) error
	GetEventByID(ctx context.Context, id int64) (*GroupEvent, error)
	GetGroupEvents(ctx context.Context, groupID int64) ([]GroupEvent, error)
	RespondToEvent(ctx context.Context, resp *GroupEventResponse) error

	// Group Messages
	SaveGroupMessage(ctx context.Context, msg *GroupMessage) error
	GetGroupMessages(ctx context.Context, groupID int64, limit, offset int) ([]GroupMessage, error)
	UpdateLastReadID(ctx context.Context, groupID, userID, msgID int64) error
	GetUnreadCounts(ctx context.Context, userID int64) ([]GroupUnreadCount, error)
}

//--------------------------------------------------------------------------------------|

type GroupService interface {
	CreateGroup(ctx context.Context, creatorID int64, title, description string) (*Group, error)
	GetGroup(ctx context.Context, id, userID int64) (*Group, error)
	ListGroups(ctx context.Context, limit, offset int) ([]Group, error)
	GetMembers(ctx context.Context, groupID, userID int64) ([]GroupMember, error)
	LeaveGroup(ctx context.Context, groupID, userID int64) error
	InviteUser(ctx context.Context, groupID, inviterID, inviteeID int64) error
	GetPendingInvitations(ctx context.Context, userID int64) ([]GroupInvitation, error)
	RespondToInvitation(ctx context.Context, invitationID, userID int64, accept bool) error
	RequestJoin(ctx context.Context, groupID, userID int64) error
	GetPendingJoinRequests(ctx context.Context, groupID, requestingUserID int64) ([]GroupJoinRequest, error)
	RespondToJoinRequest(ctx context.Context, requestID, creatorID int64, accept bool) error
	CreateEvent(ctx context.Context, event *GroupEvent) error
	GetGroupEvents(ctx context.Context, groupID, userID int64) ([]GroupEvent, error)
	RespondToEvent(ctx context.Context, eventID, userID int64, response string) error
	SendGroupMessage(ctx context.Context, groupID, senderID int64, body string, imageURL *string) (*GroupMessage, error)
	GetGroupMessages(ctx context.Context, groupID, userID int64, limit, offset int) ([]GroupMessage, error)
	GetUnreadCounts(ctx context.Context, userID int64) ([]GroupUnreadCount, error)
	MarkAsRead(ctx context.Context, groupID, userID int64) error
}
