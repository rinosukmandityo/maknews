package repositories

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type NewsRepository interface {
	GetBy(filter map[string]interface{}) (*m.News, error)
	Store(data *m.News) error
	Update(data map[string]interface{}, id int) (*m.News, error)
	Delete(id int) error
}
