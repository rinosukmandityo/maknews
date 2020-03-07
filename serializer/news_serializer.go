package serializer

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type UserSerializer interface {
	Decode(input []byte) (*m.News, error)
	Encode(input *m.News) ([]byte, error)
	DecodeGetPayload(input []byte) (m.GetPayload, error)
	EncodeGetPayload(input *m.GetPayload) ([]byte, error)
	EncodeGetData(input []m.News) ([]byte, error)
}
