package system

import (
	"sync"

	"net/http"

	"github.com/gorilla/mux"
)

var routers = make([]*RouteManager, 0)
var routerLock sync.Mutex

type Route struct {
	prefix     string
	path       string
	host       string
	methods    []string
	queries    []string
	subrouters []*Router
	handler    http.Handler
}

func (r *Route) Subrouter() *Router {
	router := &Router{}
	r.subrouters = append(r.subrouters, router)
	return router
}

func (r *Route) Path(path string) *Route {
	r.path = path
	return r
}

func (r *Route) PathPrefix(prefix string) *Route {
	r.prefix = prefix
	return r
}

func (r *Route) Host(host string) *Route {
	r.host = host
	return r
}

func (r *Route) Methods(methods ...string) *Route {
	r.methods = append(r.methods, methods...)
	return r
}

func (r *Route) Queries(pairs ...string) *Route {
	r.queries = append(r.queries, pairs...)
	return r
}

func (r *Route) Handler(handler http.Handler) *Route {
	r.handler = handler
	return r
}

func (r *Route) HandlerFunc(f func(http.ResponseWriter, *http.Request)) *Route {
	return r.Handler(http.HandlerFunc(f))
}

type Router struct {
	routes []*Route
}

func (r *Router) NewRoute() *Route {
	route := &Route{}
	r.routes = append(r.routes, route)
	return route
}

func (r *Router) Path(path string) *Route {
	return r.NewRoute().Path(path)
}

func (r *Router) Host(host string) *Route {
	return r.NewRoute().Host(host)
}

func (r *Router) Methods(methods ...string) *Route {
	return r.NewRoute().Methods(methods...)
}

func (r *Router) Queries(pairs ...string) *Route {
	return r.NewRoute().Queries(pairs...)
}

func (r *Router) Handler(path string, handler http.Handler) *Route {
	return r.NewRoute().Path(path).Handler(handler)
}

func (r *Router) handlerFunc(path string, f func(http.ResponseWriter, *http.Request)) *Route {
	return r.NewRoute().Path(path).HandlerFunc(f)
}

func NewRouter() *Router {
	return &Router{}
}

type RouteManager struct {
	PathPrefix string
	router     *Router
	root       *Router
}

func NewRouteManager(path string) *RouteManager {
	rm := &RouteManager{PathPrefix: path}
	routerLock.Lock()
	defer routerLock.Unlock()
	routers = append(routers, rm)
	return rm
}

func (rm *RouteManager) Router() *Router {
	if rm.router == nil {
		rm.router = NewRouter()
	}
	return rm.router
}

func (rm *RouteManager) Root() *Router {
	if rm.root == nil {
		rm.root = NewRouter()
	}
	return rm.root
}

func attachRouter(target *mux.Router, router *Router) {
	if router != nil && len(router.routes) > 0 {
		for _, route := range router.routes {
			if route.handler != nil || len(route.subrouters) > 0 {
				r := target.NewRoute()
				if len(route.path) > 0 {
					r.Path(route.path)
				} else if len(route.prefix) > 0 {
					r.PathPrefix(route.prefix)
				}

				if len(route.host) > 0 {
					r.Host(route.host)
				}

				if len(route.methods) > 0 {
					r.Methods(route.methods...)
				}

				if len(route.queries) > 0 {
					r.Queries(route.queries...)
				}

				if route.handler != nil {
					r.Handler(route.handler)
				}

				if len(route.subrouters) > 0 {
					srr := r.Subrouter()
					for _, sr := range route.subrouters {
						attachRouter(srr, sr)
					}
				}
			}
		}
	}
}

func (rm *RouteManager) attachRouter(target *mux.Router) {
	if rm.router != nil && len(rm.router.routes) > 0 {
		attachRouter(target.PathPrefix(rm.PathPrefix).Subrouter(), rm.router)
	}
	if rm.root != nil && len(rm.root.routes) > 0 {
		attachRouter(target, rm.root)
	}
	return
}

func AttachRouterManager(router *mux.Router) {
	if len(routers) > 0 {
		for _, rm := range routers {
			rm.attachRouter(router)
		}
	}
}
