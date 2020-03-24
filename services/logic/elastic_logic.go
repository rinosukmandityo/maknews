package logic

import (
	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"
	svc "github.com/rinosukmandityo/maknews/services"

	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
)

type elasticService struct {
	repo repo.ElasticRepository
}

func NewElasticService(repo repo.ElasticRepository) svc.ElasticService {
	return &elasticService{
		repo,
	}
}

func (u *elasticService) GetBy(payload m.GetPayload) ([]m.ElasticNews, error) {
	payload.Order = map[string]bool{"created": false}
	res, e := u.repo.GetBy(payload)
	if e != nil {
		return res, errs.Wrap(e, "service.ElasticNews.GetBy")
	}
	return res, nil
}
func (u *elasticService) Store(data m.ElasticNews) error {
	if e := validate.Validate(data); e != nil {
		return errs.Wrap(helper.ErrDataInvalid, "service.ElasticNews.Store")
	}
	if e := u.repo.Store(data); e != nil {
		return errs.Wrap(e, "service.ElasticNews.Store")
	}
	return nil

}
func (u *elasticService) Update(data m.ElasticNews) error {
	if e := u.repo.Update(data, data.ID); e != nil {
		return errs.Wrap(e, "service.ElasticNews.Update")
	}
	return nil

}
func (u *elasticService) Delete(data m.ElasticNews) error {
	if e := u.repo.Delete(data.ID); e != nil {
		return errs.Wrap(e, "service.ElasticNews.Delete")
	}
	return nil

}
