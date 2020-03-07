package redis

import (
	"encoding/json"
	"github.com/rinosukmandityo/maknews/helper"
	repo "github.com/rinosukmandityo/maknews/repositories"
	"time"

	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

type newsRedisRepository struct {
	client     *redis.Client
	expiration time.Duration
}

func newNewsClient(redisURL string) (*redis.Client, error) {
	// opt, err := redis.ParseURL("redis://:qwerty@localhost:6379/1")
	opt, e := redis.ParseURL(redisURL)
	if e != nil {
		return nil, e
	}
	client := redis.NewClient(opt)
	if _, e = client.Ping().Result(); e != nil {
		return nil, e
	}
	return client, e
}

func NewNewsRepository(redisURL string, expiration int) (repo.NewsRepository, error) {
	repo := &newsRedisRepository{
		expiration: time.Duration(expiration) * time.Second,
	}
	client, e := newNewsClient(redisURL)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewNewsRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *newsRedisRepository) GetBy(param repo.GetParam) error {
	key := generateKeyOffsetLimit(param.Offset, param.Limit)
	dataRedis, e := r.client.HGetAll(key).Result()
	if e != nil {
		return errors.Wrap(e, "repository.News.GetById")
	}
	data := []byte(dataRedis["data"])
	if len(data) == 0 {
		return errors.Wrap(helper.ErrDataNotFound, "repository.News.GetById")
	}
	if e := json.Unmarshal(data, &param.Result); e != nil {
		return errors.Wrap(e, "repository.News.GetBy")
	}
	return nil
}
func (r *newsRedisRepository) Store(param repo.StoreParam) error {
	data := param.Data.(map[string]interface{})
	key := generateKeyOffsetLimit(data["offset"].(int), data["limit"].(int))
	dataBytes := data["data"].([]byte)
	dataItem := map[string]interface{}{"data": dataBytes}
	if _, e := r.client.HMSet(key, dataItem).Result(); e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	r.client.Expire(key, r.expiration)
	return nil

}
func (r *newsRedisRepository) Update(param repo.UpdateParam) error {
	data := param.Data.(map[string]interface{})
	key := generateKey(data["id"].(string))
	if _, e := r.client.HMSet(key, data).Result(); e != nil {
		return errors.Wrap(e, "repository.News.Update")
	}
	return nil

}
func (r *newsRedisRepository) Delete(param repo.DeleteParam) error {
	key := generateKey(param.Filter["id"].(string))
	if _, e := r.client.HDel(key).Result(); e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	}

	return nil

}
