// +build repo_test

package repositories_test

import (
	"fmt"
	"strings"
	"sync"
	"testing"
	"time"

	_ "github.com/rinosukmandityo/maknews/api"
	"github.com/rinosukmandityo/maknews/helper"
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

type TestTable struct {
	name        string
	expected    string
	expectedErr error
	errMsg      string
	updatedData map[string]interface{}
	filter      map[string]interface{}
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
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to save data",
			data:        ListTestData(),
		},
	}

	for _, tt := range tts {
		testdata := tt.data

		// Clean test data if any
		for _, data := range testdata {
			repo.Delete(data.ID)
		}

		t.Run(tt.name, func(t *testing.T) {
			for _, data := range testdata {
				if e := repo.Store(&data); e != tt.expectedErr {
					t.Errorf("%s %s ", tt.errMsg, e.Error())
				}
			}

			for _, data := range testdata {
				if res, e := repo.GetBy(map[string]interface{}{
					"id": data.ID,
				}); e != nil || res.ID == 0 {
					t.Errorf("[ERROR] - Failed to get data")
				}
			}
		})
	}
}

func UpdateData(t *testing.T) {
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to update data",
			data:        []m.News{ListTestData()[0]},
			updatedData: map[string]interface{}{"author": ListTestData()[0].Author + "UPDATED"},
		},
		{
			name:        "Case: Negative Test",
			expected:    "",
			expectedErr: helper.ErrDataNotFound,
			errMsg:      fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			data:        []m.News{{ID: -9999}},
			updatedData: map[string]interface{}{"author": "Data Not Exists"},
		},
	}

	for _, tt := range tts {
		testdata := tt.data
		t.Run(tt.name, func(t *testing.T) {
			for _, data := range testdata {
				if _, e := repo.Update(tt.updatedData, data.ID); e != tt.expectedErr {
					if !strings.Contains(e.Error(), tt.expectedErr.Error()) {
						t.Errorf("%s %s ", tt.errMsg, e.Error())
					}
				}
			}
		})
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
		t.Run(tt.name, func(t *testing.T) {
			for _, data := range testdata {
				if e := repo.Delete(data.ID); e != tt.expectedErr {
					if !strings.Contains(e.Error(), tt.expectedErr.Error()) {
						t.Errorf("%s %s ", tt.errMsg, e.Error())
					}
				}
			}
		})
	}
}

func GetData(t *testing.T) {
	tts := []TestTable{
		{
			name:        "Case: Positive Test",
			expected:    "",
			expectedErr: nil,
			errMsg:      "[ERROR] - Failed to get data",
			filter:      map[string]interface{}{"id": ListTestData()[0].ID},
		},
		{
			name:        "Case: Negative Test",
			expected:    "",
			expectedErr: helper.ErrDataNotFound,
			errMsg:      fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			filter:      map[string]interface{}{"id": -9999},
		},
	}

	for _, tt := range tts {
		t.Run(tt.name, func(t *testing.T) {
			if _, e := repo.GetBy(tt.filter); e != tt.expectedErr {
				if !strings.Contains(e.Error(), tt.expectedErr.Error()) {
					t.Errorf("%s %s ", tt.errMsg, e.Error())
				}
			}
		})
	}
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
