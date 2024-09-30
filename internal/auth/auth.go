package auth

import (
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/google"
	"log"
	"os"
)

const (
	key    = "randomString"
	maxAge = 86400 * 30
	isProd = false
)

// todo - fix bugs + refactor

func NewAuth() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	GoogleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	GoogleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")

	store := sessions.NewCookieStore([]byte(key))
	store.MaxAge(maxAge)

	store.Options.Path = "/"
	store.Options.HttpOnly = true
	store.Options.Secure = isProd

	gothic.Store = store

	goth.UseProviders(
		google.New(GoogleClientId, GoogleClientSecret, "http://localhost:8080/auth/google/callback "),
	)
}
