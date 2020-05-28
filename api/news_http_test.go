// +build news_http

package api_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	. "github.com/rinosukmandityo/maknews/api"
	"github.com/rinosukmandityo/maknews/helper"
	m "github.com/rinosukmandityo/maknews/models"
	repo "github.com/rinosukmandityo/maknews/repositories"
	rh "github.com/rinosukmandityo/maknews/repositories/helper"
	"github.com/rinosukmandityo/maknews/services/logic"

	"github.com/go-chi/chi"
)

/*
	==================
	RUN FROM TERMINAL
	==================
	go test -v -tags=news_http

	===================================
	TO SET DATABASE INFO FROM TERMINAL
	===================================
	set url=root:Password.1@tcp(127.0.0.1:3306)/news
	set timeout=10
	set db=news
	set driver=mysql
	set redis_url=redis://:@localhost:6379/0
	set redis_timeout=10
	set elastic_url=http://localhost:9200
	set elastic_timeout=10
	set elastic_index=news
	set kafka_url=localhost:9092
	set kafka_timeout=10
	set kafka_topic=news
*/

var (
	newsRepo    repo.NewsRepository
	elasticRepo repo.ElasticRepository
	cacheRepo   repo.CacheRepository
	kafkaRepo   repo.KafkaRepository
	r           *chi.Mux
	ts          *httptest.Server
)

type TestTable struct {
	name               string
	expectedStatusCode int
	expectedErr        error
	errMsg             string
	updatedData        []map[string]interface{}
	filter             []map[string]interface{}
	data               []m.News
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
	newsRepo = rh.ChooseRepo()
	elasticRepo = rh.ElasticRepo()
	cacheRepo = rh.RedisRepo()
	kafkaRepo = rh.KafkaConnection()
	r = RegisterHandler()
}

func TestNewsHTTP(t *testing.T) {
	ts = httptest.NewServer(r)
	defer ts.Close()

	t.Run("Insert Data", InsertData)
	t.Run("Get All Data", GetAllData)
	t.Run("Update Data", UpdateData)
	t.Run("Delete Data", DeleteData)
	// t.Run("Get Data", GetDataByID) // not implemented yet
}

func PostReq(t *testing.T, ts *httptest.Server, url string, _data m.News) (*http.Response, string, error) {
	dataBytes, e := getBytes(_data)
	if e != nil {
		return nil, "", e
	}
	resp, respBody, e := makeRequest(t, ts, "POST", url, bytes.NewReader(dataBytes))
	if e != nil {
		return resp, respBody, e
	}

	return resp, respBody, nil
}

func PutReq(t *testing.T, ts *httptest.Server, url string, _data map[string]interface{}) (*http.Response, string, error) {
	dataBytes, e := json.Marshal(_data)
	if e != nil {
		return nil, "", e
	}
	resp, respBody, e := makeRequest(t, ts, "PUT", url, bytes.NewReader(dataBytes))
	if e != nil {
		return resp, respBody, e
	}

	return resp, respBody, nil
}

func DeleteReq(t *testing.T, ts *httptest.Server, url string) (*http.Response, string, error) {
	resp, respBody, e := makeRequest(t, ts, "DELETE", url, nil)
	if e != nil {
		return resp, respBody, e
	}

	return resp, respBody, nil
}

func GetReq(t *testing.T, ts *httptest.Server, url string) (*http.Response, string, error) {
	resp, respBody, e := makeRequest(t, ts, "GET", url, nil)
	if e != nil {
		return resp, respBody, e
	}

	return resp, respBody, nil
}

func getBytes(_data m.News) ([]byte, error) {
	dataBytes, e := GetSerializer(ContentTypeJson).Encode(&_data)
	if e != nil {
		return dataBytes, e
	}
	return dataBytes, nil
}

func InsertData(t *testing.T) {
	newsService := logic.NewNewsService(newsRepo, cacheRepo, elasticRepo, kafkaRepo)

	tts := []TestTable{
		{
			name:               "Case: Positive Test",
			expectedStatusCode: http.StatusCreated,
			expectedErr:        nil,
			errMsg:             "[ERROR] - Status should be 'Status Created' (201)",
			data:               ListTestData(),
		},
	}

	testdata := ListTestData()
	// Clean test data if any
	for _, _data := range testdata {
		elData := m.ElasticNews{
			ID:      _data.ID,
			Created: _data.Created,
		}
		newsService.Delete(_data)
		elasticRepo.Delete(elData.ID)
		cacheRepo.Delete(_data)
	}
	time.Sleep(time.Second * 1)

	for _, tt := range tts {
		for _, _data := range tt.data {
			t.Run(tt.name, func(t *testing.T) {
				if resp, _, e := PostReq(t, ts, "/news", _data); e != tt.expectedErr && resp.StatusCode != tt.expectedStatusCode {
					t.Errorf("%s %s ", tt.errMsg, e.Error())
				}

				res, e := newsService.GetById(_data.ID)
				if e != nil || res.ID == 0 {
					t.Errorf("[ERROR] - Failed to get data")
				}

				payload := m.GetPayload{
					Offset: 0,
					Limit:  10,
				}

				if elRes, e := elasticRepo.GetBy(payload); e != nil || len(elRes) == 0 {
					t.Errorf("[ERROR] - Failed to get data from elastic search")
				}
			})
		}
	}

	time.Sleep(time.Second * 1)
}

func UpdateData(t *testing.T) {
	tts := []TestTable{
		{
			name:               "Case: Positive Test",
			expectedStatusCode: http.StatusOK,
			expectedErr:        nil,
			errMsg:             "[ERROR] - Status should be 'Status OK' (200)",
			data:               []m.News{ListTestData()[0]},
			updatedData:        []map[string]interface{}{{"author": ListTestData()[0].Author + "UPDATED"}},
		},
		{
			name:               "Case: Negative Test",
			expectedStatusCode: http.StatusNotFound,
			expectedErr:        nil,
			errMsg:             fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			data:               []m.News{{ID: -9999}},
			updatedData:        []map[string]interface{}{{"author": "Data Not Exists"}},
		},
	}

	for _, tt := range tts {
		for i, _data := range tt.data {
			t.Run(tt.name, func(t *testing.T) {
				if resp, _, e := PutReq(t, ts, fmt.Sprintf("/news/%d", _data.ID), tt.updatedData[i]); e != tt.expectedErr && resp.StatusCode != tt.expectedStatusCode {
					t.Errorf("%s %s ", tt.errMsg, e.Error())
				}
			})
		}
	}
	time.Sleep(time.Second * 1)
}

func DeleteData(t *testing.T) {
	tts := []TestTable{
		{
			name:               "Case: Positive Test",
			expectedStatusCode: http.StatusOK,
			expectedErr:        nil,
			errMsg:             "[ERROR] - Failed to delete data",
			data:               []m.News{ListTestData()[1]},
		},
		{
			name:               "Case: Negative Test",
			expectedStatusCode: http.StatusNotFound,
			expectedErr:        nil,
			errMsg:             fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			data:               []m.News{ListTestData()[1]},
		},
	}

	for _, tt := range tts {
		for _, _data := range tt.data {
			t.Run(tt.name, func(t *testing.T) {
				if resp, _, e := DeleteReq(t, ts, fmt.Sprintf("/news/%d", _data.ID)); e != tt.expectedErr && resp.StatusCode != tt.expectedStatusCode {
					t.Errorf("%s %s ", tt.errMsg, e.Error())
				}
			})
		}
	}
	time.Sleep(time.Second * 1)
}

func GetDataByID(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get Data", func(t *testing.T) {
		_data := testdata[0]
		if _, _, e := GetReq(t, ts, fmt.Sprintf("/news/?id=%d&offset=0&limit=10", _data.ID)); e != nil {
			t.Errorf("[ERROR] - Failed to get data %s", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		if resp, _, _ := GetReq(t, ts, "/news/?id=-999&offset=0&limit=10"); resp.StatusCode != http.StatusNotFound {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func GetAllData(t *testing.T) {
	tts := []TestTable{
		{
			name:               "Case: Positive Test",
			expectedStatusCode: http.StatusOK,
			expectedErr:        nil,
			errMsg:             "[ERROR] - Failed to get data",
			filter:             []map[string]interface{}{{"offset": 0, "limit": 10}},
		},
		{
			name:               "Case: Negative Test",
			expectedStatusCode: http.StatusBadRequest,
			expectedErr:        nil,
			errMsg:             fmt.Sprintf("[ERROR] - It should be error '%s'", helper.ErrDataNotFound.Error()),
			filter:             []map[string]interface{}{{"offset": -1, "limit": -10}},
		},
	}

	for _, tt := range tts {
		for _, filter := range tt.filter {
			t.Run(tt.name, func(t *testing.T) {
				if resp, _, e := GetReq(t, ts, fmt.Sprintf("/news?offset=%v&limit=%v", filter["offset"], filter["limit"])); e != tt.expectedErr && resp.StatusCode != tt.expectedStatusCode {
					t.Errorf("%s %s ", tt.errMsg, e.Error())
				}
			})
		}
	}
}

func makeRequest(t *testing.T, ts *httptest.Server, method, path string, body io.Reader) (*http.Response, string, error) {
	req, e := http.NewRequest(method, ts.URL+path, body)
	if e != nil {
		return nil, "", e
	}
	req.Header.Set("Content-Type", ContentTypeJson)

	var resp *http.Response
	switch method {
	case "GET":
		resp, e = http.DefaultTransport.RoundTrip(req)
	default:
		resp, e = http.DefaultClient.Do(req)
	}
	if e != nil {
		return nil, "", e
	}

	respBody, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		return nil, "", e
	}
	defer resp.Body.Close()

	return resp, string(respBody), nil
}
