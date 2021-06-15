package main

import (
	"github.com/gorilla/mux"
	"goim/handler"
	"goim/middleware"
	"net/http"
)

func main() {
	r := mux.NewRouter()
	handler.RegisterRoutes(r)
	r.Use(middleware.Cors, mux.CORSMethodMiddleware(r))
	err := http.ListenAndServe(":8081", r)
	if err != nil {
		panic(err)
	}
}
