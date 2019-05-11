package main

import (
	"log"
	"net/http"
	"seed/example/admiral"
	"seed/example/admiral/gen"
)

func main() {
	service := gen.New(&admiral.Server{})

	log.Fatal(http.ListenAndServe(":8080", service))
}
