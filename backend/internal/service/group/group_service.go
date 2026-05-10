package groupsvc

import (
	"context"
	"database/sql"
	"errors"
	"social-network/internal/models"
	"social-network/internal/utils"
)

//--------------------------------------------------------------------------------------|

type GroupService struct {
	Repo         models.GroupRepo
	NotifService models.NotificationService
	Hub          models.Hub
	UserRepo     models.UserRepo
	DB           *sql.DB
}

//--------------------------------------------------------------------------------------|

func NewGroupService(repo models.GroupRepo, notif models.NotificationService, hub models.Hub, userRepo models.UserRepo, db *sql.DB) *GroupService {
	return &GroupService{Repo: repo, NotifService: notif, Hub: hub, UserRepo: userRepo, DB: db}
}

//--------------------------------------------------------------------------------------|
// Group CRUD
//--------------------------------------------------------------------------------------|

// CreateGroup creates a new group and automatically adds the creator as a member.
func (s *GroupService) CreateGroup(ctx context.Context, creatorID int64, title, description string) (*models.Group, error) {
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

	if err == nil {
		s.Hub.UpdateUserGroups(creatorID)
	}

	return group, err
}

//--------------------------------------------------------------------------------------|

// GetGroup returns a group by ID. Verifies that the user is a member or the creator.
func (s *GroupService) GetGroup(ctx context.Context, id, userID int64) (*models.Group, error) {
	group, err := s.Repo.GetGroupByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrGroupNotFound
		}
		return nil, err
	}

	isMember, err := s.Repo.IsMember(ctx, id, userID)
	if err != nil {
		return nil, err
	}

	// Guard: Only creator or members can access
	if !isMember && group.CreatorID != userID {
		return nil, models.ErrNotMember
	}

	return group, nil
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
func (s *GroupService) GetMembers(ctx context.Context, groupID, userID int64) ([]models.GroupMember, error) {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group.CreatorID != userID {
		isMember, err := s.Repo.IsMember(ctx, groupID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, models.ErrNotMember
		}
	}
	members, err := s.Repo.GetMembers(ctx, groupID)
	if err != nil {
		return nil, err
	}

	for i := range members {
		members[i].AvatarURL = utils.FormatAvatarURL(members[i].AvatarURL)
	}

	return members, nil
}

//--------------------------------------------------------------------------------------|

// LeaveGroup removes a user from a group. Creators cannot leave their own group.
func (s *GroupService) LeaveGroup(ctx context.Context, groupID, userID int64) error {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return err
	}
	if group.CreatorID == userID {
		return &models.AuthorizationError{UserID: userID}
	}
	if err := s.Repo.RemoveMember(ctx, groupID, userID); err != nil {
		return err
	}
	s.Hub.UpdateUserGroups(userID)
	return nil
}

//--------------------------------------------------------------------------------------|
// Invitations
//--------------------------------------------------------------------------------------|

// InviteUser invites a user to a group. Only members can invite.
func (s *GroupService) InviteUser(ctx context.Context, groupID, inviterID, inviteeID int64) error {
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
	s.notify(inviteeID, inviterID, "group", inv.ID, "invite")

	return nil
}

//--------------------------------------------------------------------------------------|

// GetPendingInvitations returns all pending invitations for a user.
func (s *GroupService) GetPendingInvitations(ctx context.Context, userID int64) ([]models.GroupInvitation, error) {
	return s.Repo.GetPendingInvitations(ctx, userID)
}

//--------------------------------------------------------------------------------------|

// RespondToInvitation processes an accept/decline response to a group invitation.
func (s *GroupService) RespondToInvitation(ctx context.Context, invitationID, userID int64, accept bool) error {
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

	if err := utils.WithTx(ctx, s.DB, func(tx *sql.Tx) error {
		txRepo := s.Repo.WithTx(tx)

		if err := txRepo.UpdateInvitationStatus(ctx, invitationID, status); err != nil {
			return err
		}

		if accept {
			return txRepo.AddMember(ctx, inv.GroupID, userID, "member")
		}
		return nil
	}); err != nil {
		return err
	}

	if accept {
		s.Hub.UpdateUserGroups(userID)
	}
	return nil
}

//--------------------------------------------------------------------------------------|
// Join Requests
//--------------------------------------------------------------------------------------|

// RequestJoin creates a join request for a group.
func (s *GroupService) RequestJoin(ctx context.Context, groupID, userID int64) error {
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
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err == nil {
		s.notify(group.CreatorID, userID, "group", req.ID, "request")
	}

	return nil
}

//--------------------------------------------------------------------------------------|

// GetPendingJoinRequests returns all pending join requests for a group.
func (s *GroupService) GetPendingJoinRequests(ctx context.Context, groupID, requestingUserID int64) ([]models.GroupJoinRequest, error) {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
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
func (s *GroupService) RespondToJoinRequest(ctx context.Context, requestID, creatorID int64, accept bool) error {
	req, err := s.Repo.GetJoinRequestByID(ctx, requestID)
	if errors.Is(err, sql.ErrNoRows) {
		return models.ErrJoinRequestNotFound
	}
	if err != nil {
		return err
	}

	group, err := s.Repo.GetGroupByID(ctx, req.GroupID)
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
		s.Hub.UpdateUserGroups(req.UserID)
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
	if err := s.Repo.CreateEvent(ctx, event); err != nil {
		return err
	}

	// Notify all members
	members, err := s.Repo.GetMembers(ctx, event.GroupID)
	if err == nil {
		for _, member := range members {
			if member.UserID != event.CreatorID {
				s.notify(member.UserID, event.CreatorID, "group", event.GroupID, "event")
			}
		}
	}

	return nil
}

//--------------------------------------------------------------------------------------|

func (s *GroupService) GetGroupEvents(ctx context.Context, groupID, userID int64) ([]models.GroupEvent, error) {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group.CreatorID != userID {
		isMember, err := s.Repo.IsMember(ctx, groupID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, models.ErrNotMember
		}
	}
	return s.Repo.GetGroupEvents(ctx, groupID)
}

//--------------------------------------------------------------------------------------|

func (s *GroupService) RespondToEvent(ctx context.Context, eventID, userID int64, response string) error {
	event, err := s.Repo.GetEventByID(ctx, eventID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.ErrEventNotFound
		}
		return err
	}

	isMember, err := s.Repo.IsMember(ctx, event.GroupID, userID)
	if err != nil {
		return err
	}

	// Creator of the group also counts as a member for access
	group, err := s.Repo.GetGroupByID(ctx, event.GroupID)
	if err != nil {
		return err
	}

	if !isMember && group.CreatorID != userID {
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
// Group Messages
//--------------------------------------------------------------------------------------|

func (s *GroupService) SendGroupMessage(ctx context.Context, groupID, senderID int64, body string, imageURL *string) (*models.GroupMessage, error) {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	isMember, err := s.Repo.IsMember(ctx, groupID, senderID)
	if err != nil {
		return nil, err
	}

	if !isMember && group.CreatorID != senderID {
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
	// Mark as read for the sender
	_ = s.Repo.UpdateLastReadID(ctx, groupID, senderID, msg.ID)

	return msg, nil
}

//--------------------------------------------------------------------------------------|

// GetGroupMessages returns paginated messages for a group.
func (s *GroupService) GetGroupMessages(ctx context.Context, groupID, userID int64, limit, offset int) ([]models.GroupMessage, error) {
	group, err := s.Repo.GetGroupByID(ctx, groupID)
	if err != nil {
		return nil, err
	}

	if group.CreatorID != userID {
		isMember, err := s.Repo.IsMember(ctx, groupID, userID)
		if err != nil {
			return nil, err
		}
		if !isMember {
			return nil, models.ErrNotMember
		}
	}
	return s.Repo.GetGroupMessages(ctx, groupID, limit, offset)
}

//--------------------------------------------------------------------------------------|

func (s *GroupService) GetUnreadCounts(ctx context.Context, userID int64) ([]models.GroupUnreadCount, error) {
	return s.Repo.GetUnreadCounts(ctx, userID)
}

//--------------------------------------------------------------------------------------|

func (s *GroupService) MarkAsRead(ctx context.Context, groupID, userID int64) error {
	msgs, err := s.Repo.GetGroupMessages(ctx, groupID, 1, 0)
	if err != nil || len(msgs) == 0 {
		return err
	}
	return s.Repo.UpdateLastReadID(ctx, groupID, userID, msgs[0].ID)
}

//--------------------------------------------------------------------------------------|

func (s *GroupService) notify(userID, actorID int64, targetType string, targetID int64, notifType string) {
	go func() {
		username := ""
		if actor, err := s.UserRepo.GetByID(context.Background(), actorID); err == nil {
			username = actor.Nickname
		}
		s.NotifService.Notify(context.Background(), userID, actorID, username, targetType, targetID, notifType)
	}()
}
