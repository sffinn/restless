package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/sffinn/restless/internal/datastore"
)

type API struct {
	store *datastore.Store
}

func New(store *datastore.Store) http.Handler {
	api := &API{store: store}
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", api.health)
	mux.HandleFunc("GET /", api.info)
	mux.HandleFunc("GET /api/items", api.listItems)
	mux.HandleFunc("POST /api/items", api.createItem)
	mux.HandleFunc("GET /api/items/{id}", api.getItem)
	mux.HandleFunc("PUT /api/items/{id}", api.updateItem)
	mux.HandleFunc("DELETE /api/items/{id}", api.deleteItem)

	return requestLogger(mux)
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (r *responseRecorder) WriteHeader(status int) {
	if r.status != 0 {
		return
	}
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func (r *responseRecorder) Write(body []byte) (int, error) {
	if r.status == 0 {
		r.WriteHeader(http.StatusOK)
	}
	n, err := r.ResponseWriter.Write(body)
	r.bytes += n
	return n, err
}

func requestLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		recorder := &responseRecorder{ResponseWriter: w}

		next.ServeHTTP(recorder, r)

		status := recorder.status
		if status == 0 {
			status = http.StatusOK
		}
		log.Printf("request method=%s path=%s status=%d bytes=%d duration=%s",
			r.Method, r.URL.Path, status, recorder.bytes, time.Since(start))
	})
}

func (a *API) health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (a *API) info(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]any{
		"message": "Go REST API example for DigitalOcean App Platform",
		"endpoints": []string{
			"GET  /health",
			"GET  /api/items",
			"POST /api/items",
			"GET  /api/items/{id}",
			"PUT  /api/items/{id}",
			"DELETE /api/items/{id}",
		},
	})
}

func (a *API) listItems(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, a.store.List())
}

func (a *API) createItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil || body.Title == "" {
		http.Error(w, `{"error":"title is required"}`, http.StatusBadRequest)
		return
	}
	writeJSON(w, http.StatusCreated, a.store.Create(body.Title))
}

func (a *API) getItem(w http.ResponseWriter, r *http.Request) {
	item, ok := a.store.Get(r.PathValue("id"))
	if !ok {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) updateItem(w http.ResponseWriter, r *http.Request) {
	var body struct {
		Title     string `json:"title"`
		Completed bool   `json:"completed"`
	}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, `{"error":"invalid JSON"}`, http.StatusBadRequest)
		return
	}
	item, ok := a.store.Update(r.PathValue("id"), body.Title, body.Completed)
	if !ok {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, item)
}

func (a *API) deleteItem(w http.ResponseWriter, r *http.Request) {
	if !a.store.Delete(r.PathValue("id")) {
		http.Error(w, `{"error":"item not found"}`, http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}
