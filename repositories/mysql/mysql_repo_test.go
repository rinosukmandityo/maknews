// +build mysql_test

package mysql_test

import (
	"sync"
	"testing"
	"time"

	_ "github.com/rinosukmandityo/maknews/api"
	m "github.com/rinosukmandityo/maknews/models"
	. "github.com/rinosukmandityo/maknews/repositories"
	rh "github.com/rinosukmandityo/maknews/repositories/helper"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=mysql_test

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	set url=root:Password.1@tcp(127.0.0.1:3306)/tes
	set timeout=10
	set db=tes
	set driver=mysql
*/

var (
	repo NewsRepository
)

func ListTestData() []m.News {
	return []m.News{{
		ID:      1,
		Author:  "Alex",
		Body:    "Hello this is news from Alex",
		Created: time.Now().UTC(),
	}, {
		ID:      2,
		Author:  "Bacca",
		Body:    "Hello this is news from Bacca",
		Created: time.Now().UTC().Add(time.Second * 3),
	}, {
		ID:      3,
		Author:  "Chicarito",
		Body:    "Hello this is news from Chicarito",
		Created: time.Now().UTC().Add(time.Second * 5),
	}}
}

func init() {
	repo = rh.ChooseRepo()
}

func TestService(t *testing.T) {
	t.Run("Insert Data", InsertData)
	t.Run("Update Data", UpdateData)
	t.Run("Delete Data", DeleteData)
	t.Run("Get Data", GetData)
}

func InsertData(t *testing.T) {
	testdata := ListTestData()
	wg := sync.WaitGroup{}

	// Clean test data if any
	for _, data := range testdata {
		wg.Add(1)
		go func(_data m.News) {
			param := DeleteParam{
				Tablename: _data.TableName(),
				Filter: map[string]interface{}{
					"id": _data.ID,
				},
			}
			repo.Delete(param)
			wg.Done()
		}(data)
	}
	wg.Wait()

	t.Run("Case 1: Save data", func(t *testing.T) {
		for _, data := range testdata {
			wg.Add(1)
			go func(_data m.News) {
				param := StoreParam{
					Tablename: _data.TableName(),
					Data: []interface{}{
						_data.ID, _data.Author, _data.Body, _data.Created,
					},
				}
				if e := repo.Store(param); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		for _, data := range testdata {
			res := new(m.News)
			param := GetParam{
				Tablename: data.TableName(),
				Filter: map[string]interface{}{
					"id": data.ID,
				},
				Result: res,
			}
			if e := repo.GetBy(param); e != nil || res.ID == 0 {
				t.Errorf("[ERROR] - Failed to get data")
			}
		}
	})
}

func UpdateData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Update data", func(t *testing.T) {
		_data := testdata[0]
		_data.Author += "UPDATED"
		param := UpdateParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
			Data: map[string]interface{}{
				"id":      _data.ID,
				"author":  _data.Author,
				"body":    _data.Body,
				"created": _data.Created,
			},
		}

		if e := repo.Update(param); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.News{ID: -9999}
		param := UpdateParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
			Data: map[string]interface{}{
				"author":  _data.Author,
				"body":    _data.Body,
				"created": _data.Created,
			},
		}
		if e := repo.Update(param); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func DeleteData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		param := DeleteParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
		}
		if e := repo.Delete(param); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		param := DeleteParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
		}
		if e := repo.Delete(param); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func GetData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get data", func(t *testing.T) {
		_data := testdata[0]
		res := new(m.News)
		param := GetParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
			Result: res,
		}
		if e := repo.GetBy(param); e != nil || res.ID == 0 {
			t.Errorf("[ERROR] - Failed to get data")
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		res := new(m.News)
		param := GetParam{
			Tablename: res.TableName(),
			Filter: map[string]interface{}{
				"id": -9999,
			},
			Result: res,
		}
		if e := repo.GetBy(param); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}
