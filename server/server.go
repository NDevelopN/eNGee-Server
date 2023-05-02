package server

import (
	"net/http"
	"strings"
	"time"

	l "Engee-Server/lobby"
)

type myHandler struct{}

/**TODO: Review and see which are needed */
func enableCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "*")
	(*w).Header().Set("Access-Control-Allow-Headers", "*")
}

var mux = map[string]func(http.ResponseWriter, *http.Request){
	"":       ReMux,
	"server": ReMux,
	"game":   l.ReMux,
}

func (h *myHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Handle CORS issues TODO: Review CORS in detail
	enableCors(&writer)
	if request.Method == "OPTIONS" {
		writer.WriteHeader(http.StatusOK)
		writer.Write([]byte(""))
		return
	}

	// Get the base path and pass to relevant ReMux
	path := strings.Split(request.URL.Path, "/")[1]

	//Implement route forwarding, ensure there is a route established for the request
	if handler, ok := mux[path]; ok {
		handler(writer, request)
		return
	}
	http.Error(writer, "Invalid route: "+request.URL.Path, http.StatusNotFound)
}

func Serve() error {
	server := http.Server{
		Addr:        ":8080",
		Handler:     &myHandler{},
		ReadTimeout: 5 * time.Second,
	}
	return server.ListenAndServe()
}
