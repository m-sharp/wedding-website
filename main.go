package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"go.uber.org/zap"

	"github.com/m-sharp/wedding-website/lib"
	"github.com/m-sharp/wedding-website/lib/migrations"
)

const (
	Port         = "8081"
	TemplatesDir = "templates"
)

type RenderContext struct {
	TargetDate string
	TargetYear int
}

var (
	// Rendering context variables
	pageContext = &RenderContext{
		TargetYear: 2023,
		TargetDate: "10.7.23",
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := getCfg()
	logger := getLogger(cfg)
	client := getDBClient(ctx, cfg, logger)

	if err := migrations.RunAll(ctx, client, logger); err != nil {
		log.Fatal("Failed to run DB migrations", zap.Error(err))
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
	http.HandleFunc("/", handlePageRequests(logger, layoutTplPath, pageContext))

	// Start the web server
	logger.Info("Now listening!", zap.String("Port", Port))
	if err := http.ListenAndServe(fmt.Sprintf(":%v", Port), nil); err != nil {
		logger.Panic("Server stopped listening", zap.Error(err))
	}
}

func getCfg() *lib.Config {
	cfg, err := lib.NewConfig()
	if err != nil {
		log.Fatalf("Error creating Config: %s", err.Error())
	}

	return cfg
}

func getLogger(cfg *lib.Config) *zap.Logger {
	dev, err := cfg.Get(lib.Development)
	if err != nil {
		dev = "false"
	}

	var logger *zap.Logger
	if dev == "true" {
		logger, err = zap.NewDevelopment()

	} else {
		logger, err = zap.NewProduction()
	}
	if err != nil {
		log.Fatalf("Error creating Logger: %s", err.Error())
	}

	return logger
}

func getDBClient(ctx context.Context, cfg *lib.Config, log *zap.Logger) *lib.DBClient {
	client, err := lib.NewDBClient(ctx, cfg, log)
	if err != nil {
		log.Fatal("Error creating DB client", zap.Error(err))
	}
	if err = client.CheckConnection(); err != nil {
		log.Fatal("DB connection check failed", zap.Error(err))
	}

	return client
}

func handlePageRequests(
	log *zap.Logger,
	layoutTplPath string,
	pageContext *RenderContext,
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		targetTplPath := getTargetTplPath(r.URL.Path)
		log = log.With(zap.String("Path", targetTplPath))

		if is404(targetTplPath) {
			log.Info("404 Request")
			http.NotFound(w, r)
			return
		}

		tmpl, err := template.ParseFiles(layoutTplPath, filepath.Join(TemplatesDir, "nav"), targetTplPath)
		if err != nil {
			log.Error("Error parsing template files", zap.Error(err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.ExecuteTemplate(w, "layout", pageContext); err != nil {
			log.Error("Error executing template files", zap.Error(err))
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
