package models

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type EnvelopeQueue struct {
	Action string          `json:"action"`
	Data   json.RawMessage `json:"data"`
}

func NewEnvelopeQueue() *EnvelopeQueue {
	return &EnvelopeQueue{}
}

func (e *EnvelopeQueue) Convert(data []byte) (string, []byte, error) {
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()

	if err := dec.Decode(e); err != nil {
		return "", nil, err
	}
	if e.Action == "" {
		return "", nil, fmt.Errorf("empty action")
	}
	if len(e.Data) == 0 {
		return e.Action, nil, fmt.Errorf("empty data")
	}

	return e.Action, e.Data, nil
}

func (e *EnvelopeQueue) ConvertData(action, data string) []byte {
	raw := []byte(data)
	// если data — валидный JSON, вставим как RawMessage
	if json.Valid(raw) {
		payload := EnvelopeQueue{
			Action: action,
			Data:   json.RawMessage(raw),
		}
		b, _ := json.Marshal(payload)
		return b
	}

	// иначе положим как строку
	payload := struct {
		Action string `json:"action"`
		Data   string `json:"data"`
	}{
		Action: action,
		Data:   data,
	}
	b, _ := json.Marshal(payload)
	return b
}
