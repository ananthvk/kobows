package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
)

var baseTemplate = template.Must(template.ParseFiles("templates/base.gohtml"))

func main() {
	portPtr := flag.Uint("port", 8000, "specify the port on which to listen")
	hostPtr := flag.String("host", "0.0.0.0", "specify the bind address")
	flag.Parse()
	address := fmt.Sprintf("%s:%d", *hostPtr, *portPtr)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /system-info", sysinfoHandler)
	mux.HandleFunc("POST /upload", ebookFileUploadHandler)
	mux.HandleFunc("GET /", indexHandler)

	slog.Info("listening on", "address", address)
	http.ListenAndServe(address, mux)
}
