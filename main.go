package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	router := chi.NewRouter() // New em go é geralmente um construtor, é preferivel criar um construtor do que instanciar uma propriedde pois geralmente inicializam propriedades privadas desse tipo
	router.Use(middleware.Logger)

	router.Get("/hello", basicHandler)
	server := &http.Server{ // storing as a pointer and accessing the memory address
		Addr:    ":3001",
		Handler: router,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(" Failed to connect to server", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) { // o * é um pointer
	w.Write([]byte("Gosto muito de bananas"))
}
