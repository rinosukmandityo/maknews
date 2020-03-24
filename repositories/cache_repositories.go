package repositories

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type CacheRepository interface {
	GetBy(param m.GetPayload) ([]m.News, error)
	Store(data []m.News) error
	Update(data m.News) error
	Delete(data m.News) error
}
