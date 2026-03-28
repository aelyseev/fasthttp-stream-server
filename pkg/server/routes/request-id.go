package routes

import (
	"crypto/rand"
	"encoding/hex"

	"github.com/valyala/fasthttp"
)

func requestID(req *fasthttp.RequestCtx) string {
	if id := req.Request.Header.Peek("X-Request-ID"); len(id) > 0 {
		return string(id)
	}
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
