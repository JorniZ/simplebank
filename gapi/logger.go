package gapi

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	startTime := time.Now()

	resp, err = handler(ctx, req)

	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "gRPC").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received gRPC request")

	return
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(body []byte) (int, error) {
	rec.Body = body
	return rec.ResponseWriter.Write(body)
}

var successCodes = []int{http.StatusOK, http.StatusCreated, http.StatusAccepted}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(resp http.ResponseWriter, req *http.Request) {

		startTime := time.Now()
		rec := &ResponseRecorder{
			ResponseWriter: resp,
			StatusCode:     http.StatusOK,
		}

		handler.ServeHTTP(rec, req)
		duration := time.Since(startTime)

		statusCode := rec.StatusCode

		logger := log.Info()
		if !slices.Contains(successCodes, statusCode) {
			logger = log.Error().Bytes("body", rec.Body)
		}

		logger.Str("protocol", "http").
			Str("path", req.RequestURI).
			Str("method", req.Method).
			Int("status_code", statusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration).
			Msg("received an HTTP request")
	})
}
