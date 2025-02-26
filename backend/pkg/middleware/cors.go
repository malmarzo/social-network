package middleware

import ("net/http"
//"social-network/pkg/utils"
//datamodels "social-network/pkg/dataModels"
"social-network/pkg/db/queries"

)

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


func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cookie, err1 := r.Cookie("session_id")
		if err1 != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login
			return
		}
		
		userID, err2:= queries.ValidateSession(cookie.Value)
		if err2 != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther) // Redirect to login
			return
		}
		r.Header.Set("User-ID", userID)
		next(w, r)
	}
}

