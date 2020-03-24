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
	filter := map[string]interface{}{"ID": id}
	res, e := u.repo.GetBy(filter)
	if e != nil {
		return res, errs.Wrap(e, "service.News.GetById")
	}

	return res, nil

}
func (u *newsService) Store(data *m.News) error {
	if e := validate.Validate(data); e != nil {
		return errs.Wrap(helper.ErrDataInvalid, "service.News.Store")
	}
	return u.repo.Store(data)

}
func (u *newsService) Update(data map[string]interface{}, id int) (*m.News, error) {
	updatedData, e := u.repo.Update(data, id)
	if e != nil {
		return updatedData, errs.Wrap(e, "service.News.Update")
	}
	return updatedData, nil

}
func (u *newsService) Delete(id int) error {
	if e := u.repo.Delete(id); e != nil {
		return e
	}
	return nil

}
