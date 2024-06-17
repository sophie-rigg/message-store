package messages

import (
	"encoding/json"
)

type postResponse struct {
	ID int64 `json:"id"`
}

func newPostResponse(id int64) *postResponse {
	return &postResponse{
		ID: id,
	}
}

// MarshalToJson marshall the response to json
func (r *postResponse) MarshalToJson() ([]byte, error) {
	return json.Marshal(r)
}
