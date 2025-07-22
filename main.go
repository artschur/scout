package main

import (
	"go-observability-tool/internal/routes"
	"net/http"
)

func main() {
	mux := http.DefaultServeMux
	routes.AddRoutes(mux)

	http.ListenAndServe("8082", mux)
}
