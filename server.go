package main

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/jwtauth"
)

func main() {

}

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

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}
