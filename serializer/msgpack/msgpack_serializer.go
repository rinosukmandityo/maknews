package msgpack

import (
	m "github.com/rinosukmandityo/maknews/models"

	"github.com/pkg/errors"
	"github.com/vmihailenco/msgpack"
)

type News struct{}

func (u *News) Decode(input []byte) (*m.News, error) {
	user := new(m.News)
	if e := msgpack.Unmarshal(input, user); e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.Decode")
	}
	return user, nil
}

func (u *News) Encode(input *m.News) ([]byte, error) {
	rawMsg, e := msgpack.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.Encode")
	}
	return rawMsg, nil
}

func (u *News) DecodeMap(input []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	if e := msgpack.Unmarshal(input, &res); e != nil {
		return res, errors.Wrap(e, "serializer.Logic.DecodeMap")
	}
	return res, nil
}

func (u *News) EncodeMap(input map[string]interface{}) ([]byte, error) {
	rawMsg, e := msgpack.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.EncodeMap")
	}
	return rawMsg, nil
}

func (u *News) EncodeGetData(input []m.News) ([]byte, error) {
	rawMsg, e := msgpack.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.EncodeGetData")
	}
	return rawMsg, nil
}
