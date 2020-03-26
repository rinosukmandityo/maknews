package models

import (
	"time"
)

type News struct {
	ID      int       `json:"id" bson:"_id" msgpack:"_id" db:"id"`
	Author  string    `json:"author" bson:"author" msgpack:"author" db:"author"`
	Body    string    `json:"body" bson:"body" msgpack:"body" db:"body"`
	Created time.Time `json:"created" bson:"created" msgpack:"created" db:"created"`
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
