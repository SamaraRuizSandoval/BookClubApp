package app

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Config represents the server configuration
type Config struct {
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

// DefaultConfig returns default configurations
func DefaultConfig() *Config {
	return &Config{
		Port:         ":8080",
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
}

// Server represents the HTTP server
type Server struct {
	config *Config
	router *chi.Mux
}

// NewServer creates a new server instance
func NewServer(config *Config) *Server {
	if config == nil {
		config = DefaultConfig()
	}

	server := &Server{
		config: config,
		router: chi.NewRouter(),
	}

	server.setupMiddlewares()
	server.setupRoutes()

	return server
}

// setupMiddlewares configures server middlewares
func (s *Server) setupMiddlewares() {
	// Basic middlewares
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// Configure CORS
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
}

// setupRoutes configures server routes
func (s *Server) setupRoutes() {
	// Root route
	s.router.Get("/", s.handleRoot)
	
	// Health check route
	s.router.Get("/health", s.handleHealth)

	// API routes group
	s.router.Route("/api", func(r chi.Router) {
		// Books routes
		r.Route("/books", func(r chi.Router) {
			r.Get("/", s.handleGetBooks)
			r.Get("/{id}", s.handleGetBookByID)
			r.Post("/", s.handleCreateBook)
			r.Put("/{id}", s.handleUpdateBook)
			r.Delete("/{id}", s.handleDeleteBook)
			
			// Book chapters routes
			r.Route("/{bookId}/chapters", func(r chi.Router) {
				r.Get("/", s.handleGetChapters)
				r.Get("/{chapterId}", s.handleGetChapterByID)
				r.Post("/", s.handleCreateChapter)
				r.Put("/{chapterId}", s.handleUpdateChapter)
				r.Delete("/{chapterId}", s.handleDeleteChapter)
			})
		})

		// Users routes
		r.Route("/users", func(r chi.Router) {
			r.Get("/", s.handleGetUsers)
			r.Get("/{id}", s.handleGetUserByID)
			r.Post("/", s.handleCreateUser)
			r.Put("/{id}", s.handleUpdateUser)
			r.Delete("/{id}", s.handleDeleteUser)
		})

		// Comments routes
		r.Route("/comments", func(r chi.Router) {
			r.Get("/", s.handleGetComments)
			r.Get("/{id}", s.handleGetCommentByID)
			r.Post("/", s.handleCreateComment)
			r.Put("/{id}", s.handleUpdateComment)
			r.Delete("/{id}", s.handleDeleteComment)
			
			// Comments by chapter
			r.Route("/chapter/{chapterId}", func(r chi.Router) {
				r.Get("/", s.handleGetCommentsByChapter)
			})
		})
	})
}

// Start starts the HTTP server
func (s *Server) Start() error {
	log.Printf("🚀 HTTP Server starting on port %s", s.config.Port)
	log.Printf("📚 BookClubApp API available at http://localhost%s", s.config.Port)
	log.Printf("🔍 Health check: http://localhost%s/health", s.config.Port)
	log.Printf("📖 API Books: http://localhost%s/api/books", s.config.Port)
	log.Printf("📑 API Chapters: http://localhost%s/api/books/{id}/chapters", s.config.Port)
	log.Printf("💬 API Comments: http://localhost%s/api/comments", s.config.Port)
	log.Printf("👥 API Users: http://localhost%s/api/users", s.config.Port)

	server := &http.Server{
		Addr:         s.config.Port,
		Handler:      s.router,
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
	}

	return server.ListenAndServe()
}

// Handlers
func (s *Server) handleRoot(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "BookClubApp API Server", "status": "running", "version": "1.0.0", "features": ["books", "chapters", "comments", "users"]}`))
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "healthy", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`))
}

// Book handlers
func (s *Server) handleGetBooks(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Books list", "data": [], "count": 0}`))
}

func (s *Server) handleGetBookByID(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Book details", "id": "` + bookID + `", "data": {}}`))
}

func (s *Server) handleCreateBook(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Book created successfully", "data": {}}`))
}

func (s *Server) handleUpdateBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Book updated successfully", "id": "` + bookID + `"}`))
}

func (s *Server) handleDeleteBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Book deleted successfully", "id": "` + bookID + `"}`))
}

// Chapter handlers
func (s *Server) handleGetChapters(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Book chapters list", "bookId": "` + bookID + `", "data": [], "count": 0}`))
}

func (s *Server) handleGetChapterByID(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	chapterID := chi.URLParam(r, "chapterId")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Chapter details", "bookId": "` + bookID + `", "chapterId": "` + chapterID + `", "data": {}}`))
}

func (s *Server) handleCreateChapter(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Chapter created successfully", "bookId": "` + bookID + `", "data": {}}`))
}

func (s *Server) handleUpdateChapter(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	chapterID := chi.URLParam(r, "chapterId")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Chapter updated successfully", "bookId": "` + bookID + `", "chapterId": "` + chapterID + `"}`))
}

func (s *Server) handleDeleteChapter(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookId")
	chapterID := chi.URLParam(r, "chapterId")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Chapter deleted successfully", "bookId": "` + bookID + `", "chapterId": "` + chapterID + `"}`))
}

// User handlers
func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Users list", "data": [], "count": 0}`))
}

func (s *Server) handleGetUserByID(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User details", "id": "` + userID + `", "data": {}}`))
}

func (s *Server) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "User created successfully", "data": {}}`))
}

func (s *Server) handleUpdateUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User updated successfully", "id": "` + userID + `"}`))
}

func (s *Server) handleDeleteUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "User deleted successfully", "id": "` + userID + `"}`))
}

// Comment handlers
func (s *Server) handleGetComments(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Comments list", "data": [], "count": 0}`))
}

func (s *Server) handleGetCommentByID(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Comment details", "id": "` + commentID + `", "data": {}}`))
}

func (s *Server) handleCreateComment(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"message": "Comment created successfully", "data": {}}`))
}

func (s *Server) handleUpdateComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Comment updated successfully", "id": "` + commentID + `"}`))
}

func (s *Server) handleDeleteComment(w http.ResponseWriter, r *http.Request) {
	commentID := chi.URLParam(r, "id")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Comment deleted successfully", "id": "` + commentID + `"}`))
}

func (s *Server) handleGetCommentsByChapter(w http.ResponseWriter, r *http.Request) {
	chapterID := chi.URLParam(r, "chapterId")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Chapter comments", "chapterId": "` + chapterID + `", "data": [], "count": 0}`))
}
