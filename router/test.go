package router

import (
	"net/http"
		"github.com/gorilla/mux"
)

func Test(w http.ResponseWriter, r *http.Request)  {
	muxVars := mux.Vars(r)
	textKey := muxVars["key"]
	switch textKey {
	case "createGroup":
		return
	default:
		return
	}
}
