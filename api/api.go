package api

import (
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/spacelift-io/homework-object-storage/api/storageapi"
	"github.com/spacelift-io/homework-object-storage/storage"
)

type API struct {
	router       chi.Router
	server       *http.Server
	serverClosed chan struct{}
}

func NewApi(s storage.Storage) *API {
	r := chi.NewRouter()
	storageapi.RegisterRoutes(r, s)

	r.Get("/ok", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return &API{
		router: r,
	}
}

func (a *API) Start() {
	go func() {
		a.server = &http.Server{
			Addr:    ":3000",
			Handler: a.router,
		}

		if err := a.server.ListenAndServe(); err != nil {
			log.Fatalf("Server error: %s", err)
		}
		log.Println("server closed")
		close(a.serverClosed)
	}()
}

func (a *API) Stop() {
	a.server.Close()
	<-a.serverClosed
}
