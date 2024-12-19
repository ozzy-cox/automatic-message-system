package consumer

import (
	"encoding/json"
)

type MessageResponse struct {
	Message   *string `json:"message"`
	MessageId *string `json:"messageId"`
}

func (i MessageResponse) MarshalBinary() ([]byte, error) {
	return json.Marshal(i)
}
