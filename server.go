package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

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

func (c *Config) APIMessagesPostHandler(w http.ResponseWriter, r *http.Request) {
	var message Message
	var redirect bool

	if r.Header.Get("Content-type") == "application/json" {
		err := json.NewDecoder(r.Body).Decode(&message)
		if err != nil {
			http.Error(w, "malformed request body", http.StatusBadRequest)
			return
		}
	} else {
		r.ParseForm()
		messageText := r.PostForm.Get("message")
		if messageText == "" {
			http.Error(w, "malformed request body", http.StatusBadRequest)
			return
		}
		message.Text = messageText
		redirect = true
	}

	userName := GetUserNameFromContext(r.Context())
	if _, err := c.CreateMessage(userName, message.Text); err != nil {
		log.Printf("db create message failure: %q", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if redirect {
		http.Redirect(w, r, "/chat", http.StatusSeeOther)
	} else {
		w.WriteHeader(http.StatusCreated)
	}
}

func (c *Config) APIMessagesHandler(w http.ResponseWriter, r *http.Request) {
	var (
		err        error
		timeParam  string
		parsedTime int64  = time.Now().UnixNano() / 1e6
		order      ByTime = older
		attempts   int    = 1
	)

	r.ParseForm()
	timeParam = r.Form.Get("before")
	since := r.Form.Get("since")

	if since != "" {
		attempts = 10
		order = newer
		timeParam = since
	}

	if timeParam != "" {
		if parsedTime, err = strconv.ParseInt(timeParam, 10, 64); err != nil {
			http.Error(w, "invalid 'before' or 'since' parameter sent", http.StatusBadRequest)
		}
	}

	var messages []Message
	for i := 0; i < attempts; i++ {
		messages, err = c.GetMessagesByTime(10, parsedTime, order)
		if err != nil {
			log.Printf("db failed to get messages: %q", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if len(messages) != 0 {
			break
		}
		time.Sleep(1 * time.Second)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	body, _ := json.Marshal(messages)
	w.Write(body)
}
