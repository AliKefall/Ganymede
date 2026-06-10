package endpoints

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"
)

type ErrorResponse struct {
	Error ErrorBody `json:"error"`
}

type ErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"request_id,omitempty"`
}

func RespondWithJSON(
	w http.ResponseWriter,
	status int,
	payload any,
) {
	body, err := json.Marshal(payload)
	if err != nil {
		slog.Error(
			"failed to marshal response",
			"error", err,
		)

		body = []byte(`{
			"error":{
				"code":"internal_server_error",
				"message":"Internal server error"
			}
		}`)

		status = http.StatusInternalServerError
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)

	if _, err := w.Write(body); err != nil {
		slog.Error(
			"failed to write response",
			"error", err,
		)
	}
}

func RespondWithError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
	requestID string,
	err error,
) {
	if err != nil {
		slog.Error("request failed",
			"status", status,
			"code", code,
			"error", err,
			"request_id", requestID,
		)
	}

	resp := ErrorResponse{
		Error: ErrorBody{
			Code:      code,
			Message:   message,
			RequestID: requestID,
		},
	}

	RespondWithJSON(w, status, resp)

}

func shouldUserSecureCookie(r *http.Request) bool {
	if r.TLS != nil {
		return true
	}
	if strings.EqualFold(r.Header.Get("X-Forwarded-For"), "https"){
		return true
	}
	return false
}
