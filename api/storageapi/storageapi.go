package storageapi

import (
	"bytes"
	"io"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spacelift-io/homework-object-storage/storage"
)

func RegisterRoutes(r chi.Router, s storage.Storage) {
	r.Route("/object/{id}", func(r chi.Router) {
		r.Get("/", Get(s))
		r.Put("/", Put(s))
	})
}

func Get(s storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		reader, err := s.Get(r.Context(), id)
		if err != nil {
			log.Printf("error getting from storage: %s\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}
		defer reader.Close()

		bytes, err := io.ReadAll(reader)
		if err != nil {
			log.Printf("error reading from storage: %s\n", err)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/octet-stream")
		written, err := w.Write(bytes)
		if err != nil || written != len(bytes) {
			log.Printf("writing error: %s; written: %d, should write: %d\n", err, written, len(bytes))
		}
	}
}

func Put(s storage.Storage) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")

		if r.Header.Get("Content-Type") != "application/octet-stream" {
			log.Printf("wrong content-type\n")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		data, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("reading body data error: %s\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		if err := s.Put(r.Context(), id, bytes.NewReader(data), int64(len(data))); err != nil {
			log.Printf("error writing to the storage: %s\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}
