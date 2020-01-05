package msgpack

import (
	"fmt"

	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
	"github.com/vmihailenco/msgpack"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := msgpack.Unmarshal(input, redirect); err != nil {
		return nil, fmt.Errorf("MsgPack Serializer: Decode %v -> %w", redirect, err)
	}

	return redirect, nil
}

func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := msgpack.Marshal(input)

	if err != nil {
		return nil, fmt.Errorf("MsgPack Serializer: Encode %v -> %w", rawMsg, err)
	}

	return rawMsg, nil
}
