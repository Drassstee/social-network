// Package web provides HTTP utilities including template rendering,
// error handling, image upload processing, and request parsing.
package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"social-network/internal/models"
	"social-network/internal/utils"
	"log"
	"net/http"
	"path/filepath"
	"strings"
)

//--------------------------------------------------------------------------------------|

var (
	defaultStaticDir = "./assets/static/"
)

var (
	UploadDir       = "./assets/static/uploads"
	UploadURLPrefix = "/static/uploads/"
)

const (
	// DefaultLimit is the default number of items per page for paginated requests.
	DefaultLimit = 10
	// MaxLimit is the maximum allowed items per page to prevent excessive memory usage.
	MaxLimit = 100

	// MaxImageSize is the maximum allowed size for uploaded images (20MB).
	MaxImageSize = 20 * 1024 * 1024

	// Path segment indices for ID extraction from relative URLs.
	SegmentPost         = 1 // e.g., /posts/{id} -> segment 1
	SegmentPostEdit     = 2 // e.g., /posts/edit/{id} -> segment 2
	SegmentComment      = 1 // e.g., /comments/{id} -> segment 1
	SegmentNotification = 1 // e.g., /notifications/{id} -> segment 1
)

//--------------------------------------------------------------------------------------|

// StatusError represents an HTTP error with a status code.
type StatusError struct {
	Code int
	Err  error
}

func (se StatusError) Error() string {
	return se.Err.Error()
}

//--------------------------------------------------------------------------------------|

// AppHandler is a specialized handler that returns an error and includes user identity.
type AppHandler func(http.ResponseWriter, *http.Request, *models.UserIdentity) error

// NewAppHandler wraps an AppHandler and implements http.HandlerFunc.
func NewAppHandler(h AppHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		identity := models.GetIdentity(r.Context())
		if err := h(w, r, identity); err != nil {
			HandleError(w, r, err)
		}
	}
}

//--------------------------------------------------------------------------------------|

// JSONResponse sends a JSON response with the specified status code.
func JSONResponse(w http.ResponseWriter, code int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(data)
}

//--------------------------------------------------------------------------------------|

// FileServer returns a handler that serves file contents from the given directory.
func FileServer(dir string) http.Handler {
	if dir == "" {
		dir = defaultStaticDir
	}
	return http.FileServer(http.Dir(dir))
}

//--------------------------------------------------------------------------------------|

// InternalServerError logs the error and sends a generic 500 response.
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	log.Printf("[server error] %s %s: %v", r.Method, r.URL.Path, err)
	renderStaticError(w, http.StatusInternalServerError, "Internal Server Error")
}

//--------------------------------------------------------------------------------------|

// NotFound sends a 404 response and serves the main index.html file.
func NotFound(w http.ResponseWriter, r *http.Request) {
	log.Printf("[client error] %s %s: 404 Not Found", r.Method, r.URL.Path)
	w.WriteHeader(http.StatusNotFound)
	http.ServeFile(w, r, filepath.Join(defaultStaticDir, "index.html"))
}

func renderStaticError(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(code)
	fmt.Fprintf(w, "Error %d: %s", code, message)
}

//--------------------------------------------------------------------------------------|

// GetLimitOffset extracts and validates pagination parameters from the request.
func GetLimitOffset(r *http.Request) (int, int) {
	limit := QueryInt(r, "limit", DefaultLimit)
	if limit > MaxLimit {
		limit = MaxLimit
	}
	if limit <= 0 {
		limit = DefaultLimit
	}

	offset := QueryInt(r, "offset", 0)
	if offset < 0 {
		offset = 0
	}

	return limit, offset
}

//--------------------------------------------------------------------------------------|

// ExtractIDFromPath parses an integer ID from a specific path segment.
func ExtractIDFromPath(path string, segmentIndex int) (int, error) {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) <= segmentIndex {
		return 0, errors.New("invalid URL structure")
	}

	id, err := utils.ParseInt(parts[segmentIndex])
	if err != nil || id <= 0 {
		return 0, errors.New("invalid or missing ID in path")
	}

	return id, nil
}
