package uhttp

import (
	"context"
	"net/http"
)

// CtxKeyDB is the context key to retrieve the db-info
const CtxKeyDB = ContextKey("database")

// WithDB attaches a dbSession object to the http-request context
func WithDB(session *uumongo.DbSession) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			sessionCopy := session.Copy()
			ctx := context.WithValue(r.Context(), CtxKeyDB, sessionCopy)
			next.ServeHTTP(w, r.WithContext(ctx))
			sessionCopy.Close()
		}
	}
}
