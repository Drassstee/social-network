// Package group provides HTTP handlers for group management endpoints.
package group

import (
	"encoding/json"
	"errors"
	"net/http"
	"social-network/internal/models"
	"social-network/internal/web"
	"time"
)

//--------------------------------------------------------------------------------------|

// GroupHandler handles all group-related HTTP requests.
type GroupHandler struct {
	Service models.GroupService
}

//--------------------------------------------------------------------------------------|

// NewGroupHandler creates a new GroupHandler.
func NewGroupHandler(svc models.GroupService) *GroupHandler {
	return &GroupHandler{Service: svc}
}

//--------------------------------------------------------------------------------------|
// Group CRUD
//--------------------------------------------------------------------------------------|

// CreateGroup handles POST /api/groups.
func (h *GroupHandler) CreateGroup(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	group, err := h.Service.CreateGroup(r.Context(), identity.ID, body.Title, body.Description)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusCreated, group)
	return nil
}

//--------------------------------------------------------------------------------------|

// GetGroup handles GET /api/groups/{id}.
func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2) // /api/groups/{id}
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	group, err := h.Service.GetGroup(r.Context(), groupID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, group)
	return nil
}

//--------------------------------------------------------------------------------------|

// ListGroups handles GET /api/groups.
func (h *GroupHandler) ListGroups(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	limit, offset := web.GetLimitOffset(r)

	groups, err := h.Service.ListGroups(r.Context(), limit, offset)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, groups)
	return nil
}

//--------------------------------------------------------------------------------------|
// Membership
//--------------------------------------------------------------------------------------|

// GetMembers handles GET /api/groups/{id}/members.
func (h *GroupHandler) GetMembers(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	members, err := h.Service.GetMembers(r.Context(), groupID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, members)
	return nil
}

//--------------------------------------------------------------------------------------|

// LeaveGroup handles POST /api/groups/{id}/leave.
func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if err := h.Service.LeaveGroup(r.Context(), groupID, identity.ID); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "left"})
	return nil
}

//--------------------------------------------------------------------------------------|
// Invitations
//--------------------------------------------------------------------------------------|

// InviteUser handles POST /api/groups/{id}/invite.
func (h *GroupHandler) InviteUser(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var body struct {
		UserID int `json:"user_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.UserID == 0 {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("invalid user_id")}
	}

	if err := h.Service.InviteUser(r.Context(), groupID, identity.ID, body.UserID); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusCreated, map[string]string{"status": "invited"})
	return nil
}

//--------------------------------------------------------------------------------------|

// GetPendingInvitations handles GET /api/groups/invitations.
func (h *GroupHandler) GetPendingInvitations(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	invitations, err := h.Service.GetPendingInvitations(r.Context(), identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, invitations)
	return nil
}

//--------------------------------------------------------------------------------------|

// RespondToInvitation handles POST /api/groups/invitations/{id}/respond.
func (h *GroupHandler) RespondToInvitation(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	invID, err := web.ExtractIDFromPath(r.URL.Path, 3) // /api/groups/invitations/{id}/respond
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var body struct {
		Accept bool `json:"accept"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	if err := h.Service.RespondToInvitation(r.Context(), invID, identity.ID, body.Accept); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "responded"})
	return nil
}

//--------------------------------------------------------------------------------------|
// Join Requests
//--------------------------------------------------------------------------------------|

// RequestJoin handles POST /api/groups/{id}/request.
func (h *GroupHandler) RequestJoin(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if err := h.Service.RequestJoin(r.Context(), groupID, identity.ID); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusCreated, map[string]string{"status": "requested"})
	return nil
}

//--------------------------------------------------------------------------------------|

// GetPendingJoinRequests handles GET /api/groups/{id}/requests.
func (h *GroupHandler) GetPendingJoinRequests(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	requests, err := h.Service.GetPendingJoinRequests(r.Context(), groupID, identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, requests)
	return nil
}

//--------------------------------------------------------------------------------------|

// RespondToJoinRequest handles POST /api/groups/requests/{id}/respond.
func (h *GroupHandler) RespondToJoinRequest(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	reqID, err := web.ExtractIDFromPath(r.URL.Path, 3) // /api/groups/requests/{id}/respond
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var body struct {
		Accept bool `json:"accept"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	if err := h.Service.RespondToJoinRequest(r.Context(), reqID, identity.ID, body.Accept); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "responded"})
	return nil
}

//--------------------------------------------------------------------------------------|
// Events
//--------------------------------------------------------------------------------------|

// CreateEvent handles POST /api/groups/{id}/events.
func (h *GroupHandler) CreateEvent(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var body struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		EventTime   string `json:"event_time"` // RFC3339 format
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	eventTime, err := time.Parse(time.RFC3339, body.EventTime)
	if err != nil {
		return &models.ValidationError{Field: "event_time", Message: "must be RFC3339 format"}
	}

	event := &models.GroupEvent{
		GroupID:     groupID,
		CreatorID:   identity.ID,
		Title:       body.Title,
		Description: body.Description,
		EventTime:   eventTime,
	}

	if err := h.Service.CreateEvent(r.Context(), event); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusCreated, event)
	return nil
}

//--------------------------------------------------------------------------------------|

// GetGroupEvents handles GET /api/groups/{id}/events.
func (h *GroupHandler) GetGroupEvents(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	events, err := h.Service.GetGroupEvents(r.Context(), groupID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, events)
	return nil
}

//--------------------------------------------------------------------------------------|

// RespondToEvent handles POST /api/groups/events/{id}/respond.
func (h *GroupHandler) RespondToEvent(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	eventID, err := web.ExtractIDFromPath(r.URL.Path, 3) // /api/groups/events/{id}/respond
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var body struct {
		Response string `json:"response"` // "going" or "not_going"
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: errors.New("invalid request body")}
	}

	if err := h.Service.RespondToEvent(r.Context(), eventID, identity.ID, body.Response); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "responded"})
	return nil
}

//--------------------------------------------------------------------------------------|
// Group Messages
//--------------------------------------------------------------------------------------|

// GetGroupMessages handles GET /api/groups/{id}/messages.
func (h *GroupHandler) GetGroupMessages(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.ExtractIDFromPath(r.URL.Path, 2)
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	limit, offset := web.GetLimitOffset(r)

	messages, err := h.Service.GetGroupMessages(r.Context(), groupID, limit, offset)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, messages)
	return nil
}
