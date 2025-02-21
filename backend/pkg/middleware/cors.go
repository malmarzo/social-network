package middleware

import "net/http"

// CORS Middleware
func CorsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// enableCors(w)
		// Handle preflight (OPTIONS) request for CORS
		if r.Method == http.MethodOptions {
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS") // Allow HTTP methods
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// func enableCors(w http.ResponseWriter) {
// 	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5005")
// 	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
// 	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
// 	w.Header().Set("Access-Control-Allow-Credentials", "true")
// }
