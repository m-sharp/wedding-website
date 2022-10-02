package web

import (
	"go.uber.org/zap"
	"net/http"
)

type Route struct {
	Path    string
	Method  string
	Handler func(w http.ResponseWriter, r *http.Request)
}

type ApiRouter struct {
	log    *zap.Logger
	routes []*Route
}

// ToDo - finish
func NewApiRouter(log *zap.Logger) *ApiRouter {
	inst := &ApiRouter{log: log.Named("ApiRoute")}
	inst.routes = []*Route{
		{
			Path:    "/rsvp",
			Method:  http.MethodPost,
			Handler: inst.RSVPPost,
		},
	}
	return inst
}

func (a *ApiRouter) SetupRoutes() {
	//for _, route := range a.routes {
	//	http.HandleFunc()
	//}
}

func (a *ApiRouter) RSVPPost(w http.ResponseWriter, r *http.Request) {

}
