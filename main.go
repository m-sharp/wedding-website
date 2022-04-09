package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DEBUG        = true
	TemplatesDir = "templates"
)

type RenderContext struct {
	TargetYear int
}

func main() {
	// Declare rendering context variables
	pageContext := &RenderContext{
		TargetYear: 2023,
	}

	// Handle static asset requests
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle other requests
	layoutTplPath := filepath.Join(TemplatesDir, "layout")
	http.HandleFunc("/", handlePageRequests(layoutTplPath, pageContext))

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	log.Println("Listening on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Panic(fmt.Errorf("stopped listening with error: %w", err).Error())
	}
}

func handlePageRequests(layoutTplPath string, pageContext *RenderContext) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetTplPath := getTargetTplPath(r.URL.Path)

		if is404(targetTplPath) {
			if DEBUG {
				log.Println(fmt.Sprintf("404: %q", targetTplPath))
			}
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(layoutTplPath, targetTplPath)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if DEBUG {
			log.Println(fmt.Sprintf("Serving %q (resolved to %q)", r.URL.Path, targetTplPath))
		}
		if err := tmpl.ExecuteTemplate(w, "layout", pageContext); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func getTargetTplPath(urlPath string) string {
	cleanedPath := filepath.Clean(urlPath)
	if cleanedPath == "\\" {
		cleanedPath = "\\home"
	}
	return filepath.Join(TemplatesDir, cleanedPath)
}

func is404(targetTplPath string) bool {
	info, err := os.Stat(targetTplPath)
	if err != nil && os.IsNotExist(err) {
		return true
	}
	if info.IsDir() {
		return true
	}
	return false
}
