package mongo

import (
	"context"
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"reflect"
	"time"

	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type newsMongoRepository struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

func newNewsClient(mongoURL string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, e := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if e != nil {
		return nil, e
	}
	if e = client.Ping(ctx, readpref.Primary()); e != nil {
		return nil, e
	}
	return client, e
}

func NewNewsRepository(mongoURL, mongoDB string, mongoTimeout int) (repo.NewsRepository, error) {
	repo := &newsMongoRepository{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDB,
	}
	client, e := newNewsClient(mongoURL, mongoTimeout)
	if e != nil {
		return nil, errors.Wrap(e, "repository.NewNewsRepository")
	}
	repo.client = client
	return repo, nil
}

func (r *newsMongoRepository) GetBy(filter map[string]interface{}) (*m.News, error) {
	res := new(m.News)
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(res.TableName())
	convertID(filter)
	if e := c.FindOne(ctx, filter).Decode(res); e != nil {
		if e == mongo.ErrNoDocuments {
			return res, errors.Wrap(helper.ErrDataNotFound, "repository.News.GetById")
		}
		return res, errors.Wrap(e, "repository.News.GetById")
	}
	return res, nil

}
func (r *newsMongoRepository) Store(data *m.News) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	c := r.client.Database(r.database).Collection(data.TableName())
	if _, e := c.InsertOne(ctx, data); e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}

	return nil

}
func (r *newsMongoRepository) Update(data map[string]interface{}, id int) (*m.News, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	news := new(m.News)
	c := r.client.Database(r.database).Collection(news.TableName())
	filter := map[string]interface{}{"_id": id}
	convertID(data)
	if fmt.Sprintf("%T", data["created"]) == "string" {
		created, _ := time.Parse("2006-01-02T15:04:05Z", data["created"].(string))
		data["created"] = created
	}
	if res, e := c.UpdateOne(ctx, filter, bson.M{"$set": data}, options.Update().SetUpsert(false)); e != nil {
		return news, errors.Wrap(e, "repository.News.Update")
	} else {
		if res.MatchedCount == 0 && res.ModifiedCount == 0 {
			return news, errors.Wrap(helper.ErrDataNotFound, "repository.News.Update")
		}
	}
	news, e := r.GetBy(filter)
	if e != nil {
		return news, errors.Wrap(e, "repository.News.Update")
	}

	return news, nil

}
func (r *newsMongoRepository) Delete(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	filter := map[string]interface{}{"_id": id}
	c := r.client.Database(r.database).Collection(new(m.News).TableName())
	if res, e := c.DeleteOne(ctx, filter); e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	} else {
		if res.DeletedCount == 0 {
			return errors.Wrap(helper.ErrDataNotFound, "repository.News.Delete")
		}
	}

	return nil
}

func convertID(data map[string]interface{}) {
	if _, ok := data["id"]; ok {
		switch helper.Kind(data["id"]) {
		case reflect.Float64:
			data["_id"] = int(data["id"].(float64))
		case reflect.Int:
			data["_id"] = data["id"].(int)
		}
		delete(data, "id")
	}
}
