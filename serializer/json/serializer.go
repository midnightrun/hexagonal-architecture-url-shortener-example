package json

import (
	"github.com/midnightrun/hexagonal-architecture-url-shortener-example/shortener"
)

type Redirect struct{}

func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	return nil, nil
}
