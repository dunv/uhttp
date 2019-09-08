package models

import (
	"context"
	"net/http"
)

// Handler configured
type Handler struct {
	Pattern                   string
	PostHandler               http.HandlerFunc
	PostModel                 interface{}
	GetHandler                http.HandlerFunc
	GetModel                  interface{}
	DeleteHandler             http.HandlerFunc
	DeleteModel               interface{}
	RequiredParams            Params
	OptionalParams            Params
	AdditionalContextRequired []ContextKey
	AuthRequired              bool
	AuthMiddleware            *Middleware
	PreProcess                func(ctx context.Context) error
}
