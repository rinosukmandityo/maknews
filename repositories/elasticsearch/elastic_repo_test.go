// +build elastic_test

package elastic_test

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
	go test -v -tags=elastic_test
*/

var (
	repo NewsRepository
)

func ListTestData() []m.ElasticNews {
	return []m.ElasticNews{{
		ID:      1,
		Created: time.Now().UTC(),
	}, {
		ID:      2,
		Created: time.Now().UTC().Add(time.Second * 3),
	}, {
		ID:      3,
		Created: time.Now().UTC().Add(time.Second * 5),
	}}
}

func init() {
	repo = rh.ElasticRepo()
}

func TestService(t *testing.T) {
	t.Run("Insert Data", InsertData)
	t.Run("Get All", GetAll)
	t.Run("Update Data", UpdateData)
	t.Run("Delete Data", DeleteData)
	t.Run("Get Data", GetData)
	// t.Run("Delete All", DeleteAll)
}

func InsertData(t *testing.T) {
	testdata := ListTestData()
	wg := sync.WaitGroup{}

	// Clean test data if any
	for _, data := range testdata {
		wg.Add(1)
		go func(_data m.ElasticNews) {
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
			go func(_data m.ElasticNews) {
				param := StoreParam{
					Tablename: _data.TableName(),
					Data:      _data,
				}
				if e := repo.Store(param); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		time.Sleep(time.Second * 2)
		for _, data := range testdata {
			res := []m.ElasticNews{}
			param := GetParam{
				Tablename: data.TableName(),
				Filter: map[string]interface{}{
					"id": data.ID,
				},
				Result: &res,
			}
			if e := repo.GetBy(param); e != nil || len(res) == 0 {
				t.Errorf("[ERROR] - Failed to get data")
			}
		}
	})
}

func UpdateData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Update data", func(t *testing.T) {
		_data := testdata[0]
		_data.Created = time.Time{}
		param := UpdateParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
			Data: map[string]interface{}{
				"id":      _data.ID,
				"created": _data.Created,
			},
		}

		if e := repo.Update(param); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.ElasticNews{ID: -9999}
		param := UpdateParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
			Data: map[string]interface{}{
				"id":      _data.ID,
				"created": _data.Created,
			},
		}
		if e := repo.Update(param); e == nil {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
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
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func GetData(t *testing.T) {
	testdata := ListTestData()
	_data := testdata[0]
	t.Run("Case 1: Get data", func(t *testing.T) {
		res := []m.ElasticNews{}
		param := GetParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
			Result: &res,
			Order: map[string]bool{
				"created": true,
			},
			Offset: 0,
			Limit:  10,
		}
		if e := repo.GetBy(param); e != nil || len(res) == 0 {
			t.Errorf("[ERROR] - Failed to get data")
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		res := []m.ElasticNews{}
		param := GetParam{
			Tablename: _data.TableName(),
			Filter: map[string]interface{}{
				"id": -9999,
			},
			Result: &res,
			Order: map[string]bool{
				"created": true,
			},
			Offset: 0,
			Limit:  10,
		}
		if e := repo.GetBy(param); e == nil {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func GetAll(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get all data", func(t *testing.T) {
		res := []m.ElasticNews{}
		param := GetParam{
			Tablename: new(m.ElasticNews).TableName(),
			Result:    &res,
			Order: map[string]bool{
				"created": false,
			},
			Offset: 0,
			Limit:  10,
		}
		if e := repo.GetBy(param); e != nil || len(res) == 0 {
			t.Errorf("[ERROR] - Failed to get data")
		}
		idList := []int{}
		for i := len(testdata) - 1; i >= 0; i-- {
			idList = append(idList, testdata[i].ID)
		}
		for i, v := range res {
			if idList[i] != v.ID {
				t.Errorf("[ERROR] - Incorrect order data")
			}
		}
	})
}

func DeleteAll(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		for _, _data := range testdata {
			param := DeleteParam{
				Tablename: _data.TableName(),
				Filter: map[string]interface{}{
					"id": _data.ID,
				},
			}
			if e := repo.Delete(param); e != nil {
				t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
			}
		}
	})
}
