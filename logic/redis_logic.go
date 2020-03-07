package logic

import (
	"encoding/json"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"
	svc "github.com/rinosukmandityo/maknews/services"
)

type redisService struct {
	repo repo.NewsRepository
}

func NewRedisService(repo repo.NewsRepository) svc.RedisService {
	return &redisService{
		repo,
	}
}

func (u *redisService) GetData(payload m.GetPayload) ([]m.News, error) {
	res := []m.News{}
	param := repo.GetParam{
		Tablename: new(m.ElasticNews).TableName(),
		Filter:    payload.Filter,
		Result:    &res,
		Offset:    payload.Offset,
		Limit:     payload.Limit,
	}
	if e := u.repo.GetBy(param); e != nil {
		return res, e
	}

	return res, nil

}
func (u *redisService) StoreData(data []m.News, payload m.GetPayload) error {
	dataBytes, e := json.Marshal(data)
	if e != nil {
		return e
	}
	payloadData := map[string]interface{}{
		"data":   dataBytes,
		"offset": payload.Offset,
		"limit":  payload.Limit,
	}

	param := repo.StoreParam{
		Data: payloadData,
	}
	return u.repo.Store(param)

}
