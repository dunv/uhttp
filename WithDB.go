package uhttp

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// WithDB attaches a dbSession object to the http-request context
func WithDB(dbName ContextKey, mongoClient *mongo.Client) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			err := mongoClient.Ping(ctx, readpref.Primary())
			if err != nil {
				RenderErrorWithStatusCode(w, r, http.StatusInternalServerError, fmt.Errorf("Could not connect to db"))
			}

			httpContext := context.WithValue(r.Context(), dbName, mongoClient)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}
