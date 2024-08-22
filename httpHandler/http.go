package httpHandler

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/api-challenge/handler"
)

type HttpHandler struct {
	handler  *handler.Handler
	hostname string
	port     string
}

func NewHttpHandler(hostname string, port string, handler *handler.Handler) *HttpHandler {
	return &HttpHandler{
		hostname: hostname,
		port:     port,
		handler:  handler,
	}
}

func (r *HttpHandler) Init() {
	ginEngine := gin.Default()

	InitRoutes(r.handler, ginEngine)

	srv := &http.Server{
		Addr:    r.hostname + ":" + r.port,
		Handler: ginEngine.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shut down requested...")
	cancelContext, stdServer := context.WithTimeout(context.Background(), time.Second*30)
	defer stdServer()
	if err := srv.Shutdown(cancelContext); err != nil {
		log.Fatalf("Error shutting down the server: %v", err)
	}
	log.Println("Server exited")
}
