package http

import (
	"Todo"
	"encoding/json"
	"log"
	"net/http"
)

func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Extract error code & message.
	code, message := Todo.ErrorCode(err), Todo.ErrorMessage(err)

	// Log & report internal errors.
	if code == Todo.EINTERNAL {
		Todo.ReportError(r.Context(), err, r)
		LogError(r, err)
	}

	// Print user message to response based on request accept header.
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))
	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}

var codes = map[string]int{
	Todo.ECONFLICT:       http.StatusConflict,
	Todo.EINVALID:        http.StatusBadRequest,
	Todo.ENOTFOUND:       http.StatusNotFound,
	Todo.ENOTIMPLEMENTED: http.StatusNotImplemented,
	Todo.EUNAUTHORIZED:   http.StatusUnauthorized,
	Todo.EINTERNAL:       http.StatusInternalServerError,
}

func LogError(r *http.Request, err error) {
	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}
