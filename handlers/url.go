package handlers

import (
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/r6m/shorten/models"
	"github.com/r6m/shorten/store"
	"github.com/sirupsen/logrus"
)

type shortenURLRequest struct {
	URL string `json:"url" form:"url"`
	Key string `json:"key" form:"key"`
}

func (api *API) shortenHandler(w http.ResponseWriter, r *http.Request) any {
	in := new(shortenURLRequest)
	if err := render.Decode(r, in); err != nil {
		return badRequestError("can't parse body: %v", err)
	}

	_, err := url.Parse(in.URL)
	if err != nil {
		return badRequestError("invalid url: %v", err)
	}

	url := &models.URL{
		Key:         in.Key,
		OriginalURL: in.URL,
		Details:     make([]*models.Detail, 0),
	}

	isUserKey := in.Key != ""
	if url.Key == "" {
		url.GenerateKey()
	}

	retry := 3
	for i := 0; i < retry; i++ {
	REDO:
		if err := api.repo.Save(url); err != nil {
			if errors.Is(err, store.ErrDuplicate) {
				if isUserKey {
					return badRequestError("key '%s' already exists", in.Key)
				}
				url.GenerateKey()
				goto REDO
			}
			return internalError("can't save url").withInternalError(err)
		}

		break
	}

	logrus.Printf("save key: '%s', url: '%s'", url.Key, url.OriginalURL)
	return url
}

func (api *API) redirectHandler(w http.ResponseWriter, r *http.Request) any {
	key := chi.URLParam(r, "key")

	ua := r.Header.Get("User-Agent")

	detail := &models.Detail{
		UserAgent: ua,
		CreatedAt: time.Now(),
	}

	url, err := api.repo.Load(key, detail)
	if err != nil {
		logrus.Printf("load: key '%s' not found", key)
		return notFoundError("key not found")
	}
	logrus.Printf("redirect: '%s', url: '%s'", url.Key, url.OriginalURL)

	http.Redirect(w, r, url.OriginalURL, http.StatusTemporaryRedirect)
	return nil
}

func (api *API) infoHandler(w http.ResponseWriter, r *http.Request) any {
	key := chi.URLParam(r, "key")

	url, err := api.repo.Load(key, nil)
	if err != nil {
		return notFoundError("url not found")
	}

	return url.Details
}
