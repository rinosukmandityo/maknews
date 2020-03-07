package elastic

import (
	repo "github.com/rinosukmandityo/maknews/repositories"

	elasticapi "github.com/olivere/elastic/v7"
)

func constructDeleteQuery(param repo.DeleteParam) *elasticapi.BoolQuery {
	q := elasticapi.NewBoolQuery()
	queries := []elasticapi.Query{}
	for k, v := range param.Filter {
		queries = append(queries, elasticapi.NewTermQuery(k, v))
	}
	q = q.Must(queries...)

	return q
}

func constructGetBy(param repo.GetParam) *elasticapi.BoolQuery {
	q := elasticapi.NewBoolQuery()
	queries := []elasticapi.Query{}
	for k, v := range param.Filter {
		queries = append(queries, elasticapi.NewTermQuery(k, v))
	}
	q = q.Must(queries...)

	return q
}
