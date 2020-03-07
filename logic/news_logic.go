package logic

import (
	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"
	svc "github.com/rinosukmandityo/maknews/services"

	errs "github.com/pkg/errors"
	"gopkg.in/dealancer/validate.v2"
)

type newsService struct {
	repo repo.NewsRepository
}

func NewNewsService(repo repo.NewsRepository) svc.NewsService {
	return &newsService{
		repo,
	}
}

func (u *newsService) GetById(id int) (*m.News, error) {
	res := new(m.News)
	param := repo.GetParam{
		Tablename: res.TableName(),
		Filter:    map[string]interface{}{"id": id},
		Result:    res,
	}
	if e := u.repo.GetBy(param); e != nil {
		return res, e
	}

	return res, nil

}
func (u *newsService) Store(data *m.News) error {
	if e := validate.Validate(data); e != nil {
		return errs.Wrap(helper.ErrDataInvalid, "service.News.Store")
	}
	if data.ID == 0 {
		data.ID = int(data.Created.UTC().Unix())
	}
	param := repo.StoreParam{
		Tablename: data.TableName(),
		Data: []interface{}{
			data.ID, data.Author, data.Body, data.Created,
		},
	}
	return u.repo.Store(param)

}
func (u *newsService) Update(data *m.News) error {
	if e := validate.Validate(data); e != nil {
		return errs.Wrap(helper.ErrDataInvalid, "service.News.Update")
	}
	if data.ID == 0 {
		data.ID = int(data.Created.UTC().Unix())
	}
	param := repo.UpdateParam{
		Tablename: data.TableName(),
		Filter:    map[string]interface{}{"id": data.ID},
		Data: map[string]interface{}{
			"author":  data.Author,
			"body":    data.Body,
			"created": data.Created,
		},
	}
	return u.repo.Update(param)

}
func (u *newsService) Delete(data *m.News) error {
	if data.ID == 0 {
		return errs.Wrap(helper.ErrDataNotFound, "service.News.Delete")
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
