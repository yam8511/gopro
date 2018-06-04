package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)
import "net/http/pprof"

// 比較偏向機器or連線上的錯誤，而且給外部看
const (
	// Parse error
	E_PARSE int = -32700
	// Invalid Request
	E_INVALID_REQ int = -32600
	// Method not found
	E_NO_METHOD int = -32601
	// Invalid params
	E_BAD_PARAMS int = -32602
	// Internal error
	E_INTERNAL int = -32603
	// Server error (-32000 to -32099)
	E_SERVER int = -32000
)

// Handler HTTP處理
type Handler int

// ServeHTTP 服務處理
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.RequestURI == "/favicon.ico" {
		w.WriteHeader(http.StatusOK)
		return
	}
	// log.Println("Req -> ", r.URL.EscapedPath())
	switch r.URL.EscapedPath() {
	case "/debug/pprof/cmdline":
		pprof.Cmdline(w, r)
		return
	case "/debug/pprof/profile":
		pprof.Profile(w, r)
		return
	case "/debug/pprof/symbol":
		pprof.Symbol(w, r)
		return
	case "/debug/pprof/trace":
		pprof.Trace(w, r)
		return
	}
	if strings.HasPrefix(r.RequestURI, "/debug/") {
		r.URL.Path = strings.Replace(r.URL.Path, "/debug/pprof", "/debug", 1)
		// log.Println("#1 Req -> ", r.URL.Path)
		r.URL.Path = strings.Replace(r.URL.Path, "/debug", "/debug/pprof", 1)
		// log.Println("#2 Req -> ", r.URL.Path)
		pprof.Index(w, r)
		return
	}
	// log.Println("Not Debug Req -> ", r.RequestURI)

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
		JSONRPC string      `json:"jsonrpc"`
		Result  interface{} `json:"result"`
		Error   interface{} `json:"error"`
		ID      int         `json:"id"`
	}
	w.Header().Set("Content-Type", "application/json")
	var data Input
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println("JSON Error ->", err)
		json.NewEncoder(w).Encode(Output{
			JSONRPC: "2.0",
			Result:  nil,
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
			JSONRPC: "2.0",
			Result:  nil,
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
		JSONRPC: "2.0",
		Result:  res,
		Error:   nil,
		ID:      data.ID,
	})
	if err != nil {
		log.Println("Response Encode Error ->", err)
		return
	}
}
