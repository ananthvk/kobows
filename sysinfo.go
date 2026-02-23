package main

import (
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"os/user"
	"runtime"
	"time"
)

var sysinfoTemplate = template.Must(template.Must(baseTemplate.Clone()).ParseFiles("templates/sysinfo.gohtml"))

type sysinfoData struct {
	Time          time.Time
	Hostname      string
	OS            string
	Arch          string
	NumCPU        int
	GoVersion     string
	MemStats      runtime.MemStats
	NumGoroutines int
	User          string
	WorkingDir    string
	Uptime        time.Duration
	PID           int
}

func sysinfoHandler(w http.ResponseWriter, r *http.Request) {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	slog.Info("request details", "method", r.Method, "path", r.URL.Path, "from", r.RemoteAddr)

	data := sysinfoData{
		Time:          time.Now(),
		Hostname:      getHostname(),
		OS:            runtime.GOOS,
		Arch:          runtime.GOARCH,
		NumCPU:        runtime.NumCPU(),
		GoVersion:     runtime.Version(),
		MemStats:      m,
		NumGoroutines: runtime.NumGoroutine(),
		User:          getUser(),
		WorkingDir:    getWorkingDir(),
		Uptime:        time.Since(startTime),
		PID:           os.Getpid(),
	}
	sysinfoTemplate.ExecuteTemplate(w, "base", data)
}

func getHostname() string {
	hostname, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return hostname
}

func getUser() string {
	user, err := user.Current()
	if err != nil {
		return "unknown"
	}
	return user.Name
}

func getWorkingDir() string {
	wd, err := os.Getwd()
	if err != nil {
		return "unknown"
	}
	return wd
}
