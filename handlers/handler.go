package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// HandlerInterface - Interface for all handlers
type HandlerInterface interface {
	GetName() string
	GetRoute() string
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Handler - Base struct for all handlers
type Handler struct {
	Name  string
	Route string
	Log   *logrus.Logger
}

// GetName - Returns Handler.Name
func (h Handler) GetName() string { return h.Name }

// GetRoute - Returns Handler.Route
func (h Handler) GetRoute() string { return h.Route }
