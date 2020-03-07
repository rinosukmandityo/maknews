package mysql

import (
	"fmt"
	"strings"

	repo "github.com/rinosukmandityo/maknews/repositories"
)

func constructUpdateQuery(param repo.UpdateParam) (string, []interface{}) {
	// 	"UPDATE <tablename> SET field1=?, field2=?  WHERE filter1=?"
	data := param.Data.(map[string]interface{})
	q := fmt.Sprintf("UPDATE %s SET", param.Tablename)
	values := []interface{}{}
	for k, v := range data {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")
	q += " WHERE"
	for k, v := range param.Filter {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, values
}

func constructDeleteQuery(param repo.DeleteParam) (string, []interface{}) {
	// 	"DELETE <tablename> WHERE filter1=?"
	q := fmt.Sprintf("DELETE FROM %s WHERE", param.Tablename)
	values := []interface{}{}
	for k, v := range param.Filter {
		q += fmt.Sprintf(" %s=?,", k)
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q, values
}

func constructStoreQuery(param repo.StoreParam) (string, []interface{}) {
	// 	"INSERT INTO <tablename> VALUES(?, ?, ?, ?)"
	data := param.Data.([]interface{})
	q := fmt.Sprintf("INSERT INTO %s VALUES(", param.Tablename)
	values := []interface{}{}
	for _, v := range data {
		q += "?,"
		values = append(values, v)
	}
	q = strings.TrimSuffix(q, ",") + ")"

	return q, values
}

func constructGetBy(param repo.GetParam) string {
	// SELET * FROM <tablename> WHERE filter1=filtervalue
	q := fmt.Sprintf("SELECT * FROM %s WHERE", param.Tablename)
	for k, v := range param.Filter {
		q += fmt.Sprintf(" %s=%v,", k, v)
	}
	q = strings.TrimSuffix(q, ",")

	return q
}
