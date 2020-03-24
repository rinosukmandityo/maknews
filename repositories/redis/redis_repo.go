package redis

import (
	"encoding/json"
	"time"

	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"

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

func NewNewsRepository(redisURL string, expiration int) (repo.CacheRepository, error) {
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

func (r *newsRedisRepository) GetBy(param m.GetPayload) ([]m.News, error) {
	res := []m.News{}
	keyList, e := r.client.ZRevRange(helper.REDIS_KEY_SET, int64(param.Offset), int64(param.Offset+param.Limit)).Result()
	if e != nil {
		return res, errors.Wrap(e, "repository.News.GetBy")
	}
	for _, key := range keyList {
		dataRedis, e := r.client.Get(key).Result()
		if e != nil {
			return res, errors.Wrap(e, "repository.News.GetBy")
		}
		dataByte := []byte(dataRedis)
		if len(dataByte) == 0 {
			return res, errors.Wrap(helper.ErrDataNotFound, "repository.News.GetBy")
		}
		_res := m.News{}
		if e := json.Unmarshal(dataByte, &_res); e != nil {
			return res, errors.Wrap(e, "repository.News.GetBy")
		}
		res = append(res, _res)
	}
	return res, nil
}

func (r *newsRedisRepository) Store(data []m.News) error {
	for _, v := range data {
		key, score := generateKeyScore(v)
		if _, e := r.client.Set(key, v, r.expiration).Result(); e != nil {
			return errors.Wrap(e, "repository.News.Update")
		}
		r.client.Expire(key, r.expiration)

		member := &redis.Z{
			Score:  score,
			Member: key,
		}
		r.client.ZAdd(helper.REDIS_KEY_SET, member)
	}
	return nil

}

func (r *newsRedisRepository) Update(data m.News) error {
	key, _ := generateKeyScore(data)
	if _, e := r.client.Set(key, data, r.expiration).Result(); e != nil {
		return errors.Wrap(e, "repository.News.Update")
	}
	r.client.Expire(key, r.expiration)
	return nil

}
func (r *newsRedisRepository) Delete(data m.News) error {
	key, score := generateKeyScore(data)
	if _, e := r.client.Del(key).Result(); e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	}
	member := &redis.Z{
		Score:  score,
		Member: key,
	}
	r.client.ZRem(helper.REDIS_KEY_SET, member)

	return nil

}
