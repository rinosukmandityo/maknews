package repositories

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type ElasticRepository interface {
	GetBy(param m.GetPayload) ([]m.ElasticNews, error)
	Store(data m.ElasticNews) error
	Update(data m.ElasticNews, id int) error
	Delete(id int) error
}
