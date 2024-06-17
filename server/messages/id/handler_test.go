package id

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sophie-rigg/message-store/cache"
)

func Test_handler_ServeHTTP(t *testing.T) {
	type fields struct {
		localCache func() *cache.Cache
	}
	type args struct {
		writer  *httptest.ResponseRecorder
		request func() *http.Request
	}
	tests := []struct {
		name         string
		fields       fields
		args         args
		expectedCode int
		expectedBody string
	}{
		{
			name: "Happy path",
			fields: fields{
				localCache: func() *cache.Cache {
					c, err := cache.NewCache()
					if err != nil {
						t.Errorf("Error creating cache: %v", err)
					}
					c.Add("test message")
					return c
				},
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					r := httptest.NewRequest("GET", "/messages/1", nil)
					return mux.SetURLVars(r, map[string]string{"id": "1"})
				},
			},
			expectedCode: 200,
			expectedBody: "test message",
		},
		{
			name: "bad id",
			fields: fields{
				localCache: func() *cache.Cache {
					c, err := cache.NewCache()
					if err != nil {
						t.Errorf("Error creating cache: %v", err)
					}
					c.Add("test message")
					return c
				},
			},
			args: args{
				writer: httptest.NewRecorder(),
				request: func() *http.Request {
					r := httptest.NewRequest("GET", "/messages/hello", nil)
					return mux.SetURLVars(r, map[string]string{"id": "hello"})
				},
			},
			expectedCode: 400,
			expectedBody: "Error converting id: hello to int64",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewHandler(tt.fields.localCache())

			h.ServeHTTP(tt.args.writer, tt.args.request())

			if tt.args.writer.Code != tt.expectedCode {
				t.Errorf("expected status code 200, got %d", tt.args.writer.Code)
			}

			if reflect.DeepEqual(tt.args.writer.Body.Bytes(), []byte(tt.expectedBody)) && tt.args.writer.Body.String() != tt.expectedBody {
				t.Errorf("expected message \n%s\ngot\n%s", tt.expectedBody, tt.args.writer.Body.String())
			}
		})
	}
}
