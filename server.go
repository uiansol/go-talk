package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
	"gorm.io/gorm"
)

// Contains pointers to the database connection and token auth info.
type Config struct {
	DB        *gorm.DB
	TokenAuth *jwtauth.JWTAuth
}

// Setup and Start the server
func main() {
	var config Config

	// These variable are here for easy demonstration.
	// It's not good to have it hardcoded for security issues.
	port := "8080"
	dsn := "go-talk.db"
	jwtSecret := "go-servers"

	if fromEnv := os.Getenv("PORT"); fromEnv != "" {
		port = fromEnv
	}

	if fromEnv := os.Getenv("DB_DSN"); fromEnv != "" {
		dsn = fromEnv
	}

	if fromEnv := os.Getenv("JWT_SECRET"); fromEnv != "" {
		jwtSecret = fromEnv
	}
	config.TokenAuth = jwtauth.New("HS256", []byte(jwtSecret), nil)

	// Setup log server activity
	f, err := os.OpenFile("server.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(f)

	// Setup data base
	db, err := DBConnect(dsn)
	if err != nil {
		log.Fatalf("db conn failure: %v", err)
	}
	config.DB = db

	err = config.EnsureDBSetup()
	if err != nil {
		log.Fatalf("db setup failure: %v", err)
	}

	// Setup routes with end points
	router := config.SetupRoutes()
	log.Printf("Starting up on http://localhost:%s", port)
	log.Println(http.ListenAndServe(":"+port, router))
}

// Creates the router and define end-points.
func (c *Config) SetupRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	r.Get("/", IndexHandler)
	r.Post("/newuser", c.NewUserPostHandler)
	r.Post("/login", c.LoginPostHandler)
	r.Route("/", func(r chi.Router) {
		r.Use(jwtauth.Verifier(c.TokenAuth))
		r.Use(jwtauth.Authenticator)

		r.Get("/chat", APIChatHandler)
		r.Route("/api", func(r chi.Router) {
			r.Get("/messages", c.APIMessagesHandler)
			r.Post("/messages", c.APIMessagesPostHandler)
		})
	})

	return r
}
