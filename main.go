package main

import (
	"flag"
	"fmt"
	"log/slog"
	"net/http"
)

func main() {
	portPtr := flag.Uint("port", 8000, "specify the port on which to listen")
	hostPtr := flag.String("host", "0.0.0.0", "specify the bind address")
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *hostPtr, *portPtr)

	mux := http.NewServeMux()
	mux.HandleFunc("/", indexHandler)

	slog.Info("listening on", "address", address)
	http.ListenAndServe(address, mux)
}
