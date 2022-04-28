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
	DEBUG        = false
	TemplatesDir = "templates"
)

type RenderContext struct {
	TargetDate string
	TargetYear int
}

func main() {
	// Declare rendering context variables
	pageContext := &RenderContext{
		TargetYear: 2023,
		TargetDate: "??.??.23",
	}

	// Handle static asset requests
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Handle special files
	http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./site_files/robots.txt")
	})
	http.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./site_files/sitemap.xml")
	})

	// Handle other requests
	layoutTplPath := filepath.Join(TemplatesDir, "layout")
	http.HandleFunc("/", handlePageRequests(layoutTplPath, pageContext))

	//Start the web server, set the port to listen to 8080. Without a path it assumes localhost
	//Print any errors from starting the webserver using fmt
	log.Println("Listening on port 8081")
	if err := http.ListenAndServe(":8081", nil); err != nil {
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

		tmpl, err := template.ParseFiles(layoutTplPath, filepath.Join(TemplatesDir, "nav"), targetTplPath)
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
	switch cleanedPath {
	case "\\":
		cleanedPath = "\\home"
	case "/":
		cleanedPath = "/home"
	default:
	}
	return filepath.Join(TemplatesDir, cleanedPath)
}

var reservedPaths = []string{
	"templates\\layout",
	"templates/layout",
	"templates\\nav",
	"templates/nav",
}

func is404(targetTplPath string) bool {
	info, err := os.Stat(targetTplPath)
	if err != nil && os.IsNotExist(err) {
		return true
	}
	if info.IsDir() {
		return true
	}

	for _, p := range reservedPaths {
		if targetTplPath == p {
			return true
		}
	}
	return false
}
