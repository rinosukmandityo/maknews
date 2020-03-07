package services

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type KafkaService interface {
	WriteMessage(data *m.News) error
	ReadMessage(newsSvc NewsService, elasticSvc ElasticService) error
}
