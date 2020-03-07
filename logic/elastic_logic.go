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
	repo repo.NewsRepository
}

func NewElasticService(repo repo.NewsRepository) svc.ElasticService {
	return &elasticService{
		repo,
	}
}

func (u *elasticService) GetBy(payload m.GetPayload) ([]m.ElasticNews, error) {
	res := []m.ElasticNews{}
	param := repo.GetParam{
		Tablename: new(m.ElasticNews).TableName(),
		Filter:    payload.Filter,
		Result:    &res,
		Offset:    payload.Offset,
		Limit:     payload.Limit,
		Order:     map[string]bool{"created": false},
	}
	if e := u.repo.GetBy(param); e != nil {
		return res, e
	}

	return res, nil

}
func (u *elasticService) Store(data m.ElasticNews) error {
	if e := validate.Validate(data); e != nil {
		return errs.Wrap(helper.ErrDataInvalid, "service.ElasticNews.Store")
	}
	if data.ID == 0 {
		data.ID = int(data.Created.UTC().Unix())
	}
	param := repo.StoreParam{
		Tablename: data.TableName(),
		Data:      data,
	}
	return u.repo.Store(param)

}
func (u *elasticService) Update(data m.ElasticNews) error {
	if e := validate.Validate(data); e != nil {
		return errs.Wrap(helper.ErrDataInvalid, "service.ElasticNews.Update")
	}
	if data.ID == 0 {
		data.ID = int(data.Created.UTC().Unix())
	}
	param := repo.UpdateParam{
		Tablename: data.TableName(),
		Filter: map[string]interface{}{
			"id": data.ID,
		},
		Data: map[string]interface{}{
			"id":      data.ID,
			"created": data.Created,
		},
	}
	return u.repo.Update(param)

}
func (u *elasticService) Delete(data m.ElasticNews) error {
	if data.ID == 0 {
		return errs.Wrap(helper.ErrDataNotFound, "service.ElasticNews.Delete")
	}
	param := repo.DeleteParam{
		Tablename: data.TableName(),
		Filter:    map[string]interface{}{"id": data.ID},
	}
	if e := u.repo.Delete(param); e != nil {
		return e
	}
	return nil

}
