package server

import (
	"fmt"
	"github.com/skwizi4/lib/logs"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"MusicApp/internal/server_database"
)

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type Server struct {
	port   int
	logger logs.GoLogger
	db     server_database.Service
	Client *http.Client
}

func NewServer() *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	db, err := server_database.New()
	if err != nil {
		panic(err)
	}
	newServer := &Server{
		port:   port,
		logger: logs.InitLogger(),
		db:     db,
		Client: &http.Client{},
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	return server
}
