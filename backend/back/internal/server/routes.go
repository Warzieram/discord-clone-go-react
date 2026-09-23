package server

import (
	"back/internal/handlers"
	"back/internal/middleware"
	messagerepository "back/internal/repositories/message_repository"
	"net/http"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	// Apply CORS middleware
	r.Use(s.corsMiddleware)

	// Compose dependencies: build the message repository on the injected DB
	// handle, then hand it to the message handlers.
	messageRepo := messagerepository.NewPostgresMessageRepository(s.db.Conn())
	messageHandlers := handlers.NewMessageHandlers(messageRepo)

	// Public Routes
	r.HandleFunc("/api/register", handlers.Register).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/login", handlers.Login).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/verify", handlers.VerifyEmail).Methods("GET", "OPTIONS")

	// Protected Routes
	r.HandleFunc("/api/profile", middleware.AuthMiddleware(handlers.Profile)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/messages", middleware.AuthMiddleware(messageHandlers.RetrieveMessages)).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/message", middleware.WSAuthMiddleware(messageHandlers.MessageHandler)).Methods("GET", "OPTIONS")
	go handlers.SendMessage()

	return r
}

// CORS middleware
func (s *Server) corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// CORS Headers
		w.Header().Set("Access-Control-Allow-Origin", "*") // Wildcard allows all origins
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS, PATCH")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type")
		w.Header().Set("Access-Control-Allow-Credentials", "false") // Credentials not allowed with wildcard origins

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return

		}

		next.ServeHTTP(w, r)
	})
}
