package middleware

import (
	"fmt"
	"io"
	"net/http"

	"github.com/felixge/httpsnoop"
)

type errorLoggingHandler struct {
	writer  io.Writer
	handler http.Handler
}

type errorLogger struct {
	w        http.ResponseWriter
	status   int
	size     int
	response string
}

func (l *errorLogger) Write(b []byte) (int, error) {
	//wrap Write to get error text
	l.response = string(b)

	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *errorLogger) WriteHeader(s int) {
	//wrap WriteHeader to get error code
	l.w.WriteHeader(s)
	l.status = s
}

func (h errorLoggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	logger, w := makeLogger(w)
	//url := *req.URL

	h.handler.ServeHTTP(w, req)
	if req.MultipartForm != nil {
		req.MultipartForm.RemoveAll()
	}

	if logger.status != http.StatusOK {
		fmt.Fprintf(h.writer, "Error response: %s", logger.response)
	}
}

func makeLogger(w http.ResponseWriter) (*errorLogger, http.ResponseWriter) {
	logger := &errorLogger{w: w, status: http.StatusOK}
	return logger, httpsnoop.Wrap(w, httpsnoop.Hooks{
		Write: func(httpsnoop.WriteFunc) httpsnoop.WriteFunc {
			return logger.Write
		},
		WriteHeader: func(httpsnoop.WriteHeaderFunc) httpsnoop.WriteHeaderFunc {
			return logger.WriteHeader
		},
	})
}

func ErrorLoggingHandler(out io.Writer, h http.Handler) http.Handler {
	return errorLoggingHandler{out, h}
}
