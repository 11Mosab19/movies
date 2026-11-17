package logging

import (
	"log/slog"
	"regexp"
)

var (
	credentialURLPattern  = regexp.MustCompile(`(?i)((?:postgres(?:ql)?|cloudinary)://[^:\s/@]+:)[^@\s/]+@`)
	sensitiveValuePattern = regexp.MustCompile(`(?i)\b(password|access[_-]?token|refresh[_-]?token|token|secret|api[_-]?key|authorization)\b\s*([=:])\s*("[^"]*"|'[^']*'|[^,\s;]+)`)
	emailValuePattern     = regexp.MustCompile(`(?i)(email\s*\)?\s*=\s*\(?)[^\s,;\)]+`)
)

func Error(component, operation string, err error) {
	attributes := []any{
		"component", component,
		"operation", operation,
	}
	if err != nil {
		attributes = append(attributes, "cause", redact(err.Error()))
	}

	slog.Error("internal server error", attributes...)
}

func redact(message string) string {
	message = credentialURLPattern.ReplaceAllString(message, "${1}[REDACTED]@")
	message = sensitiveValuePattern.ReplaceAllString(message, "${1}${2}[REDACTED]")
	return emailValuePattern.ReplaceAllString(message, "${1}[REDACTED]")
}
