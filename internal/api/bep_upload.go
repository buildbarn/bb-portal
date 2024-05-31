package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/pkg/processing"
)

const (
	MB            = 1024 * 1024
	MaxUploadSize = 500 * MB
)

type bepUploadHandler struct {
	client       *ent.Client
	blobArchiver processing.BlobMultiArchiver
}

func NewBEPUploadHandler(client *ent.Client, blobArchiver processing.BlobMultiArchiver) http.Handler {
	return &bepUploadHandler{
		client:       client,
		blobArchiver: blobArchiver,
	}
}

func (b bepUploadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, MaxUploadSize)
	if err := r.ParseMultipartForm(MaxUploadSize); err != nil {
		msg := fmt.Sprintf("The uploaded file is too big. Please choose an file that's less than %dMB in size", MaxUploadSize/MB)
		http.Error(w, msg, http.StatusBadRequest)
		return
	}

	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer file.Close()

	slog.Info("Received file", "name", fileHeader.Filename, "size", fileHeader.Size)

	tmpFile, err := os.CreateTemp("", fileHeader.Filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = io.Copy(tmpFile, file)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer os.Remove(tmpFile.Name())

	workflow := processing.New(b.client, b.blobArchiver)
	invocation, err := workflow.ProcessFile(r.Context(), tmpFile.Name())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	location := fmt.Sprintf("/bazel-invocations/%s", invocation.InvocationID)
	// NOTE: Want to do http.Redirect(w, r, location, http.StatusSeeOther), but can't get it working with antd Upload widget.
	writeLocationResponse(w, location)
}

func writeLocationResponse(w http.ResponseWriter, location string) {
	w.WriteHeader(http.StatusOK)
	resp := struct {
		Location string
	}{
		Location: location,
	}
	respBody, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, err = w.Write(respBody)
	if err != nil {
		slog.Error("failed to write response", "err", err)
	}
}
