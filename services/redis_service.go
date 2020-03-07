package services

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type RedisService interface {
	StoreData(data []m.News, payload m.GetPayload) error
	GetData(payload m.GetPayload) ([]m.News, error)
}
