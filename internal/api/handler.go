package api

import (
	"encoding/json"
	"net/http"

	"github.com/XAER/osu-ha-integration/internal/osu"
)

type Handler struct {
	Client *osu.Client
}

func NewHandler(client *osu.Client) *Handler {
	return &Handler{
		Client: client,
	}
}

func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "username is required", http.StatusBadRequest)
		return
	}

	user, err := h.Client.GetUser(r.Context(), username)
	if err != nil {
		http.Error(w, "failed to fetch user stats", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
