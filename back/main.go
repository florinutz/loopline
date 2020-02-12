package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	. "back/pkg"
)

var cliConfig = struct {
	listenAddr string
}{}

func init() {
	cliConfig.listenAddr = *flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
}

func main() {
	api := NewController(*NewInMemoryStorage())
	api.Router.Use(loggingMiddleware)
	fmt.Printf("listening on %s\n", cliConfig.listenAddr)
	log.Println(http.ListenAndServe(cliConfig.listenAddr, api.Router))
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[%s] %s\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
