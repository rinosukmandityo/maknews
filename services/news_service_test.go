// +build news_service

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
	go test -v -tags=news_service

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	set url=root:Password.1@tcp(127.0.0.1:3306)/tes
	set timeout=10
	set db=tes
	set driver=mysql
*/

var (
	newsService NewsService
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
	repo := rh.ChooseRepo()
	newsService = logic.NewNewsService(repo)
}

func TestNewsService(t *testing.T) {
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
			newsService.Delete(&_data)
			wg.Done()
		}(data)
	}
	wg.Wait()

	t.Run("Case 1: Save data", func(t *testing.T) {
		for _, data := range testdata {
			wg.Add(1)
			go func(_data m.News) {
				if e := newsService.Store(&_data); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		for _, data := range testdata {
			res, e := newsService.GetById(data.ID)
			if e != nil || res.ID == 0 {
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
		if e := newsService.Update(&_data); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.News{ID: -999}
		if e := newsService.Update(&_data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func DeleteData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		if e := newsService.Delete(&_data); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		if e := newsService.Delete(&_data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func GetData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get data", func(t *testing.T) {
		_data := testdata[0]
		if _, e := newsService.GetById(_data.ID); e != nil {
			t.Errorf("[ERROR] - Failed to get data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		if _, e := newsService.GetById(-999); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}
