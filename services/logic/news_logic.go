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
	repo        repo.NewsRepository
	redisRepo   repo.CacheRepository
	elasticRepo repo.ElasticRepository
	kafkaRepo   repo.KafkaRepository
}

func NewNewsService(repo repo.NewsRepository, redisRepo repo.CacheRepository, elasticRepo repo.ElasticRepository,
	kafkaRepo repo.KafkaRepository) svc.NewsService {
	return &newsService{
		repo,
		redisRepo,
		elasticRepo,
		kafkaRepo,
	}
}

func IsDataEmpty(data m.News) bool {
	return data.ID == 0 && data.Author == "" && data.Body == "" && data.Created.IsZero()
}

func (u *newsService) worker(jobs <-chan int, res chan<- *m.News, err chan<- error) {
	for id := range jobs {
		result, e := u.GetById(id)
		if e != nil {
			err <- e
		} else {
			err <- nil
		}
		res <- result
	}
}

func (u *newsService) getDataWithWorker(elasticData []m.ElasticNews) []m.News {
	lenData := len(elasticData)
	jobs := make(chan int, lenData)
	res := make(chan *m.News, lenData)
	err := make(chan error, lenData)
	resMap := map[int]m.News{}
	for i := 0; i < 3; i++ {
		go u.worker(jobs, res, err)
	}

	for _, v := range elasticData {
		jobs <- v.ID
	}
	close(jobs)

	for i := 0; i < lenData; i++ {
		result := <-res
		resMap[result.ID] = *result
	}

	errs := []error{}
	for i := 0; i < lenData; i++ {
		e := <-err
		errs = append(errs, e)
	}

	data := []m.News{}
	for _, v := range elasticData {
		if !IsDataEmpty(resMap[v.ID]) {
			data = append(data, resMap[v.ID])
		}
	}
	return data
}

func (u *newsService) GetData(payload m.GetPayload) ([]m.News, error) {
	data, e := u.redisRepo.GetBy(payload)

	if e != nil || len(data) == 0 {
		payload.Order = map[string]bool{"created": false}
		elasticData, e := u.elasticRepo.GetBy(payload)
		if e != nil {
			if errs.Cause(e) == helper.ErrDataNotFound {
				e = helper.ErrDataNotFound
			}
			return data, e
		}
		data = u.getDataWithWorker(elasticData)
	}
	if e := u.redisRepo.Store(data); e != nil {
		return data, e
	}

	return data, nil
}

func (u *newsService) GetById(id int) (*m.News, error) {
	filter := map[string]interface{}{"id": id}
	res, e := u.repo.GetBy(filter)
	if e != nil {
		return res, errs.Wrap(e, "service.News.GetById")
	}

	return res, nil

}
func (u *newsService) Store(data *m.News) error {
	if e := u.kafkaRepo.WriteMessage(data); e != nil {
		return errs.Wrap(e, "service.News.Store")
	}

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
	if e := u.redisRepo.Update(*updatedData); e != nil {
		return updatedData, e
	}
	return updatedData, nil

}
func (u *newsService) Delete(existingData m.News) error {
	if e := u.repo.Delete(existingData.ID); e != nil {
		return e
	}
	if e := u.redisRepo.Delete(existingData); e != nil {
		return e
	}
	return nil

}
