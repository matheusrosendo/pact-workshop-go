package service

import (
	"encoding/json"
	"net/http"

	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pact-foundation/pact-workshop-go/provider/internal/model"
	"github.com/pact-foundation/pact-workshop-go/provider/internal/repository"
)

var userRepository = &repository.UserRepository{
	Users: map[string]*model.User{
		"sally": &model.User{
			FirstName: "Jean-Marie",
			LastName:  "de La Beaujardière😀😍",
			Username:  "sally",
			Type:      "admin",
			ID:        10,
		},
	},
}

// Crude time-bound "bearer" token
func getAuthToken() string {
	return fmt.Sprintf("Bearer %s", time.Now().Format("2006-01-02T15:04"))
}

// GetAuthToken returns the current auth token (exported for testing)
func GetAuthToken() string {
	return getAuthToken()
}

// SetUserRepository sets the user repository (exported for testing)
func SetUserRepository(repo *repository.UserRepository) {
	userRepository = repo
}

// IsAuthenticated checks for a correct bearer token
func WithCorrelationID(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uuid := uuid.New()
		w.Header().Set("X-Api-Correlation-Id", uuid.String())
		h.ServeHTTP(w, r)
	}
}

// IsAuthenticated checks for a correct bearer token
func IsAuthenticated(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") == getAuthToken() {
			h.ServeHTTP(w, r)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
		}
	}
}

// GetUser fetches a user if authenticated and exists
func GetUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Get username from path
	a := strings.Split(r.URL.Path, "/")
	id, _ := strconv.Atoi(a[len(a)-1])

	user, err := userRepository.ByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusOK)
		resBody, _ := json.Marshal(user)
		w.Write(resBody)
	}
}

// GetUsers fetches all users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resBody, _ := json.Marshal(userRepository.GetUsers())
	w.Write(resBody)
}

func Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func commonMiddleware(f http.HandlerFunc) http.HandlerFunc {
	return WithCorrelationID(IsAuthenticated(f))
}

func GetHTTPHandler() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/user/", commonMiddleware(GetUser))
	mux.HandleFunc("/users/", commonMiddleware(GetUsers))
	mux.HandleFunc("/healthasdf", Health)

	return mux
}
