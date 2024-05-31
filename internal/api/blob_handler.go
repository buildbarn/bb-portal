package api

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/ent/gen/ent/blob"
	"github.com/buildbarn/bb-portal/pkg/cas"
)

type blobHandler struct {
	client     *ent.Client
	casManager *cas.ConnectionManager
}

func NewBlobHandler(client *ent.Client, casManager *cas.ConnectionManager) http.Handler {
	return &blobHandler{client: client, casManager: casManager}
}

func (b *blobHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	blobIDPathValue := request.PathValue("blobID")
	name := request.PathValue("name")
	blobID, err := strconv.Atoi(blobIDPathValue)
	if err != nil {
		writeErr(writer, request, http.StatusBadRequest, fmt.Sprintf("Invalid blobID: %s", blobIDPathValue))
		return
	}

	// TODO: We probably want semantic IDs, not row IDs.
	blobRecord, err := b.client.Blob.Get(request.Context(), blobID)
	if err != nil {
		writeErr(
			writer,
			request,
			http.StatusNotFound,
			fmt.Sprintf("Could not find blob with blobID: %s", blobIDPathValue),
		)
		return
	}

	b.serveBlob(writer, request, name, blobRecord)
}

func (b *blobHandler) serveBlob(writer http.ResponseWriter, request *http.Request, name string, blobRecord *ent.Blob) {
	if blobRecord.ArchivingStatus == blob.ArchivingStatusSUCCESS {
		http.ServeFile(writer, request, blobRecord.ArchiveURL)
		return
	}

	// Fallback to reading original.
	uri, err := url.Parse(blobRecord.URI)
	if err != nil {
		writeErr(
			writer,
			request,
			http.StatusInternalServerError,
			fmt.Sprintf("Blob %d had an invalid URI: %s", blobRecord.ID, blobRecord.URI),
		)
		return
	}
	switch uri.Scheme {
	case "file":
		http.ServeFile(writer, request, uri.Path)
	case "bytestream":
		b.serveFromBytestream(writer, request, name, uri)
	default:
		writeErr(writer, request, http.StatusInternalServerError, fmt.Sprintf("unsupported URI scheme: %s", uri.Scheme))
	}
}

func (b *blobHandler) serveFromBytestream(writer http.ResponseWriter, request *http.Request, name string, uri *url.URL) {
	casClient, err := b.casManager.GetClientForURI(request.Context(), uri)
	if err != nil {
		writeErr(writer, request, http.StatusInternalServerError, err.Error())
		return
	}
	defer casClient.Close()

	tmpFile, err := os.CreateTemp("", filepath.Base(uri.Path))
	if err != nil {
		writeErr(writer, request, http.StatusInternalServerError, err.Error())
		return
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	err = casClient.ReadBlobToFile(request.Context(), uri, tmpFile.Name())
	if err != nil {
		writeErr(writer, request, http.StatusInternalServerError, err.Error())
		return
	}

	if _, err = tmpFile.Seek(0, io.SeekStart); err != nil {
		writeErr(writer, request, http.StatusInternalServerError, err.Error())
		return
	}
	http.ServeContent(writer, request, name, time.Time{}, tmpFile)
}

func writeErr(writer http.ResponseWriter, request *http.Request, statusCode int, msg string) {
	writer.WriteHeader(statusCode)
	if _, err := writer.Write([]byte(msg)); err != nil {
		slog.ErrorContext(
			request.Context(),
			"could not write response",
			"statusCode", statusCode, "msg", msg,
		)
	}
}
