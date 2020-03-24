package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"

	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type newsMySQLRepository struct {
	url     string
	timeout time.Duration
}

func newNewsClient(URL string) (*sql.DB, error) {
	db, e := sql.Open("mysql", URL)
	if e != nil {
		return nil, e
	}
	if e = db.Ping(); e != nil {
		return nil, e
	}
	return db, e
}

func NewNewsRepository(URL, DB string, timeout int) (repo.NewsRepository, error) {
	repo := &newsMySQLRepository{
		url:     fmt.Sprintf("%s?parseTime=true", URL),
		timeout: time.Duration(timeout) * time.Second,
	}
	return repo, nil
}

func (r *newsMySQLRepository) GetBy(filter map[string]interface{}) (*m.News, error) {
	res := new(m.News)
	db, e := sqlx.Connect("mysql", r.url)
	if e != nil {
		return res, errors.Wrap(e, "repository.News.GetBy")
	}
	defer db.Close()
	q := constructGetBy(filter)

	if e = db.Get(res, q); e != nil {
		return res, errors.Wrap(e, "repository.News.GetBy")
	}
	return res, nil

}
func (r *newsMySQLRepository) Store(data *m.News) error {
	db, e := newNewsClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	defer conn.Close()

	q, dataField := constructStoreQuery(data)
	tx, e := db.Begin()
	if e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	if _, e = tx.ExecContext(ctx, q, dataField...); e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	tx.Commit()

	return nil

}

func (r *newsMySQLRepository) Update(data map[string]interface{}, id int) (*m.News, error) {
	news := new(m.News)
	db, e := newNewsClient(r.url)
	if e != nil {
		return news, errors.Wrap(e, "repository.News.Store")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return news, errors.Wrap(e, "repository.News.Update")
	}
	defer conn.Close()

	filter := map[string]interface{}{"id": id}
	q, dataField := constructUpdateQuery(data, filter)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return news, errors.Wrap(e, "repository.News.Update")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(dataField...); e != nil {
		return news, errors.Wrap(e, "repository.News.Update")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return news, errors.Wrap(e, "repository.News.Update")
		}
		if count == 0 {
			return news, errors.Wrap(helper.ErrDataNotFound, "repository.News.Update")
		}
	}
	news, e = r.GetBy(filter)
	if e != nil {
		return news, errors.Wrap(e, "repository.News.Update")
	}

	return news, nil

}
func (r *newsMySQLRepository) Delete(id int) error {
	db, e := newNewsClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	}
	defer conn.Close()

	filter := map[string]interface{}{"id": id}
	q, data := constructDeleteQuery(filter)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(data...); e != nil {
		return errors.Wrap(e, "repository.News.Delete")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.News.Delete")
		}
		if count == 0 {
			return errors.Wrap(helper.ErrDataNotFound, "repository.News.Delete")
		}
	}

	return nil

}
