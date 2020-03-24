package logic

import (
	"github.com/pkg/errors"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"
	svc "github.com/rinosukmandityo/maknews/services"
)

type redisService struct {
	repo repo.CacheRepository
}

func NewRedisService(repo repo.CacheRepository) svc.RedisService {
	return &redisService{
		repo,
	}
}

func (u *redisService) GetData(payload m.GetPayload) ([]m.News, error) {
	res, e := u.repo.GetBy(payload)
	if e != nil {
		return res, errors.Wrap(e, "service.News.GetData")
	}

	return res, nil

}

func (u *redisService) StoreData(data []m.News) error {
	if e := u.repo.Store(data); e != nil {
		return errors.Wrap(e, "service.News.Store")
	}
	return nil

}

func (u *redisService) UpdateData(data m.News) error {
	if e := u.repo.Update(data); e != nil {
		return errors.Wrap(e, "service.News.Update")
	}
	return nil
}

func (u *redisService) DeleteData(data m.News) error {
	if e := u.repo.Delete(data); e != nil {
		return errors.Wrap(e, "service.News.Update")
	}
	return nil
}
