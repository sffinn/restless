package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/example/go-rest-api/internal/datastore"
)

func newTestServer() http.Handler {
	return New(datastore.NewStore())
}

func TestHealthAndInfo(t *testing.T) {
	handler := newTestServer()

	tests := []struct {
		name string
		path string
		want map[string]any
	}{
		{
			name: "health",
			path: "/health",
			want: map[string]any{"status": "ok"},
		},
		{
			name: "info",
			path: "/",
			want: map[string]any{
				"message": "Go REST API example for DigitalOcean App Platform",
				"endpoints": []any{
					"GET  /health",
					"GET  /api/items",
					"POST /api/items",
					"GET  /api/items/{id}",
					"PUT  /api/items/{id}",
					"DELETE /api/items/{id}",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := request(t, handler, http.MethodGet, tt.path, "")
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
			}

			var got map[string]any
			decodeJSON(t, response, &got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("body = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestItemsCRUD(t *testing.T) {
	handler := newTestServer()

	createResponse := request(t, handler, http.MethodPost, "/api/items", `{"title":"Buy milk"}`)
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", createResponse.Code, http.StatusCreated)
	}
	var created datastore.Item
	decodeJSON(t, createResponse, &created)
	if created.Title != "Buy milk" || created.Completed {
		t.Fatalf("created item = %#v", created)
	}

	listResponse := request(t, handler, http.MethodGet, "/api/items", "")
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listResponse.Code, http.StatusOK)
	}
	var items []datastore.Item
	decodeJSON(t, listResponse, &items)
	if len(items) != 1 || items[0].ID != created.ID {
		t.Fatalf("listed items = %#v", items)
	}

	getResponse := request(t, handler, http.MethodGet, "/api/items/"+created.ID, "")
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", getResponse.Code, http.StatusOK)
	}

	updateResponse := request(t, handler, http.MethodPut, "/api/items/"+created.ID, `{"title":"Buy oat milk","completed":true}`)
	if updateResponse.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updateResponse.Code, http.StatusOK)
	}
	var updated datastore.Item
	decodeJSON(t, updateResponse, &updated)
	if updated.Title != "Buy oat milk" || !updated.Completed || updated.ID != created.ID {
		t.Fatalf("updated item = %#v", updated)
	}

	deleteResponse := request(t, handler, http.MethodDelete, "/api/items/"+created.ID, "")
	if deleteResponse.Code != http.StatusNoContent {
		t.Fatalf("delete status = %d, want %d", deleteResponse.Code, http.StatusNoContent)
	}
	getDeletedResponse := request(t, handler, http.MethodGet, "/api/items/"+created.ID, "")
	if getDeletedResponse.Code != http.StatusNotFound {
		t.Fatalf("get deleted status = %d, want %d", getDeletedResponse.Code, http.StatusNotFound)
	}
}

func TestItemValidationAndMissingResources(t *testing.T) {
	handler := newTestServer()

	tests := []struct {
		name   string
		method string
		path   string
		body   string
		status int
		error  string
	}{
		{
			name:   "create invalid JSON",
			method: http.MethodPost,
			path:   "/api/items",
			body:   `{`,
			status: http.StatusBadRequest,
			error:  "title is required",
		},
		{
			name:   "create missing title",
			method: http.MethodPost,
			path:   "/api/items",
			body:   `{}`,
			status: http.StatusBadRequest,
			error:  "title is required",
		},
		{
			name:   "update invalid JSON",
			method: http.MethodPut,
			path:   "/api/items/missing",
			body:   `{`,
			status: http.StatusBadRequest,
			error:  "invalid JSON",
		},
		{
			name:   "get missing item",
			method: http.MethodGet,
			path:   "/api/items/missing",
			status: http.StatusNotFound,
			error:  "item not found",
		},
		{
			name:   "update missing item",
			method: http.MethodPut,
			path:   "/api/items/missing",
			body:   `{"title":"Missing"}`,
			status: http.StatusNotFound,
			error:  "item not found",
		},
		{
			name:   "delete missing item",
			method: http.MethodDelete,
			path:   "/api/items/missing",
			status: http.StatusNotFound,
			error:  "item not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := request(t, handler, tt.method, tt.path, tt.body)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d", response.Code, tt.status)
			}
			if !strings.Contains(response.Body.String(), tt.error) {
				t.Fatalf("body = %q, want error containing %q", response.Body.String(), tt.error)
			}
		})
	}
}

func request(t *testing.T, handler http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	request := httptest.NewRequest(method, path, strings.NewReader(body))
	response := httptest.NewRecorder()
	handler.ServeHTTP(response, request)
	return response
}

func decodeJSON(t *testing.T, response *httptest.ResponseRecorder, target any) {
	t.Helper()
	if err := json.NewDecoder(response.Body).Decode(target); err != nil {
		t.Fatalf("decode response: %v", err)
	}
}
