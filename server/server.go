package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

type Message struct {
	Text string
}

type myHandler struct{}

/**TODO: Review and see which are needed */
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

func (h *myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Handle CORS issues TODO: Review CORS in detail
	enableCors(&writer)
	if request.Method == "OPTIONS" {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(""))
		return
	}

	//Implement route forwarding, ensure there is a route established for the request
	if handler, ok := mux[request.URL.Path]; ok {
		handler(writer, request)
		return
	}
	http.Error(writer, "Invalid route"+request.URL.Path, http.StatusNotFound)
}

func Write(writer http.ResponseWriter, request *http.Request) {

	//Decode body into formated form
	var msg Message
	err := json.NewDecoder(request.Body).Decode(&msg)

	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println(msg)
}

func Serve(inMux map[string]func(http.ResponseWriter, *http.Request)) error {
	server := http.Server{
		Addr:        ":8080",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
	}

	mux = inMux

	return server.ListenAndServe()
}
