package server

import (
	"github.com/gorilla/mux"
	"github.com/sophie-rigg/message-store/cache"
	"github.com/sophie-rigg/message-store/server/messages"
	"github.com/sophie-rigg/message-store/server/messages/id"
)

// Register registers the API routes for the server
func Register(cache *cache.Cache) *mux.Router {
	router := mux.NewRouter()

	router.Path("/messages").Handler(messages.NewHandler(cache))
	// id is the message id
	router.Path("/messages/{id}").Handler(id.NewHandler(cache))

	return router
}
