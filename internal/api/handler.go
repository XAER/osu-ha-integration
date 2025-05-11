package api

import (
	"encoding/json"
	"net/http"

	"github.com/XAER/osu-ha-integration/internal/osu"
)

type Handler struct {
	Client *osu.Client
	Cache  *osu.Cache
}

func NewHandler(client *osu.Client, cache *osu.Cache) *Handler {
	return &Handler{Client: client, Cache: cache}
}

func (h *Handler) GetUserStats(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "missing username", http.StatusBadRequest)
		return
	}

	if cached, ok := h.Cache.Get(username); ok {
		json.NewEncoder(w).Encode(cached)
		return
	}

	user, err := h.Client.GetUser(r.Context(), username)
	if err != nil {
		http.Error(w, "failed to fetch user", http.StatusInternalServerError)
		return
	}

	h.Cache.Set(username, user)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
