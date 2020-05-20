package logic

import (
	"encoding/json"

	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"
	svc "github.com/rinosukmandityo/maknews/services"
)

type kafkaService struct {
	repo repo.KafkaRepository
}

func NewKafkaService(repo repo.KafkaRepository) svc.KafkaService {
	return &kafkaService{
		repo,
	}
}

func (u *kafkaService) ReadMessage(newsRepo repo.NewsRepository, elasticRepo repo.ElasticRepository) error {
	dataChan := make(chan []byte) // it will be sent to ReadMessage function

	go func() {
		for {
			select {
			case dataByte := <-dataChan:
				data := new(m.News)
				if e := json.Unmarshal(dataByte, data); e != nil {
					return
				}
				elasticData := m.ElasticNews{
					ID:      data.ID,
					Created: data.Created,
				}
				if e := elasticRepo.Store(elasticData); e != nil {
					return
				}

				if e := newsRepo.Store(data); e != nil {
					return
				}
			default:
			}
		}
	}()

	u.repo.ReadMessage(dataChan)

	return nil
}
