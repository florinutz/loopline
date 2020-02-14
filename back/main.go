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

var cliConfig = struct {
	listenAddr string
}{}

func init() {
	cliConfig.listenAddr = *flag.String("addr", ":8080", "HTTP listen address")
	flag.Parse()
}

func main() {
	api := NewController(*NewInMemoryStorage())
	fmt.Printf("listening on %s\n\n", cliConfig.listenAddr)
	log.Println(http.ListenAndServe(cliConfig.listenAddr, decorateHandler(api.Router)))
}

func decorateHandler(router *mux.Router) http.Handler {
	// add apache combined logging
	handler := handlers.CombinedLoggingHandler(os.Stdout, router)

	// allow all CORS
	handler = handlers.CORS(
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "DELETE", "OPTIONS"}),
	)(handler)

	// recover panics
	handler = handlers.RecoveryHandler(handlers.PrintRecoveryStack(true))(handler)

	return handler
}
