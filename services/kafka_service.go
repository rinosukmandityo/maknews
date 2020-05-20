package services

import (
	repo "github.com/rinosukmandityo/maknews/repositories"
)

type KafkaService interface {
	ReadMessage(newsRepo repo.NewsRepository, elasticRepo repo.ElasticRepository) error
}
