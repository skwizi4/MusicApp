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
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
}

type Server struct {
	port   int
	logger logs.GoLogger
	db     server_database.Service
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
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", newServer.port),
		Handler:      newServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	fmt.Println(server.Addr)
	return server
}
