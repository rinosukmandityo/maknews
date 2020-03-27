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
	repo ElasticRepository
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
			repo.Delete(_data.ID)
			wg.Done()
		}(data)
	}
	wg.Wait()

	t.Run("Case 1: Save data", func(t *testing.T) {
		for _, data := range testdata {
			wg.Add(1)
			go func(_data m.ElasticNews) {
				if e := repo.Store(_data); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		time.Sleep(time.Second * 2)
		for _, data := range testdata {
			param := m.GetPayload{
				Filter: map[string]interface{}{
					"id": data.ID,
				},
				Offset: 0,
				Limit:  10,
			}
			if res, e := repo.GetBy(param); e != nil || len(res) == 0 {
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

		if e := repo.Update(_data, _data.ID); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.ElasticNews{ID: -9999}
		if e := repo.Update(_data, _data.ID); e == nil {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func DeleteData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		if e := repo.Delete(_data.ID); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		if e := repo.Delete(_data.ID); e == nil {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func GetData(t *testing.T) {
	testdata := ListTestData()
	_data := testdata[0]
	t.Run("Case 1: Get data", func(t *testing.T) {
		param := m.GetPayload{
			Filter: map[string]interface{}{
				"id": _data.ID,
			},
			Order: map[string]bool{
				"created": true,
			},
			Offset: 0,
			Limit:  10,
		}
		if res, e := repo.GetBy(param); e != nil || len(res) == 0 {
			t.Errorf("[ERROR] - Failed to get data")
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		param := m.GetPayload{
			Filter: map[string]interface{}{
				"id": -9999,
			},
			Order: map[string]bool{
				"created": true,
			},
			Offset: 0,
			Limit:  10,
		}
		if _, e := repo.GetBy(param); e == nil {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func GetAll(t *testing.T) {
	t.Run("Case 1: Get all data", func(t *testing.T) {
		res := []m.ElasticNews{}
		param := m.GetPayload{
			Order: map[string]bool{
				"created": false,
			},
			Offset: 0,
			Limit:  10,
		}
		res, e := repo.GetBy(param)
		if e != nil || len(res) == 0 {
			t.Errorf("[ERROR] - Failed to get data")
		}
		var prevDate time.Time
		for _, v := range res {
			if !prevDate.IsZero() && v.Created.After(prevDate) {
				t.Errorf("[ERROR] - Incorrect order data")
			}
			prevDate = v.Created
		}
	})
}

func DeleteAll(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		for _, _data := range testdata {
			if e := repo.Delete(_data.ID); e != nil {
				t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
			}
		}
	})
}
