package main

import (
	"log"
	"net/http"
	"seed/example"
	"seed/example/gen"
)

func main() {
	service := gen.New(&example.MyServer{})

	log.Fatal(http.ListenAndServe(":8080", service))
}
