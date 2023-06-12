package middleware

import (
	"fmt"
	"io"
	"net/http"
)

type errorLoggingHandler struct {
	writer  io.Writer
	handler http.Handler
}

type errorLogger struct {
	http.ResponseWriter
	status   int
	size     int
	response string
}

func (l *errorLogger) Write(b []byte) (int, error) {
	//wrap Write to get error text
	l.response = string(b)

	size, err := l.ResponseWriter.Write(b)
	l.size += size
	return size, err
}

func (l *errorLogger) WriteHeader(s int) {
	//wrap WriteHeader to get error code
	l.ResponseWriter.WriteHeader(s)
	l.status = s
}

/*
//we dont need it as ResponseWriter is embedded and errorLogger satisfies ResponseWriter interface without it
func (l *errorLogger) Header() http.Header {
	return l.ResponseWriter.Header()
}*/

func (h errorLoggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	logger := &errorLogger{ResponseWriter: w, status: http.StatusOK}

	h.handler.ServeHTTP(logger, req)
	if req.MultipartForm != nil {
		req.MultipartForm.RemoveAll()
	}

	if logger.status != http.StatusOK {
		fmt.Fprintf(h.writer, "Error response: %s", logger.response)
		if len(logger.response) == 0 || logger.response[len(logger.response)-1] != '\n' {
			fmt.Fprintf(h.writer, "\n")
		}
	}
}

func ErrorLoggingHandler(out io.Writer) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return errorLoggingHandler{out, h}
	}
}
