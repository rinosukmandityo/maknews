package api

import (
	"context"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	svc "github.com/rinosukmandityo/maknews/services"

	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

type NewsHandler interface {
	NewsCtx(http.Handler) http.Handler
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

func (u *newshandler) NewsCtx(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		idInt, e := strconv.Atoi(id)
		if e != nil {
			http.Error(w, helper.ErrDataInvalid.Error(), http.StatusBadRequest)
		}
		data, e := u.newsService.GetById(idInt)
		if e != nil {
			if errors.Cause(e) == helper.ErrDataNotFound {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, e.Error(), http.StatusBadRequest)
		}
		ctx := context.WithValue(r.Context(), "news", data)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (u *newshandler) Get(w http.ResponseWriter, r *http.Request) {
	payload := m.GetPayload{
		Offset: 0,
		Limit:  10,
	}
	q := r.URL.Query()
	if q.Get("offset") != "" {
		payload.Offset, _ = strconv.Atoi(q.Get("offset"))
	}
	if q.Get("limit") != "" {
		payload.Limit, _ = strconv.Atoi(q.Get("limit"))
	}

	contentType := r.Header.Get("Content-Type")

	data, e := u.redisService.GetData(payload)
	if e != nil {
		elasticData, e := u.elasticService.GetBy(payload)
		if e != nil {
			if errors.Cause(e) == helper.ErrDataNotFound {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			http.Error(w, e.Error(), http.StatusBadRequest)
			return
		}
		data = u.getDataWithWorker(elasticData)
	}
	if e := u.redisService.StoreData(data); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	respBody, e := GetSerializer(contentType).EncodeGetData(data)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusFound)
}

func (u *newshandler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	data, e := GetSerializer(contentType).Decode(requestBody)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	if e := u.kafkaService.WriteMessage(data); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}

	// if e := u.redisService.StoreData([]m.News{data}); e != nil {
	// 	http.Error(w, e.Error(), http.StatusBadRequest)
	// 	return
	// }

	respBody, e := GetSerializer(contentType).Encode(data)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusCreated)
}

func (u *newshandler) Update(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	existingData, ok := ctx.Value("news").(*m.News)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := existingData.ID
	contentType := r.Header.Get("Content-Type")
	requestBody, e := ioutil.ReadAll(r.Body)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	data, e := GetSerializer(contentType).DecodeMap(requestBody)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	updatedData, e := u.newsService.Update(data, id)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	respBody, e := GetSerializer(contentType).Encode(updatedData)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := u.redisService.UpdateData(*updatedData); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusOK)

}

func (u *newshandler) Delete(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	existingData, ok := ctx.Value("news").(*m.News)
	if !ok {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}
	id := existingData.ID
	contentType := r.Header.Get("Content-Type")
	if e := u.newsService.Delete(id); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	respBody, e := GetSerializer(contentType).EncodeMap(map[string]interface{}{"ID": id})
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	if e := u.redisService.DeleteData(*existingData); e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	SetupResponse(w, contentType, respBody, http.StatusOK)
}
