package helpers

// // Helper function (ONLY for usage within a middleware, a context needs to be passed on) for adding params to uhttp request-logging
// func AddToLogLine(ctx context.Context, key string, value string) context.Context {
// 	currentRaw := ctx.Value(contextkeys.CtxKeyAddLogging)
// 	var current map[string]string
// 	if currentRaw == nil {
// 		current = map[string]string{}
// 	} else {
// 		current = currentRaw.(map[string]string)
// 	}
// 	current[key] = value
// 	return context.WithValue(ctx, contextkeys.CtxKeyAddLogging, current)
// }

// // Helper function for getting additional logging parameters from the request
// func LogLineParams(r *http.Request) map[string]string {
// 	currentRaw := r.Context().Value(contextkeys.CtxKeyAddLogging)
// 	var current map[string]string
// 	if currentRaw == nil {
// 		current = map[string]string{}
// 	} else {
// 		current = currentRaw.(map[string]string)
// 	}
// 	return current
// }
