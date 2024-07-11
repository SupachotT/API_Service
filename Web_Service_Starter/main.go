package main

import (
	"fmt"
	"net/http"
	"text/template"
	"time"
)

// function for log request to show in sonsole
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.URL)
		next.ServeHTTP(w, r)
	})
}

// function for manage parameter from URL and show in HTML template
func greetHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/greet/"):]
	tmpl, err := template.New("greet").Parse("<html><body><h1>Hello, {{.}}!</h1></body></html>")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// function called Hello World
func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, World!")
	})

	mux.HandleFunc("/greet/", greetHandler)

	server := &http.Server{
		Addr:         ":8081",
		Handler:      loggingMiddleware(mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	err := server.ListenAndServe()
	if err != nil {
		fmt.Println(err)
	}
}
