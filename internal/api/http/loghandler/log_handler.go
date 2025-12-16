package loghandler

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/buildbarn/bb-portal/ent/gen/ent"
	"github.com/buildbarn/bb-portal/internal/database/dbauthservice"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/klauspost/compress/zstd"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// LogHandler serves build logs over HTTP.
type LogHandler struct {
	client        *ent.Client
	dbAuthService *dbauthservice.DbAuthService
	tracer        trace.Tracer
}

// NewLogHandler creates a new LogHandler.
func NewLogHandler(client *ent.Client, dbAuthService *dbauthservice.DbAuthService, tracerProvider trace.TracerProvider) (*LogHandler, error) {
	return &LogHandler{
		client:        client,
		dbAuthService: dbAuthService,
		tracer:        tracerProvider.Tracer("github.com/buildbarn/bb-portal/internal/api/http/loghandler"),
	}, nil
}

// ServeHTTP serves the log as a memory efficient http stream.
func (h *LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := dbauthservice.NewContextWithDbAuthService(r.Context(), h.dbAuthService)
	ctx, span := h.tracer.Start(ctx, "LogHandler.ServeLog")
	defer span.End()

	errorOut := func(err error, message string, code int) {
		http.Error(w, message, code)
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}

	vars := mux.Vars(r)

	invocationID, err := uuid.Parse(vars["invocation_id"])
	if err != nil {
		errorOut(err, "Invalid Invocation Id", http.StatusBadRequest)
		return
	}

	parseParam := func(val string, defaultVal int64) int64 {
		if val == "" {
			return defaultVal
		}
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return defaultVal
		}
		return i
	}

	query := r.URL.Query()
	startLine := parseParam(query.Get("start_line"), 0)
	if startLine < 0 {
		errorOut(fmt.Errorf("Invalid start line: %d", startLine), "Invalid start line", http.StatusBadRequest)
		return
	}
	endLine := parseParam(query.Get("end_line"), -1)
	if endLine == -1 {
		endLine = math.MaxInt64
	}

	lines, err := h.getLogLineReader(ctx, invocationID, startLine, endLine)
	if err != nil {
		if status.Code(err) == grpc_codes.NotFound {
			errorOut(err, "Log not found", http.StatusNotFound)
		}
		errorOut(err, "Error getting log lines", http.StatusInternalServerError)
		return
	}

	var writer io.Writer = w
	if strings.Contains(r.Header.Get("Accept-Encoding"), "zstd") {
		if encoder, err := zstd.NewWriter(w); err == nil {
			w.Header().Set("Content-Encoding", "zstd")
			defer encoder.Close()
			writer = encoder
		}
	}

	// Set the header of the return result, we don't need to set the
	// status code as that is implied to be 200 if we start writing a
	// response.
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if _, err := io.Copy(writer, lines); err != nil {
		// At this point returning an error code is too late but we will
		// mark the error in the otel metrics.
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}
