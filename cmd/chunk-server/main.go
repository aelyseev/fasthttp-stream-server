package main

import (
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/valyala/fasthttp"
)

const (
	readBufSize = 64 << 10
	listenAddr  = ":8081"
)

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, readBufSize)
		return &b
	},
}

func uploadHandler(ctx *fasthttp.RequestCtx) {
	boundary := string(ctx.Request.Header.MultipartFormBoundary())
	if boundary == "" {
		ctx.Error("expected multipart/form-data", fasthttp.StatusBadRequest)
		return
	}

	body := ctx.RequestBodyStream()
	if body == nil {
		ctx.Error("streaming not enabled", fasthttp.StatusInternalServerError)
		return
	}
	defer func() { _, _ = io.Copy(io.Discard, body) }()

	mr := multipart.NewReader(body, boundary)

	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	buf := *bufPtr

	var total int64
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("invalid multipart body: %v", err)
			ctx.Error("invalid multipart body", fasthttp.StatusBadRequest)
			return
		}

		if part.FormName() != "file" {
			_, _ = io.Copy(io.Discard, part)
			continue
		}

		for {
			n, readErr := part.Read(buf)
			total += int64(n)
			if readErr == io.EOF {
				break
			}
			if readErr != nil {
				log.Printf("read error: %v", readErr)
				ctx.Error("read error", fasthttp.StatusInternalServerError)
				return
			}
		}
	}

	log.Printf("upload complete: %d bytes received", total)
	ctx.SetStatusCode(fasthttp.StatusOK)
	ctx.SetContentType("text/plain; charset=utf-8")
	ctx.SetBodyString("ok\n")
}

func main() {
	s := &fasthttp.Server{
		Handler:                      uploadHandler,
		StreamRequestBody:            true,
		DisablePreParseMultipartForm: true,
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		log.Printf("chunk-server listening on %s", listenAddr)
		if err := s.ListenAndServe(listenAddr); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	<-quit
	log.Println("shutting down...")
	if err := s.Shutdown(); err != nil {
		log.Printf("shutdown error: %v", err)
	}
	log.Println("stopped")
}
