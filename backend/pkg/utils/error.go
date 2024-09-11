package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type ErrorResponseStatus struct {
	Timestamp    int64    `json:"timestamp"`
	ErrorCode    int      `json:"status_code"`
	ErrorMessage []string `json:"error_message"`
}

type ErrorResponse struct {
	Status ErrorResponseStatus `json:"status"`
}

func Err(w http.ResponseWriter, errorCode int, errorMessages []string) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(errorCode)

	errorResponse := &ErrorResponse{
		Status: ErrorResponseStatus{
			Timestamp:    time.Now().Unix(),
			ErrorCode:    errorCode,
			ErrorMessage: errorMessages,
		},
	}
	err := json.NewEncoder(w).Encode(&errorResponse)
	if err != nil {
		zap.S().Error("Error encoding error response data")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
