package httpserver

import (
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
)

type StorageItf interface{}

type HttpServer struct {
	storage     StorageItf
	siteHandler http.Handler
}

func InitAndStart(storage StorageItf) *HttpServer {
	httpServer := Init(storage)
	httpServer.Start()
	return httpServer
}

func Init(storage StorageItf) *HttpServer {
	httpServer := &HttpServer{
		storage: storage,
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/health", httpServer.health)

	// Middlewares
	httpServer.siteHandler = httpServer.metricsMiddleware(mux)
	httpServer.siteHandler = httpServer.accessControlMiddleware(httpServer.siteHandler)
	httpServer.siteHandler = httpServer.panicMiddleware(httpServer.siteHandler)

	return httpServer
}

func (httpServer *HttpServer) Start() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "80"
	}
	log.Info().Msgf("Start HTTP Server at :%v", port)
	log.Panic().Msgf("HTTP Server ListenAndServe fatal error: %v", http.ListenAndServe(":"+port, httpServer.siteHandler))
}

func (httpServer *HttpServer) health(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
