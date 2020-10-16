package security

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"io/ioutil"

	"github.com/pkg/errors"
)

type Decoder interface {
	Decode(payload []byte) (AccessControl, error)
}

type decompress struct {
	_       struct{}
	decoder Decoder
}

func (d decompress) Decode(payload []byte) (AccessControl, error) {
	r, err := zlib.NewReader(bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return d.decoder.Decode(data)
}

func newDecompress(decoder Decoder) Decoder {
	return &decompress{
		decoder: decoder,
	}
}

var _ Decoder = &decompress{}

type base64Decoder struct {
	_       struct{}
	decoder Decoder
}

func (d base64Decoder) Decode(payload []byte) (AccessControl, error) {
	data, err := base64.StdEncoding.DecodeString(string(payload))
	if err != nil {
		return nil, errors.Wrap(err, "unable to decode acls")
	}
	return d.decoder.Decode(data)
}

func newBase64Decoder(decoder Decoder) Decoder {
	return &base64Decoder{
		decoder: decoder,
	}
}

var _ Decoder = &base64Decoder{}
