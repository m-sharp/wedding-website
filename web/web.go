package web

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strconv"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/wedding-website/lib"
)

const (
	Port = "8081"
)

var (
	TemplatesDir  = filepath.FromSlash(filepath.Join("web", "templates"))
	TemplateFiles = []string{
		filepath.Join(TemplatesDir, "layout"),
		filepath.Join(TemplatesDir, "nav"),
		filepath.Join(TemplatesDir, "analytics"),
		filepath.Join(TemplatesDir, "typography"),
	}
	ReservedFiles = append(TemplateFiles, filepath.Join(TemplatesDir, "adminView"))

	StaticDir    = filepath.FromSlash(filepath.Join("web", "static"))
	SiteFilesDir = filepath.FromSlash(filepath.Join("web", "site_files"))
)

type RenderContext struct {
	TargetDate string
	TargetYear int
}

// ToDo: cache page responses somehow
type Server struct {
	cfg       *lib.Config
	log       *zap.Logger
	renderCtx *RenderContext
	router    *mux.Router
}

func NewWebServer(cfg *lib.Config, log *zap.Logger, renderCtx *RenderContext, api *ApiRouter) *Server {
	inst := &Server{
		cfg:       cfg,
		log:       log.Named("WebServer"),
		renderCtx: renderCtx,
		router:    mux.NewRouter(),
	}
	// Setup API routes first, as the web request routes have a generic catch all
	api.SetupRoutes(inst.router)
	// Setup web request routes
	inst.setupRoutes()
	return inst
}

func (s *Server) setupRoutes() {
	// Handle static asset requests
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(StaticDir))))

	// Handle special files
	s.router.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(SiteFilesDir, "robots.txt"))
	})
	s.router.HandleFunc("/sitemap.xml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(SiteFilesDir, "sitemap.xml"))
	})
	s.router.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, filepath.Join(SiteFilesDir, "favicon.ico"))
	})

	// Handle other requests
	s.router.PathPrefix("/").HandlerFunc(s.handlePageRequests)
}

func (s *Server) Serve() error {
	isDev, err := s.cfg.Get(lib.Development)
	if err != nil {
		s.log.Error("Failed to get Development value from config", zap.Error(err))
		isDev = "false"
	}
	development, err := strconv.ParseBool(isDev)
	if err != nil {
		s.log.Error("Failed to parse Development value from config as bool", zap.Error(err))
		development = false
	}

	csrfSecret, err := s.cfg.Get(lib.CSRFSecret)
	if err != nil {
		return fmt.Errorf("failed to get CSRF Secret from config: %w", err)
	}

	s.log.Info("Now listening!", zap.String("Port", Port))
	return http.ListenAndServe(
		fmt.Sprintf(":%v", Port),
		csrf.Protect([]byte(csrfSecret), csrf.Secure(!development))(s.router),
	)
}

func (s *Server) handlePageRequests(w http.ResponseWriter, r *http.Request) {
	targetPage := filepath.Clean(r.URL.Path)
	targetPage = filepath.Join(TemplatesDir, s.doRedirects(targetPage))
	log := s.log.With(zap.String("Path", targetPage))
	log.Debug("Requested template path")

	if s.is404(targetPage) {
		log.Info("404 Request")
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles(append(TemplateFiles, targetPage)...)
	if err != nil {
		log.Error("Error parsing template files", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"TargetDate":     s.renderCtx.TargetDate,
		"TargetYear":     s.renderCtx.TargetYear,
	}); err != nil {
		log.Error("Error building template files", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) doRedirects(urlPath string) string {
	switch urlPath {
	case "/":
		urlPath = "/home"
	}

	return urlPath
}

func (s *Server) is404(targetTplPath string) bool {
	info, err := os.Stat(targetTplPath)
	if err != nil && os.IsNotExist(err) {
		s.log.Debug("Could not find requested path", zap.String("Path", targetTplPath), zap.Error(err))
		return true
	}
	if info.IsDir() {
		return true
	}

	for _, p := range ReservedFiles {
		if targetTplPath == p {
			s.log.Debug("Reserved Path requested", zap.String("Path", targetTplPath))
			return true
		}
	}
	return false
}
