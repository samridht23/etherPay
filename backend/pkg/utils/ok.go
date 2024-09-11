package utils

import (
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Response struct {
	Data   interface{}    `json:"data,omitempty"` // Optional data field
	Status ResponseStatus `json:"status"`
}

type ResponseStatus struct {
	Timestamp  int64 `json:"timestamp"`
	StatusCode int   `json:"status_code"`
}

func Send(w http.ResponseWriter, statusCode int, data interface{}) {

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(statusCode)

	response := &Response{
		Data: data,
		Status: ResponseStatus{
			Timestamp:  time.Now().Unix(),
			StatusCode: statusCode,
		},
	}

	err := json.NewEncoder(w).Encode(&response)
	if err != nil {
		zap.S().Errorf("Error encoding response data: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
