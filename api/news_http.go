package api

import (
	"io/ioutil"
	"net/http"

	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	svc "github.com/rinosukmandityo/maknews/services"

	"github.com/pkg/errors"
)

type NewsHandler interface {
	Get(http.ResponseWriter, *http.Request)
	Post(http.ResponseWriter, *http.Request)
	Update(http.ResponseWriter, *http.Request)
	Delete(http.ResponseWriter, *http.Request)
}

type newshandler struct {
	newsService    svc.NewsService
	elasticService svc.ElasticService
	kafkaService   svc.KafkaService
	redisService   svc.RedisService
}

func NewNewsHandler(newsService svc.NewsService, elasticService svc.ElasticService, kafkaService svc.KafkaService,
	redisService svc.RedisService) NewsHandler {
	return &newshandler{newsService, elasticService, kafkaService, redisService}
}

func worker(jobs <-chan int, res chan<- *m.News, err chan<- error, newsService svc.NewsService) {
	for id := range jobs {
		result, e := newsService.GetById(id)
		if e != nil {
			err <- e
			return
		}
		res <- result
	}
}

func (u *newshandler) getDataWithWorker(elasticData []m.ElasticNews) []m.News {
	idList := []int{}
	for _, v := range elasticData {
		idList = append(idList, v.ID)
	}
	jobs := make(chan int, len(idList))
	res := make(chan *m.News, len(idList))
	err := make(chan error)
	resMap := map[int]m.News{}
	for _, id := range idList {
		jobs <- id
	}
	close(jobs)
	for i := 0; i < 3; i++ {
		go worker(jobs, res, err, u.newsService)
	}

	for i := 0; i < len(idList); i++ {
		result := <-res
		resMap[result.ID] = *result
	}
	data := []m.News{}
	for _, id := range idList {
		data = append(data, resMap[id])
	}
	return data
}

func (u *newshandler) Get(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	payload, e := GetSerializer(contentType).DecodeGetPayload(requestBody)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	data := []m.News{}
	data, e = u.redisService.GetData(payload)
	if e != nil {
		elasticData, e := u.elasticService.GetBy(payload)
		if e != nil {
			if errors.Cause(e) == helper.ErrDataNotFound {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		data = u.getDataWithWorker(elasticData)
	}
	if e := u.redisService.StoreData(data, payload); e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respBody, e := GetSerializer(contentType).EncodeGetData(data)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusFound)
}

func (u *newshandler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, e := GetSerializer(contentType).Decode(requestBody)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if e := u.kafkaService.WriteMessage(data); e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	respBody, e := GetSerializer(contentType).Encode(data)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusCreated)
}

func (u *newshandler) Update(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, e := GetSerializer(contentType).Decode(requestBody)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if e = u.newsService.Update(data); e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respBody, e := GetSerializer(contentType).Encode(data)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusOK)

}

func (u *newshandler) Delete(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	data, e := GetSerializer(contentType).Decode(requestBody)
	if e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	if e = u.newsService.Delete(data); e != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	respBody, e := GetSerializer(contentType).Encode(data)
	SetupResponse(w, contentType, respBody, http.StatusOK)
}
