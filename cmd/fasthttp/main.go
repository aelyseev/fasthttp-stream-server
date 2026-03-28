package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/aelyseev/assignments/fasthttp-server/internal/config"
	"github.com/aelyseev/assignments/fasthttp-server/internal/logger"
	"github.com/aelyseev/assignments/fasthttp-server/pkg/server/routes"
	"github.com/fasthttp/router"

	"github.com/valyala/fasthttp"
)

func main() {
	config.Initialize()

	conf := config.GetConfig()

	logger.Initialize(conf.Settings.Level.SlogLevel())

	r := router.New()
	r.POST("/upload", routes.UploadHandler)

	s := &fasthttp.Server{
		Handler:                      r.Handler,
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
		MaxRequestBodySize:           int(conf.Settings.MaxBodySize),
		IdleTimeout:                  conf.Settings.IdleTimeout,
		WriteTimeout:                 conf.Settings.WriteTimeout,
		Concurrency:                  conf.Settings.Concurrency,
	}

	if err := os.WriteFile("server.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		log.Fatalf("failed to write pid file: %v", err)
	}
	defer func() { _ = os.Remove("server.pid") }()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Println("listening on :8080")
		if err := s.ListenAndServe(":8080"); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")

	if err := s.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("server stopped")
}
