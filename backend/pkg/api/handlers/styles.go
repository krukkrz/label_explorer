package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/krukkrz/label_explorer/pkg/releases"
)

func (h *Handlers) Styles(w http.ResponseWriter, r *http.Request) {
	styleGenreName := r.URL.Query().Get("style")
	sortBy := r.URL.Query().Get("sort")
	order := r.URL.Query().Get("order")

	validateStyleSortBy(sortBy, w)
	if ok := validateOrder(order); !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	releasesCounts, err := releases.ByStyleGenre(h.Db, styleGenreName, sortBy, order)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal error"))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(releasesCounts)
}

func validateStyleSortBy(value string, w http.ResponseWriter) {
	allowedValues := map[string]struct{}{
		"release_count": {},
		"artist_name":   {},
	}

	if _, ok := allowedValues[value]; !ok {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
