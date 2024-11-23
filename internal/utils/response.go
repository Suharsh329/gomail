package utils

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func Response(w http.ResponseWriter, status int, message any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	var messageMap map[string]string

	if status < 400 { // Non-error status
		messageMap = message.(map[string]string)
	} else {
		messageMap = map[string]string{"error": message.(string)}
	}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(messageMap); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	trimmedResponse := strings.TrimSpace(buf.String())
	if _, err := w.Write([]byte(trimmedResponse)); err != nil {
		log.Printf("Failed to write response: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
