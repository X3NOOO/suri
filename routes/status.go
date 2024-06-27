package routes

import "net/http"

func (ctx *RoutingContext) StatusALL(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
