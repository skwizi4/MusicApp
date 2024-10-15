package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/auth/google/callback", s.getGoogleAuthCallBackFunction)
	r.HandleFunc("/auth/spotify/callback", s.getSpotifyAuthCallBackFunction)
	r.HandleFunc("/authentication_completed", s.AuthCompleted)
	return r
}

// todo -refactor + fix bugs
func (s *Server) getGoogleAuthCallBackFunction(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	TelegramId := r.URL.Query().Get("state")
	if code == "" || TelegramId == "" {
		http.Error(w, "Invalid request, missing code or UserProcess", http.StatusBadRequest)
		return
	}

	data := url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("GOOGLE_CLIENT_ID")},
		"client_secret": {os.Getenv("GOOGLE_CLIENT_SECRET")},
		"redirect_uri":  {"http://localhost:8080/auth/google/callback"},
		"grant_type":    {"authorization_code"},
	}

	// Запрос токена
	req, err := http.NewRequest("POST", "https://oauth2.googleapis.com/token", strings.NewReader(data.Encode()))
	if err != nil {
		s.logger.ErrorFrmt("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.Client.Do(req)
	if err != nil {
		s.logger.ErrorFrmt("Request failed: %v", err)
		http.Error(w, "Failed to get token from Google", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.ErrorFrmt("Failed to read response body: %v", err)
		http.Error(w, "Failed to read response from Google", http.StatusInternalServerError)
		return
	}

	var result TokenResponse
	if err = json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}
	UserProcess := fmt.Sprintf("YoutubeProcess%s", TelegramId)
	err = s.db.Add(UserProcess, result.AccessToken, result.RefreshToken)
	if err != nil {
		s.logger.ErrorFrmt("Failed to pull token into db : %v", err)
		return
	}

	http.Redirect(w, r, "http://localhost:8080/authentication_completed", http.StatusFound)
}

func (s *Server) getSpotifyAuthCallBackFunction(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	TelegramId := r.URL.Query().Get("state")
	if code == "" || TelegramId == "" {
		err := r.URL.Query().Get("error")
		if err != "" {
			s.logger.ErrorFrmt("Authorization failed: %v", err)
			return
		}
		http.Error(w, "Invalid request, missing code or UserProcess", http.StatusBadRequest)
		return
	}
	data := url.Values{
		"code":          {code},
		"client_id":     {os.Getenv("SPOTIFY_CLIENT_ID")},
		"client_secret": {os.Getenv("SPOTIFY_CLIENT_SECRET")},
		"redirect_uri":  {"http://localhost:8080/auth/spotify/callback"},
		"grant_type":    {"authorization_code"},
	}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", bytes.NewBufferString(data.Encode()))
	if err != nil {
		s.logger.ErrorFrmt("Failed to create request: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := s.Client.Do(req)
	if resp.StatusCode != 200 {
		http.Error(w, resp.Status, resp.StatusCode)
		return
	}
	if err != nil {
		http.Error(w, "err, try later ", http.StatusBadRequest)
		return
	}
	var result TokenResponse
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		s.logger.ErrorFrmt("Failed to decode response body: %v", err)
		return
	}

	UserProcess := fmt.Sprintf("SpotifyProcess%s", TelegramId)
	err = s.db.Add(UserProcess, result.AccessToken, result.RefreshToken)
	if err != nil {
		logrus.Error(err)
	}
	http.Redirect(w, r, "http://localhost:8080/authentication_completed", http.StatusFound)
}

func (s *Server) AuthCompleted(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You have been authenticated")
}
