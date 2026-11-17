package middlewares

import (
	"errors"
	"movies/apperrors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestRequireRole(t *testing.T) {
	gin.SetMode(gin.TestMode)
	requireAdmin := (&AuthMiddleware{}).RequireRole("admin")

	tests := []struct {
		name        string
		role        any
		setRole     bool
		wantError   *apperrors.AppError
		wantAborted bool
	}{
		{
			name:        "missing role",
			wantError:   apperrors.ErrUnauthorized,
			wantAborted: true,
		},
		{
			name:        "invalid role type",
			role:        1,
			setRole:     true,
			wantError:   apperrors.ErrUnauthorized,
			wantAborted: true,
		},
		{
			name:        "regular user",
			role:        "user",
			setRole:     true,
			wantError:   apperrors.ErrForbidden,
			wantAborted: true,
		},
		{
			name:        "admin user",
			role:        "admin",
			setRole:     true,
			wantAborted: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, _ := gin.CreateTestContext(httptest.NewRecorder())
			if tt.setRole {
				c.Set("Role", tt.role)
			}

			requireAdmin(c)

			if c.IsAborted() != tt.wantAborted {
				t.Fatalf("aborted = %v, want %v", c.IsAborted(), tt.wantAborted)
			}
			if tt.wantError == nil {
				if len(c.Errors) != 0 {
					t.Fatalf("unexpected error: %v", c.Errors.Last())
				}
				return
			}

			if len(c.Errors) != 1 {
				t.Fatalf("error count = %d, want 1", len(c.Errors))
			}

			var appErr *apperrors.AppError
			if !errors.As(c.Errors.Last().Err, &appErr) {
				t.Fatalf("error type = %T, want *apperrors.AppError", c.Errors.Last().Err)
			}
			if appErr.Code != tt.wantError.Code {
				t.Fatalf("error code = %d, want %d", appErr.Code, tt.wantError.Code)
			}
		})
	}
}
