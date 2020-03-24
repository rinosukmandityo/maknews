package elastic

import (
	m "github.com/rinosukmandityo/maknews/models"

	elasticapi "github.com/olivere/elastic/v7"
)

func constructDeleteQuery(filter map[string]interface{}) *elasticapi.BoolQuery {
	q := elasticapi.NewBoolQuery()
	queries := []elasticapi.Query{}
	for k, v := range filter {
		queries = append(queries, elasticapi.NewTermQuery(k, v))
	}
	q = q.Must(queries...)

	return q
}

func constructGetBy(payload m.GetPayload) *elasticapi.BoolQuery {
	q := elasticapi.NewBoolQuery()
	queries := []elasticapi.Query{}
	for k, v := range payload.Filter {
		queries = append(queries, elasticapi.NewTermQuery(k, v))
	}
	q = q.Must(queries...)

	return q
}
