package repositories

import (
	m "github.com/rinosukmandityo/maknews/models"
)

type KafkaRepository interface {
	WriteMessage(data *m.News) error
	ReadMessage(res chan<- []byte)
}
