package news

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"

	"go-news-feed/pkg/model"
)

type endpoint struct {
	service   Service
	validator *validator.Validate
}

// newEndpoint - constructor
func newEndpoint(service Service) *endpoint {
	return &endpoint{
		service:   service,
		validator: validator.New(),
	}
}

func (e *endpoint) init() *http.ServeMux {
	// Initialize HTTP request multiplexer
	mux := http.NewServeMux()

	// Routes
	mux.HandleFunc("GET /find", e.find)
	mux.HandleFunc("GET /load", e.load)

	return mux
}

func (e endpoint) find(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	if err := r.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("failed to parse request body: %v", err), http.StatusBadRequest)
		return
	}

	// Transformation from map[string][]string to map[string]string:
	m := map[string]string{}
	for k, v := range r.Form {
		m[k] = v[0]
	}

	// Marshal request body
	data, err := json.Marshal(m)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal request body: %v", err), http.StatusBadRequest)
		return
	}

	// Decode request body into a new object
	var fr model.FindRequest
	if err := json.Unmarshal(data, &fr); err != nil {
		http.Error(w, fmt.Sprintf("failed to decode request body: %v", err), http.StatusBadRequest)
		return
	}

	// Validate the request
	if err := e.validator.Struct(fr); err != nil {
		http.Error(w, fmt.Sprintf("invalid request: %v", err), http.StatusBadRequest)
		return
	}

	// Find
	response, err := e.service.Find(r.Context(), fr)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to find news: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode object as JSON and write to response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}

func (e endpoint) load(w http.ResponseWriter, r *http.Request) {
	response, err := e.service.Load(r.Context(), r.URL.Query().Get("feedUrl"))
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to load news: %v", err), http.StatusInternalServerError)
		return
	}

	// Encode object as JSON and write to response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, fmt.Sprintf("failed to encode response: %v", err), http.StatusInternalServerError)
		return
	}
}
