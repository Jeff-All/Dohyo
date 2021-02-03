package handlers

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

// HandlerInterface - Interface for all handlers
type HandlerInterface interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

// Handler - Base struct for all handlers
type Handler struct {
	Log *logrus.Logger
}
