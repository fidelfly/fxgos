package resources

import (
	"github.com/lyismydg/fxgos/system"
)

/*var handlers  = make(map[string]http.Handler, 0)
var handlerFunctions = make(map[string]http.HandlerFunc, 0)

func defineResourceHandler(method string, path string, handler http.Handler)  {
	if len(method) > 0 {
		handlers[method + ":" + path] = handler
	} else {
		handlers[path] = handler
	}
}

func defineResourceHandlerFunction(method string, path string, handlerFunction http.HandlerFunc) {
	if len(method) > 0 {
		handlerFunctions[method + ":" + path] = handlerFunction
	} else {
		handlerFunctions[path] = handlerFunction
	}
}

func resolvePathMethod(key string) (path string, method string) {
	var splitIndex = strings.Index(key, ":");
	if splitIndex > 0 {
		method = key[:splitIndex]
		path = key[splitIndex + 1:]
	} else {
		path = key
	}
	return;
}

func SetupRouter(router *mux.Router) {
	for key, handler := range handlers {
		path, method := resolvePathMethod(key)
		if len(method)  > 0 {
			router.Handle(path, handler).Methods(method)
		} else {
			router.Handle(path, handler)
		}
	}

	for key, handlerFunction := range handlerFunctions {
		path, method := resolvePathMethod(key)
		if len(method)  > 0 {
			router.HandleFunc(path, handlerFunction).Methods(method)
		} else {
			router.HandleFunc(path, handlerFunction)
		}
	}
}*/

var myRouter = system.NewRouteManager("/resource")


