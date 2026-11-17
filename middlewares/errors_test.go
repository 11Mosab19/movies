package middlewares

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"movies/apperrors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestErrorMiddleWare(t *testing.T) {
	gin.SetMode(gin.TestMode)

	previousLogger := slog.Default()
	var logs bytes.Buffer
	slog.SetDefault(slog.New(slog.NewTextHandler(&logs, nil)))
	t.Cleanup(func() {
		slog.SetDefault(previousLogger)
	})

	tests := []struct {
		name        string
		err         error
		wantStatus  int
		wantMessage string
		wantLog     bool
	}{
		{
			name:        "expected app error",
			err:         apperrors.ErrMovieNotFound,
			wantStatus:  http.StatusNotFound,
			wantMessage: "movie not found",
		},
		{
			name:        "wrapped expected app error",
			err:         fmt.Errorf("get movie: %w", apperrors.ErrMovieNotFound),
			wantStatus:  http.StatusNotFound,
			wantMessage: "movie not found",
		},
		{
			name: "expected error joined with internal failure",
			err: errors.Join(
				fmt.Errorf("create movie: %w", apperrors.ErrMovieAlreadyExists),
				errors.New("cleanup failed"),
			),
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "internal server error",
			wantLog:     true,
		},
		{
			name:        "unexpected error is redacted",
			err:         errors.New("database query failed: password=secret-value token=access-token email=user@example.com"),
			wantStatus:  http.StatusInternalServerError,
			wantMessage: "internal server error",
			wantLog:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logs.Reset()
			router := gin.New()
			router.Use(ErrorMiddleWare())
			router.GET("/test", func(c *gin.Context) {
				_ = c.Error(tt.err)
				c.Abort()
			})

			response := httptest.NewRecorder()
			request := httptest.NewRequest(http.MethodGet, "/test", nil)
			router.ServeHTTP(response, request)

			if response.Code != tt.wantStatus {
				t.Fatalf("status = %d, want %d", response.Code, tt.wantStatus)
			}
			if !strings.Contains(response.Body.String(), tt.wantMessage) {
				t.Fatalf("response = %q, want message %q", response.Body.String(), tt.wantMessage)
			}
			if tt.wantLog != (logs.Len() > 0) {
				t.Fatalf("log present = %v, want %v", logs.Len() > 0, tt.wantLog)
			}
			for _, sensitiveValue := range []string{"secret-value", "access-token", "user@example.com"} {
				if strings.Contains(logs.String(), sensitiveValue) {
					t.Fatalf("log contains a sensitive value: %q", logs.String())
				}
			}
		})
	}
}
