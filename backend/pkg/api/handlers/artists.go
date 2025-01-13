package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/krukkrz/label_explorer/pkg/releases"
)

func (h *Handlers) Artists(w http.ResponseWriter, r *http.Request) {
	artistName := r.URL.Query().Get("artist")
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	validateArtistsSortBy(sortBy, w)
	if ok := validateOrder(order); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	releasesCounts, err := releases.ByArtist(h.Db, artistName, sortBy, order)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(releasesCounts)
}

func validateArtistsSortBy(value string, w http.ResponseWriter) {
	allowedValues := map[string]struct{}{
		"release_count": {},
		"style_genre":   {},
	}

	if _, ok := allowedValues[value]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
