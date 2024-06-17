package server

import (
	"bytes"
	"net/http/httptest"
	"testing"

	"github.com/sophie-rigg/message-store/cache"
)

func TestRegister(t *testing.T) {
	t.Run("server test", func(t *testing.T) {
		c, err := cache.NewCache()
		if err != nil {
			t.Fatalf("error creating cache: %v", err)
		}
		router := Register(c)

		// Create a request to test the server
		buf := bytes.NewBufferString("test message")

		request := httptest.NewRequest("POST", "/messages", buf)
		response := httptest.NewRecorder()

		router.ServeHTTP(response, request)

		if response.Code != 200 {
			t.Fatalf("expected status code 200, got %d", response.Code)
		}

		getMessageRequest := httptest.NewRequest("GET", "/messages/1", nil)
		getMessageResponse := httptest.NewRecorder()

		router.ServeHTTP(getMessageResponse, getMessageRequest)

		if getMessageResponse.Code != 200 {
			t.Fatalf("expected status code 200, got %d", getMessageResponse.Code)
		}

		msg := getMessageResponse.Body.String()
		if msg != "test message" {
			t.Fatalf("expected message 'test message', got %s", msg)
		}
	})
}
