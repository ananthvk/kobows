package main

import (
	"flag"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
)

var baseTemplate = template.Must(template.ParseFiles("templates/base.gohtml"))
var uploadBaseDirectory = "/mnt/onboard/uploads"

func main() {
	portPtr := flag.Uint("port", 8000, "specify the port on which to listen")
	hostPtr := flag.String("host", "0.0.0.0", "specify the bind address")
	uploadsPath := flag.String("upload-path", "/mnt/onboard/uploads/", "specify the upload directory")
	flag.Parse()

	address := fmt.Sprintf("%s:%d", *hostPtr, *portPtr)
	uploadBaseDirectory = *uploadsPath

	if err := os.MkdirAll(uploadBaseDirectory, 0o755); err != nil {
		slog.Error("failed to create upload directory", "error", err, "path", uploadBaseDirectory)
		return
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /system-info", sysinfoHandler)
	mux.HandleFunc("POST /upload", ebookFileUploadHandler)
	mux.HandleFunc("GET /", indexHandler)

	slog.Info("listening on", "address", address, "upload-path", *uploadsPath)
	http.ListenAndServe(address, mux)
}
