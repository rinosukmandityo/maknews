// +build repo_test

package repositories_test

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
	go test -v -tags=repo_test

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	=======
	MySQL
	=======
	set url=root:Password.1@tcp(127.0.0.1:3306)/news
	set timeout=10
	set db=news
	set driver=mysql
	=======
	MongoDB
	=======
	set url=mongodb://localhost:27017/local
	set timeout=10
	set db=local
	set driver=mongo
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
	// t.Run("Delete All", DeleteAll)
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
			repo.Delete(_data.ID)
			wg.Done()
		}(data)
	}
	wg.Wait()

	t.Run("Case 1: Save data", func(t *testing.T) {
		for _, data := range testdata {
			wg.Add(1)
			go func(_data m.News) {
				if e := repo.Store(&_data); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		for _, data := range testdata {
			if res, e := repo.GetBy(map[string]interface{}{
				"id": data.ID,
			}); e != nil || res.ID == 0 {
				t.Errorf("[ERROR] - Failed to get data")
			}
		}
	})
}

func UpdateData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Update data", func(t *testing.T) {
		_data := testdata[0]
		data := map[string]interface{}{"author": _data.Author + "UPDATED"}

		if _, e := repo.Update(data, _data.ID); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.News{ID: -9999}
		data := map[string]interface{}{"id": _data.ID}

		if _, e := repo.Update(data, _data.ID); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
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
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func GetData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get data", func(t *testing.T) {
		_data := testdata[0]
		if _, e := repo.GetBy(map[string]interface{}{
			"id": _data.ID,
		}); e != nil {
			t.Errorf("[ERROR] - Failed to get data")
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		if _, e := repo.GetBy(map[string]interface{}{
			"id": -9999,
		}); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func DeleteAll(t *testing.T) {
	wg := sync.WaitGroup{}
	for _, data := range ListTestData() {
		wg.Add(1)
		go func(_data m.News) {
			repo.Delete(_data.ID)
			wg.Done()
		}(data)
	}
	time.Sleep(time.Second * 2)
}
