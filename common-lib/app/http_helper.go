package app

import (
	"fmt"
	"net/http"
)

func readyz() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		if cfg.isReady == nil || !cfg.isReady.Load().(bool) {
			http.Error(w, http.StatusText(http.StatusServiceUnavailable), http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func handle404() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		HttpError(
			w, http.StatusNotFound,
			http.StatusText(http.StatusNotFound),
			fmt.Errorf("path: %s", r.URL))
	})
}

func method404() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		HttpError(
			w, http.StatusNotFound,
			http.StatusText(http.StatusMethodNotAllowed),
			fmt.Errorf("path: %s", r.URL))
	})
}

func setResponseHeader(w http.ResponseWriter) {
	w.Header().Set(string(HTTP_CONTENT_TYPE), string(HTTP_CT_JSON))
	w.Header().Set(string(HTTP_SERVER), cfg.name)
}
