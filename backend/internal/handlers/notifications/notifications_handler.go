package notifications

import (
	"errors"
	"net/http"
	"social-network/internal/models"
	"social-network/internal/web"
)

//--------------------------------------------------------------------------------------|

type Handler struct {
	Service models.NotificationService
}

//--------------------------------------------------------------------------------------|

func NewHandler(svc models.NotificationService) *Handler {
	return &Handler{Service: svc}
}

//--------------------------------------------------------------------------------------|

func (h *Handler) List(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	limit, offset := web.GetLimitOffset(r)

	notifications, err := h.Service.GetNotifications(r.Context(), identity.ID, limit, offset)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, notifications)
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *Handler) GetUnreadCount(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	count, err := h.Service.GetUnreadCount(r.Context(), identity.ID)
	if err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]int{"count": count})
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *Handler) MarkAsRead(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	notificationID, err := web.PathInt64(r, "id")
	if err != nil {
		return web.StatusError{Code: http.StatusBadRequest, Err: err}
	}

	if err := h.Service.MarkAsRead(r.Context(), notificationID, identity.ID); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	return nil
}

//--------------------------------------------------------------------------------------|

func (h *Handler) MarkAllAsRead(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
	if identity == nil {
		return web.StatusError{Code: http.StatusUnauthorized, Err: errors.New("authentication required")}
	}

	if err := h.Service.MarkAllAsRead(r.Context(), identity.ID); err != nil {
		return err
	}

	web.JSONResponse(w, http.StatusOK, map[string]string{"status": "ok"})
	return nil
}
