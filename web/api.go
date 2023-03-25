package web

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/m-sharp/wedding-website/lib"
	"go.uber.org/zap"
)

type Route struct {
	Path    string
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}

type ApiRouter struct {
	ctx          context.Context // ToDo - Context probably should be held here...
	log          *zap.Logger
	rsvpProvider *lib.RSVPProvider
	routes       []*Route
}

func NewApiRouter(ctx context.Context, log *zap.Logger, client *lib.DBClient) *ApiRouter {
	inst := &ApiRouter{
		ctx:          ctx,
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
			Path:    "/api/rsvp",
			Method:  http.MethodGet,
			Handler: inst.RSVPGetAll,
		},
	}
	return inst
}

func (a *ApiRouter) SetupRoutes(router *mux.Router) {
	for _, route := range a.routes {
		router.HandleFunc(route.Path, route.Handler).Methods(route.Method)
	}
}

func (a *ApiRouter) RSVPCreate(w http.ResponseWriter, r *http.Request) {
	// ToDo: RSVP needs to be constructed from body data, verified, and inserted
}

// ToDo: This route should be locked down...
func (a *ApiRouter) RSVPGetAll(w http.ResponseWriter, _ *http.Request) {
	rsvps, err := a.rsvpProvider.GetAll(a.ctx)
	if err != nil {
		a.log.Error("Failed to get RSVP records", zap.Error(err))
		http.Error(w, "failed to get RSVP records", http.StatusInternalServerError)
	}

	respData, err := json.Marshal(rsvps)
	if err != nil {
		a.log.Error("Failed to marshal RSVP response", zap.Any("Records", rsvps), zap.Error(err))
		http.Error(w, "failed to marshal response", http.StatusInternalServerError)
	}

	if _, err := w.Write(respData); err != nil {
		a.log.Error("Failed to write response", zap.ByteString("ResponseContent", respData), zap.Error(err))
		return
	}
}
