package main

import (
	"fmt"
	"net/http"
)

func main() {
	server := &http.Server{ // storing as a pointer and accessing the memory address
		Addr:    ":3000",
		Handler: http.HandlerFunc(basicHandler),
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(" Failed to connect to server", err)
	}
}

func basicHandler(w http.ResponseWriter, r *http.Request) { // o * é um pointer
	w.Write([]byte("Gosto muito de bananas"))
}
