package json

import (
	"encoding/json"
	"fmt"

	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := json.Unmarshal(input, &redirect); err != nil {
		return nil, fmt.Errorf("JSON Serializer: Decode %v -> %w", redirect, err)
	}

	return redirect, nil
}

func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)

	if err != nil {
		return nil, fmt.Errorf("JSON Serializer: Encode %v -> %w", rawMsg, err)
	}

	return rawMsg, nil
}
