package services

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type RedisService interface {
	StoreData(data []m.News) error
	GetData(payload m.GetPayload) ([]m.News, error)
	UpdateData(data m.News) error
	DeleteData(data m.News) error
}
