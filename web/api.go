package web

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m-sharp/wedding-website/lib"
	"go.uber.org/zap"
)

type MiddlewareFunc func(nextHandler http.HandlerFunc) http.HandlerFunc

type Route struct {
	Path       string
	Method     string
	Handler    http.HandlerFunc
	MiddleWare MiddlewareFunc
}

type ApiRouter struct {
	ctx          context.Context // ToDo - Context probably shouldn't be held here...
	cfg          *lib.Config
	log          *zap.Logger
	rsvpProvider *lib.RSVPProvider
	routes       []*Route
}

func NewApiRouter(ctx context.Context, cfg *lib.Config, log *zap.Logger, client *lib.DBClient) *ApiRouter {
	inst := &ApiRouter{
		ctx:          ctx,
		cfg:          cfg,
		log:          log.Named("ApiRoute"),
		rsvpProvider: lib.NewRSVPProvider(client),
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
	return inst
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

	if err := rsvp.Validate(); err != nil {
		a.log.Error("RSVP failed validation", zap.Any("RSVP", rsvp), zap.Error(err))
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := a.rsvpProvider.Add(a.ctx, &rsvp); err != nil {
		a.log.Error("Failed to add RSVP record", zap.Any("RSVP", rsvp), zap.Error(err))
		http.Error(w, "failed to save RSVP record", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (a *ApiRouter) RSVPGetAll(w http.ResponseWriter, _ *http.Request) {
	rsvps, err := a.rsvpProvider.GetAll(a.ctx)
	if err != nil {
		a.log.Error("Failed to get RSVP records", zap.Error(err))
		http.Error(w, "failed to get RSVP records", http.StatusInternalServerError)
		return
	}

	respData, err := json.Marshal(rsvps)
	if err != nil {
		a.log.Error("Failed to marshal RSVP response", zap.Any("Records", rsvps), zap.Error(err))
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(respData); err != nil {
		a.log.Error("Failed to write response", zap.ByteString("ResponseContent", respData), zap.Error(err))
		http.Error(w, "failed to write response", http.StatusInternalServerError)
	}
}

// BasicAuthMiddleware wraps a http.HandlerFunc in a Basic Auth check.
// ToDo: Could probably live in a middleware only file
func (a *ApiRouter) BasicAuthMiddleware(nextHandler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			adminUser, err := a.cfg.Get(lib.WebAdminUser)
			if err != nil {
				http.Error(w, "cannot determine web admin username", http.StatusInternalServerError)
				return
			}
			adminPass, err := a.cfg.Get(lib.WebAdminPass)
			if err != nil {
				http.Error(w, "cannot determine web admin password", http.StatusInternalServerError)
				return
			}

			if subtleCompare(username, adminUser) && subtleCompare(password, adminPass) {
				nextHandler(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}

// subtleCompare compares to strings in constant time by first hashing the strings.
func subtleCompare(actual, expected string) bool {
	actualHash := sha256.Sum256([]byte(actual))
	expectedHash := sha256.Sum256([]byte(expected))
	return subtle.ConstantTimeCompare(actualHash[:], expectedHash[:]) == 1
}
