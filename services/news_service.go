package services

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type NewsService interface {
	GetData(payload m.GetPayload) ([]m.News, error)
	GetById(id int) (*m.News, error)
	Store(data *m.News) error
	Update(data map[string]interface{}, id int) (*m.News, error)
	Delete(data m.News) error
}
