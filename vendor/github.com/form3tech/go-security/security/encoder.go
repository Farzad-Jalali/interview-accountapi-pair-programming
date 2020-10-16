package security

import (
	"bytes"
	"compress/zlib"
	"encoding/base64"
	"fmt"
	"io"
)

type Encoder interface {
	Encode(AccessControlList, io.Writer) error
}

type compressEncoder struct {
	_       struct{}
	encoder Encoder
}

func (c *compressEncoder) Encode(acls AccessControlList, w io.Writer) error {
	buffer := &bytes.Buffer{}

	if err := c.encoder.Encode(acls, buffer); err != nil {
		return err
	}

	zlibWriter := zlib.NewWriter(w)
	if _, err := zlibWriter.Write(buffer.Bytes()); err != nil {
		return err
	}
	if err := zlibWriter.Flush(); err != nil {
		return err
	}
	return zlibWriter.Close()
}

func newCompressEncoder(e Encoder) Encoder {
	return &compressEncoder{
		encoder: e,
	}
}

var _ Encoder = &compressEncoder{}

type base64Encoder struct {
	_       struct{}
	encoder Encoder
}

func (e base64Encoder) Encode(acls AccessControlList, writer io.Writer) error {
	buffer := &bytes.Buffer{}
	if err := e.encoder.Encode(acls, buffer); err != nil {
		return err
	}
	base64Writer := base64.NewEncoder(base64.StdEncoding.WithPadding('='), writer)
	if _, err := base64Writer.Write(buffer.Bytes()); err != nil {
		return fmt.Errorf("could not encode access control list, error: %v", err)
	}
	return base64Writer.Close()
}

func newBase64Encoder(encoder Encoder) Encoder {
	return &base64Encoder{
		encoder: encoder,
	}
}

var _ Encoder = &base64Encoder{}
