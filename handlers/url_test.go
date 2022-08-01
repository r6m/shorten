package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/r6m/shorten/models"
	"github.com/r6m/shorten/store/memory"
)

func TestShortenHandler(t *testing.T) {
	server := NewServer(memory.NewStore())

	t.Run("should save valid url", func(t *testing.T) {
		req := &shortenURLRequest{
			URL: "https://google.com",
		}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(req)

		r := httptest.NewRequest(http.MethodPost, "/shorten", buf)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(server.shortenHandler)(w, r)

		res := w.Result()
		defer res.Body.Close()

		url := new(models.URL)

		err := json.NewDecoder(res.Body).Decode(url)
		if err != nil {
			t.Fatalf("can't parse json: %v", err)
		}

		if res.StatusCode != http.StatusOK {
			t.Fatalf("should return status 200, but got %d", res.StatusCode)
		}

		if url.Key == "" {
			t.Fatalf("should return a url key, but got: %s", url.Key)
		}

		if url.OriginalURL != req.URL {
			t.Fatalf("original url is not the same: %s != %s", url.OriginalURL, req.URL)
		}
	})

	t.Run("should save predefined key", func(t *testing.T) {
		req := &shortenURLRequest{
			URL: "https://google.com",
			Key: "111111",
		}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(req)

		r := httptest.NewRequest(http.MethodPost, "/shorten", buf)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(server.shortenHandler)(w, r)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			t.Fatalf("should return status 200, but got %d", res.StatusCode)
		}

		url := new(models.URL)

		err := json.NewDecoder(res.Body).Decode(url)
		if err != nil {
			t.Fatalf("can't parse json: %v", err)
		}

		if url.Key != req.Key {
			t.Fatalf("should return predefined key, but got: %s", url.Key)
		}

		if url.OriginalURL != req.URL {
			t.Fatalf("original url is not the same: %s != %s", url.OriginalURL, req.URL)
		}
	})

	t.Run("should fail on duplicate", func(t *testing.T) {
		req := &shortenURLRequest{
			URL: "https://google.com",
			Key: "111111",
		}

		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(req)

		r := httptest.NewRequest(http.MethodPost, "/shorten", buf)
		r.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		handler(server.shortenHandler)(w, r)

		res := w.Result()
		defer res.Body.Close()

		if res.StatusCode != http.StatusBadRequest {
			t.Fatalf("should return status 400, but got %d", res.StatusCode)
		}

		errResp := new(HTTPError)

		err := json.NewDecoder(res.Body).Decode(errResp)
		if err != nil {
			t.Fatalf("can't parse json: %v", err)
		}

		if !strings.Contains(errResp.Message, "already exists") {
			t.Fatalf("should return already exists, but got '%s'", errResp.Message)
		}
	})

}

func TestRedirectHandler(t *testing.T) {
	store := memory.NewStore()

	url := &models.URL{Key: "111111", OriginalURL: "https://google.com"}
	store.Save(url)

	server := NewServer(store)

	t.Run("should handle redirect", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/{key}", nil)
		w := httptest.NewRecorder()
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("key", url.Key)
		r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

		handler(server.redirectHandler)(w, r)

		res := w.Result()

		location := res.Header.Get("Location")
		if location != url.OriginalURL {
			t.Fatalf("should set redirect header to %s, but got %s", url.OriginalURL, location)
		}

		if res.StatusCode != http.StatusMovedPermanently {
			t.Fatalf("should set status to %d, but got %d", http.StatusMovedPermanently, res.StatusCode)
		}
	})
}
