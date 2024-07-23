package routes

import (
	"crypto/sha256"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/henomis/lingoose/document"
)

func (ctx *RoutingContext) KnowladgeGET(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (ctx *RoutingContext) knowladgePOSTJSON(w http.ResponseWriter, r *http.Request) {
	body_json, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var body document.Document

	err = json.Unmarshal(body_json, &body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = ctx.AI.Learn(body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctx *RoutingContext) knowladgePOSTFile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(ctx.MaxMemory)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	// The server should only be accessed by an authenticated client.
	// Normally we wouldn't trust the client with providing the (valid) header and would check it ourselves,
	// but because the http.DetectContentType lacks some of the types we might want to support, we don't really have any other choice.
	//
	// Besides, if someone already has your API key and wants to mess with you you're fucked either way lol,
	// spoofing the Content-Type header is the least they can do.
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		http.Error(w, "Content-Type header is missing", http.StatusBadRequest)
		return
	}

	hash := sha256.New()
	_, err = io.Copy(hash, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	metadata := map[string]any{
		"filename":    header.Filename,
		"upload_date": time.Now().Format(time.RFC3339),
		"sha256sum":   hash.Sum(nil),
	}

	fileContent, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = ctx.AI.LearnFile(fileContent, contentType, metadata)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (ctx *RoutingContext) KnowladgePOST(w http.ResponseWriter, r *http.Request) {
	contentType := strings.Split(strings.ToLower(r.Header.Get("Content-Type")), ";")[0]

	switch contentType {
	case "application/json":
		ctx.knowladgePOSTJSON(w, r)
	case "multipart/form-data":
		ctx.knowladgePOSTFile(w, r)
	default:
		http.Error(w, "Unsupported Content-Type", http.StatusBadRequest)
	}
}

func (ctx *RoutingContext) KnowladgeNameGET(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (ctx *RoutingContext) KnowladgeNamePUT(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}

func (ctx *RoutingContext) KnowladgeNameDELETE(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Not Implemented", http.StatusNotImplemented)
}
