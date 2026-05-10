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

type GroupHandler struct {
	Service models.GroupService
}

//--------------------------------------------------------------------------------------|

func NewGroupHandler(svc models.GroupService) *GroupHandler {
	return &GroupHandler{Service: svc}
}

//--------------------------------------------------------------------------------------|

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

func (h *GroupHandler) GetGroup(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	group, err := h.Service.GetGroup(r.Context(), groupID, identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, group)
	return nil
}

//--------------------------------------------------------------------------------------|

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

func (h *GroupHandler) GetMembers(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	members, err := h.Service.GetMembers(r.Context(), groupID, identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, members)
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *GroupHandler) LeaveGroup(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.PathInt64(r, "id")
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

func (h *GroupHandler) InviteUser(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var body struct {
		UserID int64 `json:"user_id"`
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

func (h *GroupHandler) RespondToInvitation(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	invID, err := web.PathInt64(r, "id")
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

func (h *GroupHandler) RequestJoin(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.PathInt64(r, "id")
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

func (h *GroupHandler) GetPendingJoinRequests(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.PathInt64(r, "id")
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

func (h *GroupHandler) RespondToJoinRequest(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	reqID, err := web.PathInt64(r, "id")
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

func (h *GroupHandler) CreateEvent(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	groupID, err := web.PathInt64(r, "id")
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

func (h *GroupHandler) GetGroupEvents(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	events, err := h.Service.GetGroupEvents(r.Context(), groupID, identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, events)
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *GroupHandler) RespondToEvent(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	eventID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	var body struct {
		Response string `json:"response"`
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

func (h *GroupHandler) GetGroupMessages(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	limit, offset := web.GetLimitOffset(r)

	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	messages, err := h.Service.GetGroupMessages(r.Context(), groupID, identity.ID, limit, offset)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, messages)
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *GroupHandler) GetUnreadCounts(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	counts, err := h.Service.GetUnreadCounts(r.Context(), identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, counts)
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *GroupHandler) MarkAsRead(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	groupID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	if err := h.Service.MarkAsRead(r.Context(), groupID, identity.ID); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	return nil
}
