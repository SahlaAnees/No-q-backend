package server

import (
	"context"
	"log"
	"net/http"
	"no-q-solution/utils/config"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type HTTPServer struct {
	server  *http.Server
	address string
}

func NewHTTPServer(config config.Config, r *mux.Router) HTTPServer {

	address := config.App.Host + ":" + strconv.Itoa(config.App.Port)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
	})

	handler := c.Handler(r)

	server := &http.Server{
		Addr:         address,
		WriteTimeout: time.Second * 10,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Second * 10,

		Handler: handler,
	}

	httpServer := HTTPServer{
		server:  server,
		address: address,
	}

	return httpServer
}

func (srv HTTPServer) ListnAndServe(ctx context.Context) {

	log.Printf("server listening on %s", srv.address)

	err := srv.server.ListenAndServe()
	if err != nil {
		log.Println(err)
	}
}

func (srv HTTPServer) Shutdown(ctx context.Context) {

	log.Println("stropping HTTP server")

	srv.server.SetKeepAlivesEnabled(false)

	err := srv.server.Shutdown(ctx)
	if err != nil {
		log.Println(err)
	}
}
