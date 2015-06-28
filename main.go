package main

import (
	"fmt"
	"io"
	"net/http"
	"handlers"
)

var mux map[string]func(http.ResponseWriter, *http.Request)

func printMux() {
		for key, value := range mux {
		fmt.Println("Key:", key, "Value:", value)
	}
}

func main() {
	server := http.Server{
		Addr:    "0.0.0.0:8005",
		Handler: &myHandler{},
	}

	mux = make(map[string]func(http.ResponseWriter, *http.Request))
	mux["/"] = handlers.HandleRequest
	mux["/kakki"] = handlers.HandleRequest1
	mux["/foo"] = handlers.HandleRequest2

	fmt.Println("Starting...")
	server.ListenAndServe()
}

type myHandler struct{}

func (*myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Request: " + r.URL.String())
	if h, ok := mux[r.URL.String()]; ok {
		h(w, r)
		return
	}

	io.WriteString(w, "My server: "+r.URL.String())
}
