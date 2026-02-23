package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"time"
)

var homepageTemplate = template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/index.gohtml"))
var static = http.FileServer(http.Dir("./static"))
var startTime = time.Now()

func indexHandler(w http.ResponseWriter, r *http.Request) {
	slog.Info("request details", "method", r.Method, "path", r.URL.Path, "from", r.RemoteAddr)
	if r.URL.Path == "/" {
		homepageTemplate.ExecuteTemplate(w, "base", nil)
		return
	}
	static.ServeHTTP(w, r)
}
