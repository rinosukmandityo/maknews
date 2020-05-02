package kafka

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"
	"github.com/segmentio/kafka-go"
)

type KafkaRepository struct {
	conn    *kafka.Conn
	url     string
	topic   string
	timeout time.Duration
}

func newKafkaConnection(URL, topic string, timeout int) (*kafka.Conn, error) {
	kafkaConn, e := kafka.DialLeader(context.Background(), "tcp", URL, topic, 0)
	if e != nil {
		return nil, errors.Wrap(e, "repository.newKafkaConnection")
	}
	return kafkaConn, e
}

func NewKafkaConnection(URL, topic string, timeout int) (*KafkaRepository, error) {
	repo := &KafkaRepository{
		topic:   topic,
		url:     URL,
		timeout: time.Duration(timeout) * time.Second,
	}

	conn, e := newKafkaConnection(URL, topic, timeout)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewKafkaConnection")
	}
	repo.conn = conn

	return repo, nil
}

func (k KafkaRepository) Conn() *kafka.Conn {
	return k.conn
}

func (k KafkaRepository) Timeout() time.Duration {
	return k.timeout
}

func (k KafkaRepository) WriteMessage(msgs []byte) error {
	_, e := k.conn.WriteMessages(
		kafka.Message{Value: msgs},
	)
	if e != nil {
		return e
	}
	return nil
}

func (k KafkaRepository) ReadMessage(res chan<- []byte) {
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{k.url},
		Topic:     k.topic,
		Partition: 0,
		MinBytes:  10,
		MaxBytes:  10e3,
	})
	ctx := context.Background()
	lastOffset, _ := k.conn.ReadLastOffset() // get latest offset
	r.SetOffset(lastOffset)                  // set latest offset

	for {
		m, e := r.ReadMessage(ctx)
		if e != nil {
			log.Println("kafka-repo ReadMessage", e.Error())
			break
		}
		// fmt.Printf("message at offset %d: %s = %s at %v\n", m.Offset, string(m.Key), string(m.Value), m.Time)
		res <- m.Value
	}

}
