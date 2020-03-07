package models

import (
	"time"
)

type News struct {
	ID      int       `json:"id" bson:"id" msgpack:"id"`
	Author  string    `json:"author" bson:"author" msgpack:"author"`
	Body    string    `json:"body" bson:"body" msgpack:"body"`
	Created time.Time `json:"created" bson:"created" msgpack:"created"`
}

func (m *News) TableName() string {
	return "news"
}

type ElasticNews struct {
	ID      int       `json:"id" bson:"id" msgpack:"id"`
	Created time.Time `json:"created" bson:"created" msgpack:"created"`
}

func (m *ElasticNews) TableName() string {
	return "news"
}
