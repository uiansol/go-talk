package main

import (
	"log"
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

func (c *Config) NewUserPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.PostForm.Get("user")
	userPassword := r.PostForm.Get("password")
	userPasswordConfirm := r.PostForm.Get("password_confirm")

	if userName == "" || userPassword == "" {
		http.Error(w, "missing user or password", http.StatusBadRequest)
		return
	}

	if userPassword != userPasswordConfirm {
		http.Error(w, "passwords do not match", http.StatusBadRequest)
		return
	}

	if _, err := c.CreateUser(userName, userPassword); err != nil {
		log.Printf("database query error: %q\n", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (c *Config) LoginPostHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userName := r.PostForm.Get("user")
	userPassword := r.PostForm.Get("password")

	if userName == "" || userPassword == "" {
		http.Error(w, "missing user or password", http.StatusBadRequest)
		return
	}

	if err := c.CheckLogin(userName, userPassword); err != nil {
		http.Error(w, "login unsuccessful", http.StatusBadRequest)
		return
	}

	token := c.MakeToken(userName)
	http.SetCookie(w, &http.Cookie{
		Name:  "jwt",
		Value: token,
	})
	http.Redirect(w, r, "/chat", http.StatusSeeOther)
}
