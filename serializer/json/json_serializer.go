package json

import (
	m "github.com/rinosukmandityo/maknews/models"

	"encoding/json"
	"github.com/pkg/errors"
)

type News struct{}

func (u *News) Decode(input []byte) (*m.News, error) {
	news := new(m.News)
	if e := json.Unmarshal(input, news); e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.Decode")
	}
	return news, nil
}

func (u *News) Encode(input *m.News) ([]byte, error) {
	rawMsg, e := json.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.Encode")
	}
	return rawMsg, nil
}

func (u *News) DecodeMap(input []byte) (map[string]interface{}, error) {
	res := map[string]interface{}{}
	if e := json.Unmarshal(input, &res); e != nil {
		return res, errors.Wrap(e, "serializer.Logic.DecodeMap")
	}
	return res, nil
}

func (u *News) EncodeMap(input map[string]interface{}) ([]byte, error) {
	rawMsg, e := json.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.EncodeMap")
	}
	return rawMsg, nil
}

func (u *News) EncodeGetData(input []m.News) ([]byte, error) {
	rawMsg, e := json.Marshal(input)
	if e != nil {
		return nil, errors.Wrap(e, "serializer.Logic.EncodeGetData")
	}
	return rawMsg, nil
}
