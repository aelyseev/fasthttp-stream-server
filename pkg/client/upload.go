package client

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"

	"github.com/valyala/fasthttp"
)

type Uploader struct {
	serverURL string
	chunkSize int64
	client    *fasthttp.Client
}

func NewUploader(serverURL string, chunkSize int64) *Uploader {
	return &Uploader{
		serverURL: serverURL,
		chunkSize: chunkSize,
		client:    &fasthttp.Client{},
	}
}

func (u *Uploader) Upload(filePath string) error {
	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("open file: %w", err)
	}

	pr, pw := io.Pipe()
	mw := multipart.NewWriter(pw)

	errCh := make(chan error, 1)
	go func() {
		errCh <- u.writeMultipart(mw, pw, f, filePath)
	}()

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)
	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	req.Header.SetMethod(fasthttp.MethodPost)
	req.SetRequestURI(u.serverURL)
	req.Header.SetContentType(mw.FormDataContentType())
	req.SetBodyStream(pr, -1)

	doErr := u.client.Do(req, resp)

	writeErr := <-errCh

	if writeErr != nil {
		return fmt.Errorf("write multipart: %w", writeErr)
	}
	if doErr != nil {
		return fmt.Errorf("http request: %w", doErr)
	}
	if resp.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("server returned %d: %s", resp.StatusCode(), resp.Body())
	}
	return nil
}

func (u *Uploader) writeMultipart(mw *multipart.Writer, pw *io.PipeWriter, f *os.File, filePath string) error {
	defer func() {
		_ = f.Close()
		_ = mw.Close()
		_ = pw.Close()
	}()

	fw, err := mw.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		pw.CloseWithError(err)
		return err
	}

	chunk := make([]byte, u.chunkSize)
	for {
		n, readErr := f.Read(chunk)
		if n > 0 {
			if _, werr := fw.Write(chunk[:n]); werr != nil {
				pw.CloseWithError(werr)
				return werr
			}
		}
		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			pw.CloseWithError(readErr)
			return readErr
		}
	}
	return nil
}
