package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithJson(w http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON: %v error: %v", payload, err)
		w.WriteHeader(500)
		return
	}


	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(data)

}