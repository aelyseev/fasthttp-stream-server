package routes

import (
	"errors"
	"io"
	"strings"
)

func isClientDisconnected(err error) bool {
	if errors.Is(err, io.ErrUnexpectedEOF) {
		return true
	}
	s := err.Error()
	return strings.Contains(s, "connection reset by peer") ||
		strings.Contains(s, "broken pipe")
}
