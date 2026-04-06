package web

import (
	"database/sql"
	"errors"
	"net/http"
	"social-network/internal/models"
)

//--------------------------------------------------------------------------------------|

// MapError converts standard domain errors to HTTP status codes.
func MapError(err error) int {
	var se StatusError
	if errors.As(err, &se) {
		return se.Code
	}

	var ve *models.ValidationError
	if errors.As(err, &ve) {
		return http.StatusBadRequest
	}

	var nfe *models.NotFoundError
	if errors.As(err, &nfe) {
		return http.StatusNotFound
	}

	var ae *models.AuthorizationError
	if errors.As(err, &ae) {
		if ae.UserID == 0 {
			return http.StatusUnauthorized
		}
		return http.StatusForbidden
	}

	if errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, models.ErrPostNotFound) ||
		errors.Is(err, models.ErrCommentNotFound) ||
		errors.Is(err, models.ErrParentCommentNotFound) ||
		errors.Is(err, models.ErrUserNotFound) ||
		errors.Is(err, models.ErrForeignKeyConstraint) {
		return http.StatusNotFound
	}

	if errors.Is(err, models.ErrInvalidCredentials) {
		return http.StatusUnauthorized
	}

	if errors.Is(err, models.ErrUserAlreadyExists) ||
		errors.Is(err, models.ErrUniqueConstraint) {
		return http.StatusConflict
	}

	if errors.Is(err, models.ErrUsernameLength) ||
		errors.Is(err, models.ErrUsernameFormat) ||
		errors.Is(err, models.ErrPasswordLength) ||
		errors.Is(err, models.ErrPasswordComplexity) ||
		errors.Is(err, models.ErrInvalidEmail) ||
		errors.Is(err, models.ErrInvalidPostTitle) ||
		errors.Is(err, models.ErrInvalidPostBody) ||
		errors.Is(err, models.ErrInvalidCommentBody) ||
		errors.Is(err, models.ErrInvalidCategoryName) ||
		errors.Is(err, models.ErrImageTooBig) ||
		errors.Is(err, models.ErrInvalidImageFormat) ||
		errors.Is(err, models.ErrCommentDepthExceeded) {
		return http.StatusBadRequest
	}

	return http.StatusInternalServerError
}

//--------------------------------------------------------------------------------------|

// HandleError centralizes error writing for both JSON and HTML responses.
func HandleError(w http.ResponseWriter, r *http.Request, err error) {
	code := MapError(err)

	if r.Header.Get("Accept") == "application/json" || r.Header.Get("Content-Type") == "application/json" {
		response := map[string]any{"error": err.Error()}

		// Special handling for validation errors to provide field context
		var ve *models.ValidationError
		if errors.As(err, &ve) {
			response["field"] = ve.Field
			response["message"] = ve.Message
		}

		if code == http.StatusInternalServerError {
			response["error"] = "Internal server error"
		}

		JSONResponse(w, code, response)
		return
	}

	if code == http.StatusInternalServerError {
		InternalServerError(w, r, err)
	} else if code == http.StatusNotFound {
		NotFound(w, r)
	} else {
		renderStaticError(w, code, err.Error())
	}
}
