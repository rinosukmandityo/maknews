package mysql

import (
	"context"
	"fmt"
	"time"

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

func (r *newsMySQLRepository) GetBy(param repo.GetParam) error {
	db, e := sqlx.Connect("mysql", r.url)
	if e != nil {
		return errors.Wrap(e, "repository.News.GetBy")
	}
	defer db.Close()
	q := constructGetBy(param)

	if e = db.Get(param.Result, q); e != nil {
		return errors.Wrap(e, "repository.News.GetBy")
	}
	return nil

}
func (r *newsMySQLRepository) Store(param repo.StoreParam) error {
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

	q, data := constructStoreQuery(param)
	tx, e := db.Begin()
	if e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	if _, e = tx.ExecContext(ctx, q, data...); e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	tx.Commit()

	return nil

}

func (r *newsMySQLRepository) Update(param repo.UpdateParam) error {
	db, e := newNewsClient(r.url)
	if e != nil {
		return errors.Wrap(e, "repository.News.Store")
	}
	defer db.Close()
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	conn, e := db.Conn(ctx)
	if e != nil {
		return errors.Wrap(e, "repository.News.Update")
	}
	defer conn.Close()

	q, data := constructUpdateQuery(param)
	stmt, e := conn.PrepareContext(ctx, q)
	if e != nil {
		return errors.Wrap(e, "repository.News.Update")
	}
	defer stmt.Close()
	if res, e := stmt.Exec(data...); e != nil {
		return errors.Wrap(e, "repository.News.Update")
	} else {
		count, e := res.RowsAffected()
		if e != nil {
			return errors.Wrap(e, "repository.News.Update")
		}
		if count == 0 {
			return errors.Wrap(errors.New("Data Not Found"), "repository.News.Update")
		}
	}

	return nil

}
func (r *newsMySQLRepository) Delete(param repo.DeleteParam) error {
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

	q, data := constructDeleteQuery(param)
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
			return errors.Wrap(errors.New("Data Not Found"), "repository.News.Delete")
		}
	}

	return nil

}
