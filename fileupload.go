package main

import (
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gabriel-vasile/mimetype"
)

const maxEbookFileSize = 50 * 1000 * 1000 // Max size of 50 MB
var fileUploadSuccessTemplate = template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/file_upload_done.gohtml"))
var fileUploadErrorTemplate = template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/file_upload_error.gohtml"))

func ebookFileUploadHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request details", "method", r.Method, "path", r.URL.Path, "from", r.RemoteAddr)
	r.Body = http.MaxBytesReader(w, r.Body, maxEbookFileSize)
	err := r.ParseMultipartForm(maxEbookFileSize)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fileUploadErrorTemplate.ExecuteTemplate(w, "base", map[string]string{"Error": "Could not read multipart: " + err.Error()})
		return
	}

	file, header, err := r.FormFile("ebook")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fileUploadErrorTemplate.ExecuteTemplate(w, "base", map[string]string{"Error": err.Error()})
		return
	}
	defer file.Close()

	mtype, err := mimetype.DetectReader(file)
	if err != nil && err != io.EOF {
		w.WriteHeader(http.StatusBadRequest)
		fileUploadErrorTemplate.ExecuteTemplate(w, "base", map[string]string{"Error": "Cannot read file: " + err.Error()})
		return
	}

	if mtype != mimetype.Lookup("application/epub+zip") {
		w.WriteHeader(http.StatusBadRequest)
		fileUploadErrorTemplate.ExecuteTemplate(w, "base", map[string]string{"Error": "File type not supported"})
		return
	}

	_, err = file.Seek(0, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("failed to seek file body", "error", err)
		fileUploadErrorTemplate.ExecuteTemplate(w, "base", map[string]string{"Error": "Internal error"})
		return
	}

	// Write to disk
	destinationFile, err := getDestinationEpubFile(header.Filename)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("failed to create destination file", "error", err)
		fileUploadErrorTemplate.ExecuteTemplate(w, "base", map[string]string{"Error": "Internal error"})
	}
	defer destinationFile.Close()

	if _, err := io.Copy(destinationFile, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		slog.Error("failed to copy file", "error", err)
		fileUploadErrorTemplate.ExecuteTemplate(w, "base", map[string]string{"Error": "Internal error"})
	}
	destinationFile.Sync()

	slog.Info("upload success", "filename", header.Filename, "size", header.Size, "mimeType", header.Header, "from", r.RemoteAddr)
	fileUploadSuccessTemplate.ExecuteTemplate(w, "base", map[string]string{"Filename": header.Filename})
}

func getDestinationEpubFile(filename string) (*os.File, error) {
	filename = filepath.Clean(filename)
	filename = filepath.Base(filename)

	if !filepath.IsLocal(filename) {
		return nil, errors.New("invalid filename")
	}

	if filepath.Ext(filename) != ".epub" {
		filename += ".epub"
	}

	destination := filepath.Join(uploadBaseDirectory, filename)
	return os.Create(destination)
}
