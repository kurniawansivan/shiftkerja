package handler

import "net/http"

// CORSMiddleware handles the browser's "Preflight" checks
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Allow connections from anywhere (for now)
		w.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 2. Allow specific methods
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		
		// 3. Allow specific headers (Authorization is critical for JWT)
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")

		// 4. Handle Preflight (The browser asking for permission)
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// 5. Pass the request to the real handler
		next.ServeHTTP(w, r)
	})
}