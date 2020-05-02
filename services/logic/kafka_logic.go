package logic

import (
	"encoding/json"
	"time"

	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"
	rh "github.com/rinosukmandityo/maknews/repositories/helper"
	svc "github.com/rinosukmandityo/maknews/services"
)

type kafkaService struct {
	repo repo.NewsRepository
}

func NewKafkaService() svc.KafkaService {
	return &kafkaService{}
}

func (u *kafkaService) WriteMessage(data *m.News) error {
	kafka := rh.KafkaConnection()
	conn := kafka.Conn()
	defer conn.Close()
	conn.SetWriteDeadline(time.Now().Add(kafka.Timeout()))
	msgs, _ := json.Marshal(data)
	kafka.WriteMessage(msgs)
	return nil

}
func (u *kafkaService) ReadMessage(newsSvc svc.NewsService, elasticSvc svc.ElasticService) error {
	kafka := rh.KafkaConnection()

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
				if e := elasticSvc.Store(elasticData); e != nil {
					return
				}

				if e := newsSvc.Store(data); e != nil {
					return
				}
			default:
			}
		}
	}()

	kafka.ReadMessage(dataChan)

	return nil
}
