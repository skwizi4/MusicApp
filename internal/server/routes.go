package server

import (
	"encoding/json"
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

	r.HandleFunc("/auth/google/callback", s.getAuthCallBackFunction)
	return r
}

// todo -refactor + fix bugs
func (s *Server) getAuthCallBackFunction(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	TelegramId := r.URL.Query().Get("state")
	if code == "" || TelegramId == "" {
		http.Error(w, "Invalid request, missing code or TelegramID", http.StatusBadRequest)
		return
	}

	s.logger.InfoFrmt("Code: %s\nTelegramID: %s", code, TelegramId)

	// Создание данных для запроса токена
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

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		s.logger.ErrorFrmt("Request failed: %v", err)
		http.Error(w, "Failed to get token from Google", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	s.logger.InfoFrmt("Response status: %s", resp.Status)

	// Чтение тела ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		s.logger.ErrorFrmt("Failed to read response body: %v", err)
		http.Error(w, "Failed to read response from Google", http.StatusInternalServerError)
		return
	}

	s.logger.InfoFrmt("Response body: %s", string(body))
	var result TokenResponse
	if err = json.Unmarshal(body, &result); err != nil {
		log.Fatal(err)
	}
	s.logger.InfoFrmt("Response result: %s", result)
	err = s.db.Create(result.AccessToken, TelegramId)
	if err != nil {
		s.logger.ErrorFrmt("Failed to pull token into db : %v", err)
		return
	}

	http.Redirect(w, r, "http://localhost:5317", http.StatusFound)
}
