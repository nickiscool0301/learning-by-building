package codec

import (
	"bytes"
	"encoding/gob"
)

type GobCodec struct{}

func (g *GobCodec) Encode(v any) ([]byte, error) {
	// Implementation for encoding using gob
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (g *GobCodec) Decode(data []byte, v any) error {
	// Implementation for decoding using gob
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	return decoder.Decode(v)
}

func NewGobCodec() *GobCodec {
	return &GobCodec{}
}
