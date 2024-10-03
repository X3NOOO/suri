package models

type QueryPOSTRequest struct {
	Query string `json:"query"`
}

type QueryPOSTResponse struct {
	Response string `json:"response"`
}
