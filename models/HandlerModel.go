package models

import (
	"context"
	"net/http"

	"github.com/dunv/uhttp/params"
)

// Handler configured
type Handler struct {
	Pattern        string
	PostHandler    http.HandlerFunc
	PostModel      interface{}
	GetHandler     http.HandlerFunc
	GetModel       interface{}
	DeleteHandler  http.HandlerFunc
	DeleteModel    interface{}
	RequiredGet    params.R
	OptionalGet    params.R
	AuthRequired   bool
	AuthMiddleware *Middleware
	PreProcess     func(ctx context.Context) error
}
