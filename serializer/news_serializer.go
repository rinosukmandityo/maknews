package serializer

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type UserSerializer interface {
	Decode(input []byte) (*m.News, error)
	Encode(input *m.News) ([]byte, error)
	DecodeMap(input []byte) (map[string]interface{}, error)
	EncodeMap(input map[string]interface{}) ([]byte, error)
	EncodeGetData(input []m.News) ([]byte, error)
}
