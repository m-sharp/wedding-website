package web

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"

	"github.com/gorilla/csrf"
	"github.com/gorilla/mux"
	"go.uber.org/zap"

	"github.com/m-sharp/wedding-website/lib"
)

type MiddlewareFunc func(nextHandler http.HandlerFunc) http.HandlerFunc

type Route struct {
	Path       string
	Method     string
	Handler    http.HandlerFunc
	MiddleWare MiddlewareFunc
}

type ApiRouter struct {
	cfg             *lib.Config
	recaptchaSecret string
	webAdminUser    string
	webAdminPass    string

	log          *zap.Logger
	rsvpProvider *lib.RSVPProvider
	routes       []*Route
}

func NewApiRouter(cfg *lib.Config, log *zap.Logger, client *lib.DBClient) (*ApiRouter, error) {
	recaptchaSecret, err := cfg.Get(lib.RecaptchaSecret)
	if err != nil {
		return nil, fmt.Errorf("failed to get recaptcha secret from config: %w", err)
	}
	adminUser, err := cfg.Get(lib.WebAdminUser)
	if err != nil {
		return nil, fmt.Errorf("failed to get web admin username from config: %w", err)
	}
	adminPass, err := cfg.Get(lib.WebAdminPass)
	if err != nil {
		return nil, fmt.Errorf("failed to get web admin pass from config: %w", err)
	}

	inst := &ApiRouter{
		cfg:             cfg,
		recaptchaSecret: recaptchaSecret,
		webAdminPass:    adminPass,
		webAdminUser:    adminUser,
		log:             log.Named("ApiRoute"),
		rsvpProvider:    lib.NewRSVPProvider(client),
	}
	inst.routes = []*Route{
		{
			Path:    "/api/rsvp",
			Method:  http.MethodPost,
			Handler: inst.RSVPCreate,
		},
		{
			Path:       "/api/rsvp",
			Method:     http.MethodGet,
			Handler:    inst.RSVPGetAll,
			MiddleWare: inst.BasicAuthMiddleware,
		},
	}
	return inst, nil
}

func (a *ApiRouter) SetupRoutes(router *mux.Router) {
	for _, route := range a.routes {
		if route.MiddleWare != nil {
			router.HandleFunc(route.Path, route.MiddleWare(route.Handler)).Methods(route.Method)
		} else {
			router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
		}
	}
}

func (a *ApiRouter) RSVPCreate(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	body, err := io.ReadAll(r.Body)
	if err != nil {
		a.log.Error("Failed to read RSVP post body", zap.Error(err))
		http.Error(w, "failed to read request body", http.StatusBadRequest)
		return
	}

	var rsvp lib.RSVP
	if err := json.Unmarshal(body, &rsvp); err != nil {
		a.log.Error("Failed to unmarshal RSVP body", zap.Error(err))
		http.Error(w, "failed to decode request body", http.StatusBadRequest)
		return
	}
	log := a.log.With(zap.Any("RSVP", rsvp))

	if err := rsvp.Validate(); err != nil {
		log.Error("RSVP failed validation", zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	respToken := r.URL.Query().Get("token")
	if respToken == "" {
		log.Error("RSVP post did not include verification token")
		http.Error(w, "request did not include recaptcha token", http.StatusBadRequest)
		return
	}

	verified, err := lib.Verify(log, a.recaptchaSecret, respToken, r.RemoteAddr)
	if err != nil {
		log.Error("Error verifying Recaptcha response", zap.Error(err))
		http.Error(w, "failed to verify recaptcha token", http.StatusInternalServerError)
		return
	} else if !verified {
		log.Error("Recaptcha response token not valid")
		http.Error(w, "recaptcha token failed validation", http.StatusBadRequest)
		return
	}

	log.Info("Saving RSVP record")
	if err := a.rsvpProvider.Add(ctx, &rsvp); err != nil {
		log.Error("Failed to add RSVP record", zap.Error(err))
		http.Error(w, "failed to save RSVP record", http.StatusInternalServerError)
		return
	}

	// ToDo: Compile template for email, Send off email to RSVPer
	// TODo: Compile template for email, Send off to us

	w.WriteHeader(http.StatusCreated)
}

func (a *ApiRouter) RSVPGetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	rsvps, err := a.rsvpProvider.GetAll(ctx)
	if err != nil {
		a.log.Error("Failed to get RSVP records", zap.Error(err))
		http.Error(w, "failed to get RSVP records", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles(append(TemplateFiles, filepath.Join(TemplatesDir, "adminView"))...)
	if err != nil {
		a.log.Error("Error parsing template files", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.ExecuteTemplate(w, "layout", map[string]interface{}{
		csrf.TemplateTag: csrf.TemplateField(r),
		"RSVPs":          rsvps,
	}); err != nil {
		a.log.Error("Error building template files", zap.Error(err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// BasicAuthMiddleware wraps a http.HandlerFunc in a Basic Auth check.
// ToDo: Could probably live in a middleware only file
func (a *ApiRouter) BasicAuthMiddleware(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			if subtleCompare(username, a.webAdminUser) && subtleCompare(password, a.webAdminPass) {
				nextHandler(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// subtleCompare compares two strings in constant time by first hashing the strings.
func subtleCompare(actual, expected string) bool {
	actualHash := sha256.Sum256([]byte(actual))
	expectedHash := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(actualHash[:], expectedHash[:]) == 1
}
