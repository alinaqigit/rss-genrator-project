package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func responseWithError(res http.ResponseWriter, code int, msg string){
	if(code > 499){
		log.Println("Responding with 5XX")
	}
	type errResponse struct {
		Error string `json:"error"`
	}

	responseWithJson(res, code, errResponse{
		Error: msg,
	})
}

func responseWithJson(res http.ResponseWriter, code int, payload any) {
	data, err := json.Marshal(payload)

	if err != nil {
		log.Printf("Failed to marshal JSON: %v error: %v", payload, err)
		res.WriteHeader(500)
		return
	}


	res.Header().Add("Content-Type", "application/json")
	res.WriteHeader(code)
	res.Write(data)

}