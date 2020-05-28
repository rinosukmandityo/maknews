// +build news_service

package services_test

import (
	"errors"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	rh "github.com/rinosukmandityo/maknews/repositories/helper"
	. "github.com/rinosukmandityo/maknews/services"
	"github.com/rinosukmandityo/maknews/services/logic"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=news_service

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
	newsService NewsService
)

type TestTable struct {
	name        string
	expected    string
	expectedErr error
	errMsg      string
	updatedData []map[string]interface{}
	filter      []map[string]interface{}
	data        []m.News
}

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
	cacheRepo := rh.RedisRepo()
	elasticRepo := rh.ElasticRepo()
	kafkaRepo := rh.KafkaConnection()
	newsService = logic.NewNewsService(repo, cacheRepo, elasticRepo, kafkaRepo)
}

func TestNewsService(t *testing.T) {
	t.Run("Insert Data", InsertData)
	t.Run("Update Data", UpdateData)
	t.Run("Delete Data", DeleteData)
	t.Run("Get Data By ID", GetDataById)
	t.Run("Get Data", GetData)
}

func InsertData(t *testing.T) {
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to save data",
			data:        ListTestData(),
		},
	}

	// Clean test data if any
	for _, _data := range ListTestData() {
		newsService.Delete(_data)
	}
	time.Sleep(time.Second * 1)

	for _, tt := range tts {
		for _, _data := range tt.data {
			t.Run(tt.name, func(t *testing.T) {
				if e := newsService.Store(&_data); e != tt.expectedErr {
					t.Errorf("%s %s ", tt.errMsg, e.Error())
				}
				res, e := newsService.GetById(_data.ID)
				if e != nil || res.ID == 0 {
					t.Errorf("[ERROR] - Failed to get data")
				}
			})
		}
	}
	time.Sleep(time.Second * 1)
}

func UpdateData(t *testing.T) {
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to update data",
			data:        []m.News{ListTestData()[0]},
			updatedData: []map[string]interface{}{
				{"author": ListTestData()[0].Author + "UPDATED"},
			},
		},
		{
			name:        "Case: Negative Test",
			expected:    "",
			expectedErr: helper.ErrDataNotFound,
			errMsg:      fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			data:        []m.News{{ID: -9999}},
			updatedData: []map[string]interface{}{
				{"author": "Data Not Exists"},
			},
		},
	}

	for _, tt := range tts {
		testdata := tt.data
		for i, _data := range testdata {
			t.Run(tt.name, func(t *testing.T) {
				if _, e := newsService.Update(tt.updatedData[i], _data.ID); e != tt.expectedErr {
					if !strings.Contains(e.Error(), tt.expectedErr.Error()) {
						t.Errorf("%s %s ", tt.errMsg, e.Error())
					}
				}
			})
		}
	}
}

func DeleteData(t *testing.T) {
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to delete data",
			data:        []m.News{ListTestData()[1]},
		},
		{
			name:        "Case: Negative Test",
			expected:    "",
			expectedErr: helper.ErrDataNotFound,
			errMsg:      fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			data:        []m.News{ListTestData()[1]},
		},
	}

	for _, tt := range tts {
		testdata := tt.data
		for _, data := range testdata {
			t.Run(tt.name, func(t *testing.T) {
				if e := newsService.Delete(data); e != tt.expectedErr {
					if !strings.Contains(e.Error(), tt.expectedErr.Error()) {
						t.Errorf("%s %s", tt.errMsg, e.Error())
					}
				}
			})
		}
	}
}

func GetDataById(t *testing.T) {
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to get data",
			filter:      []map[string]interface{}{{"id": ListTestData()[0].ID}},
		},
		{
			name:        "Case: Negative Test",
			expected:    "",
			expectedErr: helper.ErrDataNotFound,
			errMsg:      fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			filter:      []map[string]interface{}{{"id": -9999}},
		},
	}

	for _, tt := range tts {
		for _, filter := range tt.filter {
			t.Run(tt.name, func(t *testing.T) {
				if _, e := newsService.GetById(filter["id"].(int)); e != tt.expectedErr {
					if !strings.Contains(e.Error(), tt.expectedErr.Error()) {
						t.Errorf("%s %s ", tt.errMsg, e.Error())
					}
				}
			})
		}
	}
}

func GetData(t *testing.T) {
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to get data",
			filter:      []map[string]interface{}{{"offset": 0, "limit": 10}},
		},
		{
			name:        "Case: Negative Test",
			expected:    "",
			expectedErr: errors.New("Offset can not be less than zero"),
			errMsg:      fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			filter:      []map[string]interface{}{{"offset": -1, "limit": -10}},
		},
	}

	for _, tt := range tts {
		for _, filter := range tt.filter {
			t.Run(tt.name, func(t *testing.T) {
				payload := m.GetPayload{
					Offset: filter["offset"].(int),
					Limit:  filter["limit"].(int),
				}
				if _, e := newsService.GetData(payload); e != tt.expectedErr {
					if !strings.Contains(e.Error(), tt.expectedErr.Error()) {
						t.Errorf("%s %s ", tt.errMsg, e.Error())
					}
				}
			})
		}
	}
}
