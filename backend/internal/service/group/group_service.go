// Package groupsvc contains the business logic for group management.
package groupsvc

import (
	"context"
	"database/sql"
	"errors"
	"social-network/internal/models"
	"social-network/internal/utils"
)

//--------------------------------------------------------------------------------------|

// GroupService encapsulates group business rules.
type GroupService struct {
	Repo         models.GroupRepo
	NotifService models.NotificationService
	UserRepo     models.UserRepo
	DB           *sql.DB
}

//--------------------------------------------------------------------------------------|

// NewGroupService creates a new GroupService.
func NewGroupService(repo models.GroupRepo, notif models.NotificationService, userRepo models.UserRepo, db *sql.DB) *GroupService {
	return &GroupService{Repo: repo, NotifService: notif, UserRepo: userRepo, DB: db}
}

//--------------------------------------------------------------------------------------|
// Group CRUD
//--------------------------------------------------------------------------------------|

// CreateGroup creates a new group and automatically adds the creator as a member.
func (s *GroupService) CreateGroup(ctx context.Context, creatorID int, title, description string) (*models.Group, error) {
	if title == "" {
		return nil, &models.ValidationError{Field: "title", Message: "group title is required"}
	}

	group := &models.Group{
		CreatorID:   creatorID,
		Title:       title,
		Description: description,
	}

	err := utils.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		txRepo := s.Repo.WithTx(tx)

		if err := txRepo.CreateGroup(ctx, group); err != nil {
			return err
		}

		// Creator is automatically a member with the "creator" role.
		return txRepo.AddMember(ctx, group.ID, creatorID, "creator")
	})

	return group, err
}

//--------------------------------------------------------------------------------------|

// GetGroup returns a group by ID.
func (s *GroupService) GetGroup(ctx context.Context, id int) (*models.Group, error) {
	group, err := s.Repo.GetGroupByID(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, models.ErrGroupNotFound
	}
	return group, err
}

//--------------------------------------------------------------------------------------|

// ListGroups returns a paginated list of groups.
func (s *GroupService) ListGroups(ctx context.Context, limit, offset int) ([]models.Group, error) {
	return s.Repo.ListGroups(ctx, limit, offset)
}

//--------------------------------------------------------------------------------------|
// Membership
//--------------------------------------------------------------------------------------|

// GetMembers returns the members of a group.
func (s *GroupService) GetMembers(ctx context.Context, groupID int) ([]models.GroupMember, error) {
	return s.Repo.GetMembers(ctx, groupID)
}

//--------------------------------------------------------------------------------------|

// LeaveGroup removes a user from a group. Creators cannot leave their own group.
func (s *GroupService) LeaveGroup(ctx context.Context, groupID, userID int) error {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return err
	}
	if group.CreatorID == userID {
		return &models.AuthorizationError{UserID: userID}
	}
	return s.Repo.RemoveMember(ctx, groupID, userID)
}

//--------------------------------------------------------------------------------------|
// Invitations
//--------------------------------------------------------------------------------------|

// InviteUser invites a user to a group. Only members can invite.
func (s *GroupService) InviteUser(ctx context.Context, groupID, inviterID, inviteeID int) error {
	isMember, err := s.Repo.IsMember(ctx, groupID, inviterID)
	if err != nil {
		return err
	}
	if !isMember {
		return models.ErrNotMember
	}

	alreadyMember, err := s.Repo.IsMember(ctx, groupID, inviteeID)
	if err != nil {
		return err
	}
	if alreadyMember {
		return models.ErrAlreadyMember
	}

	inv := &models.GroupInvitation{
		GroupID:   groupID,
		InviterID: inviterID,
		InviteeID: inviteeID,
	}
	if err := s.Repo.CreateInvitation(ctx, inv); err != nil {
		return err
	}

	// Trigger notification
	s.notify(inviteeID, inviterID, "group", groupID, "invite")

	return nil
}

//--------------------------------------------------------------------------------------|

// GetPendingInvitations returns all pending invitations for a user.
func (s *GroupService) GetPendingInvitations(ctx context.Context, userID int) ([]models.GroupInvitation, error) {
	return s.Repo.GetPendingInvitations(ctx, userID)
}

//--------------------------------------------------------------------------------------|

// RespondToInvitation processes an accept/decline response to a group invitation.
func (s *GroupService) RespondToInvitation(ctx context.Context, invitationID, userID int, accept bool) error {
	inv, err := s.Repo.GetInvitationByID(ctx, invitationID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrInvitationNotFound
	}
	if err != nil {
		return err
	}

	if inv.InviteeID != userID {
		return &models.AuthorizationError{UserID: userID}
	}
	if inv.Status != "pending" {
		return models.ErrInvitationNotFound
	}

	status := "declined"
	if accept {
		status = "accepted"
	}

	return utils.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		txRepo := s.Repo.WithTx(tx)

		if err := txRepo.UpdateInvitationStatus(ctx, invitationID, status); err != nil {
			return err
		}

		if accept {
			return txRepo.AddMember(ctx, inv.GroupID, userID, "member")
		}
		return nil
	})
}

//--------------------------------------------------------------------------------------|
// Join Requests
//--------------------------------------------------------------------------------------|

// RequestJoin creates a join request for a group.
func (s *GroupService) RequestJoin(ctx context.Context, groupID, userID int) error {
	isMember, err := s.Repo.IsMember(ctx, groupID, userID)
	if err != nil {
		return err
	}
	if isMember {
		return models.ErrAlreadyMember
	}

	req := &models.GroupJoinRequest{GroupID: groupID, UserID: userID}
	if err := s.Repo.CreateJoinRequest(ctx, req); err != nil {
		return err
	}

	// Trigger notification to group creator
	group, err := s.GetGroup(ctx, groupID)
	if err == nil {
		s.notify(group.CreatorID, userID, "group", groupID, "request")
	}

	return nil
}

//--------------------------------------------------------------------------------------|

// GetPendingJoinRequests returns all pending join requests for a group.
func (s *GroupService) GetPendingJoinRequests(ctx context.Context, groupID, requestingUserID int) ([]models.GroupJoinRequest, error) {
	group, err := s.GetGroup(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group.CreatorID != requestingUserID {
		return nil, models.ErrNotGroupCreator
	}
	return s.Repo.GetPendingJoinRequests(ctx, groupID)
}

//--------------------------------------------------------------------------------------|

// RespondToJoinRequest processes a join request (creator-only action).
func (s *GroupService) RespondToJoinRequest(ctx context.Context, requestID, creatorID int, accept bool) error {
	req, err := s.Repo.GetJoinRequestByID(ctx, requestID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrJoinRequestNotFound
	}
	if err != nil {
		return err
	}

	group, err := s.GetGroup(ctx, req.GroupID)
	if err != nil {
		return err
	}
	if group.CreatorID != creatorID {
		return models.ErrNotGroupCreator
	}
	if req.Status != "pending" {
		return models.ErrJoinRequestNotFound
	}

	status := "declined"
	if accept {
		status = "accepted"
	}

	if err := utils.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		txRepo := s.Repo.WithTx(tx)

		if err := txRepo.UpdateJoinRequestStatus(ctx, requestID, status); err != nil {
			return err
		}

		if accept {
			return txRepo.AddMember(ctx, req.GroupID, req.UserID, "member")
		}
		return nil
	}); err != nil {
		return err
	}

	// Trigger notification to requester
	notifType := "decline"
	if accept {
		notifType = "accept"
	}
	s.notify(req.UserID, creatorID, "group", req.GroupID, notifType)

	return nil
}

//--------------------------------------------------------------------------------------|
// Events
//--------------------------------------------------------------------------------------|

// CreateEvent creates a new group event. Only members can create events.
func (s *GroupService) CreateEvent(ctx context.Context, event *models.GroupEvent) error {
	isMember, err := s.Repo.IsMember(ctx, event.GroupID, event.CreatorID)
	if err != nil {
		return err
	}
	if !isMember {
		return models.ErrNotMember
	}
	if event.Title == "" {
		return &models.ValidationError{Field: "title", Message: "event title is required"}
	}
	return s.Repo.CreateEvent(ctx, event)
}

//--------------------------------------------------------------------------------------|

// GetGroupEvents returns all events for a group.
func (s *GroupService) GetGroupEvents(ctx context.Context, groupID int) ([]models.GroupEvent, error) {
	return s.Repo.GetGroupEvents(ctx, groupID)
}

//--------------------------------------------------------------------------------------|

// RespondToEvent records a going/not_going RSVP. Only members can respond.
func (s *GroupService) RespondToEvent(ctx context.Context, eventID, userID int, response string) error {
	event, err := s.Repo.GetEventByID(ctx, eventID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrEventNotFound
	}
	if err != nil {
		return err
	}

	isMember, err := s.Repo.IsMember(ctx, event.GroupID, userID)
	if err != nil {
		return err
	}
	if !isMember {
		return models.ErrNotMember
	}

	if response != "going" && response != "not_going" {
		return &models.ValidationError{Field: "response", Message: "must be 'going' or 'not_going'"}
	}

	return s.Repo.RespondToEvent(ctx, &models.GroupEventResponse{
		EventID:  eventID,
		UserID:   userID,
		Response: response,
	})
}

//--------------------------------------------------------------------------------------|

// GetEventResponses returns all RSVP responses for an event.
func (s *GroupService) GetEventResponses(ctx context.Context, eventID int) ([]models.GroupEventResponse, error) {
	return s.Repo.GetEventResponses(ctx, eventID)
}

//--------------------------------------------------------------------------------------|
// Group Messages
//--------------------------------------------------------------------------------------|

// SendGroupMessage sends a message to a group. Only members can send.
func (s *GroupService) SendGroupMessage(ctx context.Context, groupID, senderID int, body string, imageURL *string) (*models.GroupMessage, error) {
	isMember, err := s.Repo.IsMember(ctx, groupID, senderID)
	if err != nil {
		return nil, err
	}
	if !isMember {
		return nil, models.ErrNotMember
	}

	msg := &models.GroupMessage{
		GroupID:  groupID,
		SenderID: senderID,
		Body:     body,
		ImageURL: imageURL,
	}

	if err := s.Repo.SaveGroupMessage(ctx, msg); err != nil {
		return nil, err
	}
	return msg, nil
}

//--------------------------------------------------------------------------------------|

// GetGroupMessages returns paginated messages for a group.
func (s *GroupService) GetGroupMessages(ctx context.Context, groupID, limit, offset int) ([]models.GroupMessage, error) {
	return s.Repo.GetGroupMessages(ctx, groupID, limit, offset)
}

//--------------------------------------------------------------------------------------|

func (s *GroupService) notify(userID, actorID int, targetType string, targetID int, notifType string) {
	go func() {
		username := ""
		if actor, err := s.UserRepo.GetByID(context.Background(), actorID); err == nil {
			username = actor.Username
		}
		s.NotifService.Notify(context.Background(), userID, actorID, username, targetType, targetID, notifType)
	}()
}
