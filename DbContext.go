package uhttp

import (
	"context"
	"net/http"

	"unverricht.net/mongo"
)

// CtxKeyDB is the context key to retrieve the db-info
const CtxKeyDB = ContextKey("database")

// WithDB attaches a dbSession object to the http-request context
func WithDB(session *mongo.DbSession) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sessionCopy := session.Copy()
			ctx := context.WithValue(r.Context(), CtxKeyDB, sessionCopy)
			next.ServeHTTP(w, r.WithContext(ctx))
			sessionCopy.Close()
		}
	}
}
