package routes

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/X3NOOO/suri/models"
)

func (ctx *RoutingContext) QueryPOST(w http.ResponseWriter, r *http.Request) {
	body_json, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var body models.QueryPOSTRequest

	err = json.Unmarshal(body_json, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	llm_response, err := ctx.AI.Query(body.Query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := models.QueryPOSTResponse{
		Response: llm_response,
	}

	response_json, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(response_json)
}
