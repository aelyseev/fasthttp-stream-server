package routes

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"sync"

	"github.com/aelyseev/assignments/fasthttp-server/internal/logger"
	"github.com/valyala/fasthttp"
)

// TcpReceiveBufferSize 64 KB matches typical kernel TCP receive buffer chunk size,
// balancing syscall overhead and memory usage per concurrent request.
const TcpReceiveBufferSize = 64 << 10

var bufPool = sync.Pool{
	New: func() any {
		b := make([]byte, TcpReceiveBufferSize)
		return &b
	},
}

func UploadHandler(req *fasthttp.RequestCtx) {
	log := logger.GetLogger()

	ctx := log.WithFields(context.Background(), slog.String("request_id", requestID(req)))

	log.Info(ctx, "upload started")
	defer log.Info(ctx, "upload finished")

	boundary := string(req.Request.Header.MultipartFormBoundary())
	if boundary == "" {
		log.Error(ctx, "expected multipart/form-data")
		req.Error("expected multipart/form-data", fasthttp.StatusBadRequest)
		return
	}

	body := req.RequestBodyStream()
	if body == nil {
		req.Error("streaming is not enabled", fasthttp.StatusInternalServerError)
		return
	}

	mr := multipart.NewReader(body, boundary)

	bufPtr := bufPool.Get().(*[]byte)
	defer bufPool.Put(bufPtr)
	buf := *bufPtr
	var total int64
	foundFile := false

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Error(ctx, "invalid multipart body", slog.String("error", err.Error()))
			req.Error("invalid multipart body", fasthttp.StatusBadRequest)
			return
		}

		if part.FormName() != "file" {
			if err := part.Close(); err != nil {
				log.Error(ctx, "failed to discard part", slog.String("error", err.Error()))
				req.Error("failed to read multipart body", fasthttp.StatusInternalServerError)
				return
			}
			continue
		}

		foundFile = true

		for {
			n, err := part.Read(buf)
			if n > 0 {
				total += int64(n)
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				if isClientDisconnected(err) {
					log.Warn(ctx, "client disconnected", slog.String("error", err.Error()), slog.Int64("bytes_read", total))
				} else {
					log.Error(ctx, "failed to read file part", slog.String("error", err.Error()), slog.Int64("bytes_read", total))
					req.Error("failed to read file part", fasthttp.StatusInternalServerError)
				}
				return
			}
		}
	}

	if !foundFile {
		log.Error(ctx, "missing [file] field")
		req.Error("missing [file] field", fasthttp.StatusBadRequest)
		return
	}

	req.SetStatusCode(fasthttp.StatusOK)
	req.SetContentType("text/plain; charset=utf-8")
	req.SetBodyString(fmt.Sprintf("ok size=%d\n", total))
}
