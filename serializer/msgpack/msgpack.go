package msgpack

import (
	"github.com/inblack67/url-shortner/shortener"
	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack/v5"
)

type Redirect struct{}

func (r *Redirect) Encode(input *shortener.Redirect) ([]byte, error) {
	encoded, err := msgpack.Marshal(input)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return encoded, nil
}

func (r *Redirect) Decode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	err := msgpack.Unmarshal(input, redirect)
	if err != nil {
		return nil, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return redirect, nil
}
