package repohelper

import (
	"log"
	"os"
	"strconv"

	repo "github.com/rinosukmandityo/maknews/repositories"
	es "github.com/rinosukmandityo/maknews/repositories/elasticsearch"
	kf "github.com/rinosukmandityo/maknews/repositories/kafka"
	mg "github.com/rinosukmandityo/maknews/repositories/mongodb"
	mr "github.com/rinosukmandityo/maknews/repositories/mysql"
	rr "github.com/rinosukmandityo/maknews/repositories/redis"
)

func ChooseRepo() repo.NewsRepository {
	url := os.Getenv("url")
	db := os.Getenv("db")
	timeout, _ := strconv.Atoi(os.Getenv("timeout"))
	switch os.Getenv("driver") {
	case "mongo":
		if url == "" {
			url = "mongodb://localhost:27017/local"
		}
		if db == "" {
			db = "local"
		}
		if timeout == 0 {
			timeout = 30
		}
		repo, e := mg.NewNewsRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}
		return repo
	default:
		if url == "" {
			url = "root:Password.1@tcp(127.0.0.1:3306)/tes"
		}
		if db == "" {
			db = "tes"
		}
		if timeout == 0 {
			timeout = 10
		}

		repo, e := mr.NewNewsRepository(url, db, timeout)
		if e != nil {
			log.Fatal(e)
		}
		return repo
	}
	return nil
}

func ElasticRepo() repo.ElasticRepository {
	timeout, _ := strconv.Atoi(os.Getenv("elastic_timeout"))
	if timeout == 0 {
		timeout = 10
	}
	url := os.Getenv("elastic_url")
	if url == "" {
		url = "http://localhost:9200"
	}
	index := os.Getenv("elastic_index")
	if index == "" {
		index = "news"
	}
	repo, e := es.NewNewsRepository(url, index, timeout)
	if e != nil {
		log.Fatal(e)
	}
	return repo
}

func KafkaConnection() *kf.KafkaRepository {
	timeout, _ := strconv.Atoi(os.Getenv("kafka_timeout"))
	if timeout == 0 {
		timeout = 10
	}
	url := os.Getenv("kafka_url")
	if url == "" {
		url = "localhost:9092"
	}
	topic := os.Getenv("kafka_topic")
	if topic == "" {
		topic = "news"
	}
	repo, e := kf.NewKafkaConnection(url, topic, timeout)
	if e != nil {
		log.Fatal(e)
	}
	return repo
}

func RedisRepo() repo.CacheRepository {
	timeout, _ := strconv.Atoi(os.Getenv("redis_expired"))
	if timeout == 0 {
		timeout = 10
	}
	url := os.Getenv("redis_url")
	if url == "" {
		url = "redis://:@localhost:6379/0"
	}
	repo, e := rr.NewNewsRepository(url, timeout)
	if e != nil {
		log.Fatal(e)
	}
	return repo
}
