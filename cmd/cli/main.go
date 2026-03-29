package main

import (
	"fmt"
	"log"
	"os"

	"github.com/aelyseev/assignments/fasthttp-server/internal/config"
	"github.com/aelyseev/assignments/fasthttp-server/internal/logger"
	"github.com/aelyseev/assignments/fasthttp-server/pkg/client"
)

func cleanup() {
	_ = os.Remove("cli.pid")
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("usage: cli <file-path>")
	}
	filePath := os.Args[len(os.Args)-1]

	config.Initialize()
	conf := config.GetConfig()

	logger.Initialize(conf.Cli.Level.SlogLevel())

	if err := os.WriteFile("cli.pid", []byte(fmt.Sprintf("%d", os.Getpid())), 0644); err != nil {
		log.Fatalf("failed to write pid file: %v", err)
	}
	defer cleanup()

	log.Printf("uploading %s in chunks of %d bytes to %s", filePath, conf.Cli.ChunkSize, conf.Cli.ServerURL)

	uploader := client.NewUploader(conf.Cli.ServerURL, conf.Cli.ChunkSize)
	if err := uploader.Upload(filePath); err != nil {
		cleanup()
		log.Fatalf("upload failed: %v", err)
	}

	log.Println("upload complete")
}
