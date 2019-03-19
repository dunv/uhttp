package uhttp

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// CtxKeyDB is the context key to retrieve the db-info
const CtxKeyDB = ContextKey("database")

// WithDB attaches a dbSession object to the http-request context
func WithDB(dbName string, mongoClient *mongo.Client) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			ctx, _ := context.WithTimeout(context.Background(), 2*time.Second)
			err := mongoClient.Ping(ctx, readpref.Primary())
			if err != nil {
				js, _ := json.Marshal(Error{"Could not connect to db"})
				w.WriteHeader(http.StatusMethodNotAllowed)
				w.Write(js)
			}

			httpContext := context.WithValue(r.Context(), ContextKey(dbName), mongoClient)
			next.ServeHTTP(w, r.WithContext(httpContext))
		}
	}
}
