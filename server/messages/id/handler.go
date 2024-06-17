package id

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sophie-rigg/message-store/cache"
)

type handler struct {
	localCache *cache.Cache
	logger     zerolog.Logger
}

var (
	_ http.Handler = (*handler)(nil)

	_errNoIDProvided = errors.New("no id provided in query")
)

func NewHandler(localCache *cache.Cache) http.Handler {
	return &handler{
		localCache: localCache,
		logger: log.With().Fields(map[string]interface{}{
			"handler": "messages/id",
		}).Logger(),
	}
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	// Only allow GET requests
	switch request.Method {
	case http.MethodGet:
		h.handleGet(writer, request)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) handleGet(writer http.ResponseWriter, request *http.Request) {
	id, err := getIDFromRequest(request)
	if err != nil {
		h.logger.Error().Err(err).Msg("error getting id from request")
		http.Error(writer, "Error getting id from request", http.StatusBadRequest)
		return
	}

	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		h.logger.Error().Err(err).Msg("error converting id to int64")
		http.Error(writer, fmt.Sprintf("Error converting id: %s to int64", id), http.StatusBadRequest)
		return
	}

	// Find message in cache
	message, ok := h.localCache.Get(ID)
	if !ok {
		// If message not found, return 404
		// TODO: Should check permanent storage once introduced
		http.Error(writer, fmt.Sprintf("Message with ID: %d not found", ID), http.StatusNotFound)
		return
	}

	// Write message to response
	var response bytes.Buffer
	_, err = response.WriteString(message)
	if err != nil {
		h.logger.Error().Err(err).Msg("error writing message to response")
		http.Error(writer, "Error writing message to response", http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(response.Bytes())
	if err != nil {
		h.logger.Error().Err(err).Msg("error writing response")
		http.Error(writer, "Error writing response", http.StatusInternalServerError)
		return
	}
}

// getIDFromRequest gets the id from the request vars
func getIDFromRequest(r *http.Request) (string, error) {
	queryVars := mux.Vars(r)
	id, ok := queryVars["id"]
	if !ok {
		return "", _errNoIDProvided
	}
	return id, nil
}
