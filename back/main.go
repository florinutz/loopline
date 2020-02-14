package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	. "back/pkg"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var conf = struct{ listenAddr string }{}

func main() {
	api := NewController(*NewInMemoryStorage())
	fmt.Printf("listening on %s\n\n", conf.listenAddr)
	log.Println(http.ListenAndServe(conf.listenAddr, addHandlers(api.Router)))
}

func addHandlers(router *mux.Router) http.Handler {
	// add apache combined logging
	handler := handlers.CombinedLoggingHandler(os.Stdout, router)

	// allow all CORS for now
	handler = handlers.CORS(
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
	)(handler)

	// recover panics
	handler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)

	return handler
}

func init() {
	defaultListen := ":8080"
	if port := os.Getenv("PORT"); port != "" {
		defaultListen = fmt.Sprintf(":%s", port)
	}
	conf.listenAddr = *flag.String("addr", defaultListen, "HTTP listen address")
	flag.Parse()
}
