// Package web provides HTTP decorators for cross-cutting concerns like identity injection.
package web

import (
	"errors"
	"social-network/internal/models"
	"net/http"
	"path/filepath"
)

//--------------------------------------------------------------------------------------|

// OnlyMethod returns a StatusError if the request method is not the expected one.

func OnlyMethod(method string, h AppHandler) AppHandler {
	return func(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
		if r.Method != method {
			return StatusError{Code: http.StatusMethodNotAllowed, Err: errors.New("method not allowed")}
		}
		return h(w, r, identity)
	}
}

//--------------------------------------------------------------------------------------|

// SPAOrAPI ensures that non-JSON GET requests return index.html (SPA shell)
// instead of executing the wrapped handler.

func SPAOrAPI(h AppHandler) AppHandler {
	return func(w http.ResponseWriter, r *http.Request, identity *models.UserIdentity) error {
		// If it's a GET request and the client does not explicitly ask for JSON
		if r.Method == http.MethodGet && r.Header.Get("Accept") != "application/json" {
			http.ServeFile(w, r, filepath.Join(defaultStaticDir, "index.html"))
			return nil
		}
		// Otherwise, proceed with the API handler
		return h(w, r, identity)
	}
}
