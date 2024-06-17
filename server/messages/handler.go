package messages

import (
	"fmt"
	"io"
	"net/http"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sophie-rigg/message-store/cache"
)

type handler struct {
	localCache *cache.Cache
	logger     zerolog.Logger
}

// handler must implement the http.Handler interface
var _ http.Handler = (*handler)(nil)

func NewHandler(localCache *cache.Cache) http.Handler {
	return &handler{
		localCache: localCache,
		logger: log.With().Fields(map[string]interface{}{
			"handler": "messages",
		}).Logger(),
	}
}

func (h *handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.Method {
	case http.MethodPost:
		h.handlePost(writer, request)
	default:
		http.Error(writer, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *handler) handlePost(writer http.ResponseWriter, request *http.Request) {
	data, err := io.ReadAll(request.Body)
	if err != nil {
		h.logger.Error().Err(err).Msg("error reading request body")
		http.Error(writer, "Error reading request body", http.StatusInternalServerError)
		return
	}

	// Add the data to the cache
	// TODO: Add permanent storage
	ID := h.localCache.Add(string(data))

	response, err := newPostResponse(ID).MarshalToJson()
	if err != nil {
		http.Error(writer, fmt.Sprintf("Error creating response, id: %d", ID), http.StatusInternalServerError)
		return
	}

	_, err = writer.Write(response)
	if err != nil {
		h.logger.Error().Err(err).Msg("error writing response")
		http.Error(writer, "Error writing response", http.StatusInternalServerError)
		return
	}
}
