package router

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/patrickn2/Clerk-Challenge/handler"
)

func Init(h *handler.Handler) {
	router := gin.Default()
	router.Use(h.Middleware())
	InitRoutes(h, router)
	port := os.Getenv("API_PORT")
	if port == "" {
		log.Fatal("api port not defined")
	}

	srv := &http.Server{
		Addr:    ":" + port,
		Handler: router.Handler(),
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	h.Shutdown.ShutdownSignal = true
	log.Println("Shut down requested...")
	if h.Shutdown.RunningProcess > 0 {
		log.Println("Shutting down after finish all current requests: ", h.Shutdown.RunningProcess)
	}
	for {
		if h.Shutdown.ShutdownSignal {
			if h.Shutdown.RunningProcess > 0 {
				fmt.Printf("Server is shutting down: %d remaining requests to shutdown\n", h.Shutdown.RunningProcess)
				time.Sleep(time.Second)
				continue
			}
			break
		}
	}
	log.Println("Server exited")
}
