package response

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Code    int         `json:"code,omitempty"`
}

func JSON(w http.ResponseWriter, statusCode int, data interface{}, message ...string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	msg := ""
	if len(message) > 0 {
		msg = message[0]
	}

	res := Response{
		Status:  "success",
		Message: msg,
		Data:    data,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}

func Error(w http.ResponseWriter, statusCode int, message string, code ...int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errCode := 0
	if len(code) > 0 {
		errCode = code[0]
	}

	res := Response{
		Status:  "error",
		Message: message,
		Code:    errCode,
	}

	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "failed to encode error response", http.StatusInternalServerError)
	}
}

func Write(w http.ResponseWriter, statusCode int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(payload); err != nil {
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
	}
}
