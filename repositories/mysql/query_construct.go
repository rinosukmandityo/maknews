package mysql

import (
	"fmt"
	"strings"

	m "github.com/rinosukmandityo/maknews/models"
)

func constructUpdateQuery(data, filter map[string]interface{}) (string, []interface{}) {
	// 	"UPDATE <tablename> SET field1=?, field2=?  WHERE filter1=?"
	q := fmt.Sprintf("UPDATE %s SET", new(m.News).TableName())
	values := []interface{}{}
	for k, v := range data {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")
	q += " WHERE"
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, values
}

func constructDeleteQuery(filter map[string]interface{}) (string, []interface{}) {
	// 	"DELETE <tablename> WHERE filter1=?"
	q := fmt.Sprintf("DELETE FROM %s WHERE", new(m.News).TableName())
	values := []interface{}{}
	for k, v := range filter {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, values
}

func constructStoreQuery(data *m.News) (string, []interface{}) {
	// 	"INSERT INTO <tablename> VALUES(?, ?, ?, ?)"
	_data := []interface{}{
		data.ID, data.Author, data.Body, data.Created,
	}
	q := fmt.Sprintf("INSERT INTO %s VALUES(", data.TableName())
	values := []interface{}{}
	for _, v := range _data {
		q += "?,"
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",") + ")"

	return q, values
}

func constructGetBy(filter map[string]interface{}) string {
	// SELET * FROM <tablename> WHERE filter1=filtervalue
	q := fmt.Sprintf("SELECT * FROM %s WHERE", new(m.News).TableName())
	for k, v := range filter {
		q += fmt.Sprintf(" %s=%v,", k, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q
}
