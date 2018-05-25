package main

import (
	"encoding/json"
	"log"
	"net/http"
)

// Handler HTTP處理
type Handler int

// ServeHTTP 服務處理
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	type Input struct {
		Method  string      `json:"method"`
		Params  interface{} `json:"params"`
		ID      int         `json:"id"`
		Address string      `json:"address"`
	}
	type ErrorDetail struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}
	type Output struct {
		Result interface{} `json:"result"`
		Error  interface{} `json:"error"`
		ID     int         `json:"id"`
	}
	w.Header().Set("Content-Type", "application/json")
	var data Input
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("JSON Error ->", err)
		json.NewEncoder(w).Encode(Output{
			Result: nil,
			Error: ErrorDetail{
				Code:    500,
				Message: err.Error(),
				Data:    err,
			},
			ID: data.ID,
		})
		return
	}
	log.Println("Input ->", data)
	address := data.Address
	if address == "" {
		address = ":50052"
	}
	res, err := transferJSONRPCClient(address, data.Method, data.Params)
	if err != nil {
		log.Println("Call Error ->", err)
		json.NewEncoder(w).Encode(Output{
			Result: nil,
			Error: ErrorDetail{
				Code:    500,
				Message: err.Error(),
				Data:    err,
			},
			ID: data.ID,
		})
		return
	}

	err = json.NewEncoder(w).Encode(Output{
		Result: res,
		Error:  nil,
		ID:     data.ID,
	})
	if err != nil {
		log.Println("Response Encode Error ->", err)
		return
	}
}
