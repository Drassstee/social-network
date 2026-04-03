// Package web provides HTTP constants and configuration defaults for the web layer.
package web

//--------------------------------------------------------------------------------------|

// Cookie and route constants used throughout the web package.
const (
	CookieSessionID = "session_id"

	// Routes
	RoutePosts         = "/posts"
	RouteCategories    = "/categories"
	RouteNotifications = "/notifications"
	RouteLogin         = "/login"
	RouteRegister      = "/register"
)
