package models

import "net/http"

type Middleware func(next http.HandlerFunc) http.HandlerFunc