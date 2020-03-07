package services

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type ElasticService interface {
	GetBy(payload m.GetPayload) ([]m.ElasticNews, error)
	Store(data m.ElasticNews) error
	Update(data m.ElasticNews) error
	Delete(data m.ElasticNews) error
}
