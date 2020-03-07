package elastic

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"

	elasticapi "github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
)

type newsElasticRepository struct {
	client  *elasticapi.Client
	index   string
	timeout time.Duration
}

func newNewsClient(URL string) (*elasticapi.Client, error) {
	client, e := elasticapi.NewClient(elasticapi.SetURL(URL))
	if e != nil {
		return nil, e
	}
	if e := ping(context.Background(), client, URL); e != nil {
		return client, e
	}
	return client, e
}

func NewNewsRepository(URL, index string, timeout int) (repo.NewsRepository, error) {
	repo := &newsElasticRepository{
		index:   index,
		timeout: time.Duration(timeout) * time.Second,
	}

	client, e := newNewsClient(URL)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewNewsRepository")
	}
	repo.client = client

	return repo, nil
}

func ping(ctx context.Context, client *elasticapi.Client, url string) error {

	// Ping the Elasticsearch server to get HttpStatus, version number
	if client != nil {
		info, code, e := client.Ping(url).Do(ctx)
		if e != nil {
			return e
		}

		fmt.Printf("Elasticsearch returned with code %d and version %s \n", code, info.Version.Number)
		return nil
	}
	return errors.New("elastic client is nil")
}

func getResult(param repo.GetParam, searchResult *elasticapi.SearchResult) (e error) {
	if searchResult.TotalHits() == 0 {
		return errors.Wrap(errors.New("Data Not Found"), "repository.News.GetBy")
	}
	res := []models.ElasticNews{}
	for _, hit := range searchResult.Hits.Hits {
		_res := models.ElasticNews{}
		if e := json.Unmarshal(hit.Source, &_res); e != nil {
			return errors.Wrap(e, "repository.News.Update")
		}
		res = append(res, _res)
	}
	resByte, _ := json.Marshal(res)
	json.Unmarshal(resByte, &param.Result)

	return nil
}

func CreateIndexIfDoesNotExist(ctx context.Context, client *elasticapi.Client, indexName string) error {
	exists, e := client.IndexExists(indexName).Do(ctx)
	if e != nil {
		return e
	}

	if exists {
		return nil
	}

	res, e := client.CreateIndex(indexName).Do(ctx)

	if e != nil {
		return e
	}

	if !res.Acknowledged {
		return errors.New("CreateIndex was not acknowledged. Check that timeout value is correct.")
	}

	return nil
}

func (r *newsElasticRepository) GetBy(param repo.GetParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	q := constructGetBy(param)
	searchService := r.client.Search().
		Index(r.index).
		Query(q).
		From(param.Offset)

	if param.Limit > 0 {
		searchService.Size(param.Limit)
	}

	if len(param.Order) > 0 {
		for k, v := range param.Order {
			searchService.Sort(k, v)
		}
	}

	searchResult, e := searchService.Do(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.News.GetBy")
	}

	if e := getResult(param, searchResult); e != nil {
		return e
	}

	return nil
}

func (r *newsElasticRepository) Store(param repo.StoreParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	data := param.Data.(models.ElasticNews)
	_, e := r.client.Index().Index(r.index).Type(param.Tablename).Id(strconv.Itoa(data.ID)).BodyJson(data).Do(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}

	return nil

}

func (r *newsElasticRepository) Update(param repo.UpdateParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	data := param.Data.(map[string]interface{})
	idInt := strconv.Itoa(param.Filter["id"].(int))

	res, e := r.client.Update().Index(r.index).
		Type(param.Tablename).Id(idInt).Doc(data).Do(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.News.Update")
	}

	if res.Result != "updated" {
		return errors.Wrap(errors.New("Data Not Found"), "repository.News.Update")
	}

	return nil

}
func (r *newsElasticRepository) Delete(param repo.DeleteParam) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()

	q := constructDeleteQuery(param)
	res, e := r.client.DeleteByQuery(r.index).Type(param.Tablename).Query(q).Do(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	}
	if res.Total == 0 {
		return errors.Wrap(errors.New("Data Not Found"), "repository.News.Delete")
	}

	// Flush data (need for refreshing data in index) after this command possible to do get.
	r.client.Flush().Index(r.index).Do(ctx)

	return nil

}
