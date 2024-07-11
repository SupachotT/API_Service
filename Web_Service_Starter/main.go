package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	http.HandleFunc("/greet/", func(w http.ResponseWriter, r *http.Request) {
        name := r.URL.Path[len("/greet/"):]
        fmt.Fprintf(w, "Hello, %s!", name)
    })

	http.ListenAndServe(":8081", nil)
}
