// +build elastic_service

package services_test

import (
	"sync"
	"testing"
	"time"

	_ "github.com/rinosukmandityo/maknews/api"
	"github.com/rinosukmandityo/maknews/logic"
	m "github.com/rinosukmandityo/maknews/models"
	rh "github.com/rinosukmandityo/maknews/repositories/helper"
	. "github.com/rinosukmandityo/maknews/services"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=elastic_service
*/

var (
	elasticService ElasticService
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
	repo := rh.ElasticRepo()
	elasticService = logic.NewElasticService(repo)
}

func TestElasticService(t *testing.T) {
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
		go func(_data m.ElasticNews) {
			elasticService.Delete(_data)
			wg.Done()
		}(data)
	}
	wg.Wait()

	t.Run("Case 1: Save data", func(t *testing.T) {
		for _, data := range testdata {
			wg.Add(1)
			go func(_data m.ElasticNews) {
				if e := elasticService.Store(_data); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		time.Sleep(time.Second * 2)

		for _, data := range testdata {
			res := []m.ElasticNews{}
			payload := m.GetPayload{
				Filter: map[string]interface{}{"id": data.ID},
			}
			res, e := elasticService.GetBy(payload)
			if e != nil || len(res) == 0 {
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
		if e := elasticService.Update(_data); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.ElasticNews{ID: -999}
		if e := elasticService.Update(_data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func DeleteData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		if e := elasticService.Delete(_data); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		if e := elasticService.Delete(_data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func GetData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get data", func(t *testing.T) {
		_data := testdata[0]
		res := []m.ElasticNews{}
		payload := m.GetPayload{
			Filter: map[string]interface{}{"id": _data.ID},
		}
		res, e := elasticService.GetBy(payload)
		if e != nil || len(res) == 0 {
			t.Errorf("[ERROR] - Failed to get data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		payload := m.GetPayload{
			Filter: map[string]interface{}{"id": -999},
		}
		if _, e := elasticService.GetBy(payload); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}
