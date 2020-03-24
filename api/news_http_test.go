// +build news_http

package api_test

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	. "github.com/rinosukmandityo/maknews/api"
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
	set url=root:Password.1@tcp(127.0.0.1:3306)/tes
	set timeout=10
	set db=tes
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
	elasticRepo repo.NewsRepository
	r           *chi.Mux
	ts          *httptest.Server
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
	newsRepo = rh.ChooseRepo()
	elasticRepo = rh.ElasticRepo()
	r = RegisterHandler()
}

func TestNewsHTTP(t *testing.T) {
	ts = httptest.NewServer(r)
	defer ts.Close()

	t.Run("Insert Data", InsertData)
	t.Run("Get All Data", GetAllData)
	t.Run("Update Data", UpdateData)
	t.Run("Delete Data", DeleteData)
	t.Run("Get Data", GetDataByID)
}

func PostData(t *testing.T, ts *httptest.Server, url string, _data m.News) error {
	dataBytes, e := getBytes(_data)
	if e != nil {
		return e
	}
	resp, _, e := makeRequest(t, ts, "POST", url, bytes.NewReader(dataBytes))
	if e != nil {
		return e
	}

	switch url {
	case "/news":
		if resp.StatusCode != http.StatusCreated {
			return errors.New("status should be 'Status Created' (201)")
		}
	default:
		if resp.StatusCode != http.StatusOK {
			return errors.New("status should be 'Status OK' (200)")
		}
	}

	return nil
}

func GetData(t *testing.T, ts *httptest.Server, url string, payload m.GetPayload, expected string) error {
	dataBytes, e := getPayloadBytes(payload)
	if e != nil {
		return e
	}
	resp, body, e := makeRequest(t, ts, "GET", url, bytes.NewReader(dataBytes))
	if e != nil {
		return e
	}
	if resp.StatusCode != http.StatusFound && strings.Contains(body, expected) {
		return errors.New("status should be 'Status Found' (302)")
	}

	return nil
}

func getBytes(_data m.News) ([]byte, error) {
	dataBytes, e := GetSerializer(ContentTypeJson).Encode(&_data)
	if e != nil {
		return dataBytes, e
	}
	return dataBytes, nil
}

func getPayloadBytes(_data m.GetPayload) ([]byte, error) {
	dataBytes, e := GetSerializer(ContentTypeJson).EncodeGetPayload(&_data)
	if e != nil {
		return dataBytes, e
	}
	return dataBytes, nil
}

func InsertData(t *testing.T) {
	newsService := logic.NewNewsService(newsRepo)
	elasticService := logic.NewElasticService(elasticRepo)

	testdata := ListTestData()
	wg := sync.WaitGroup{}

	// Clean test data if any
	for _, data := range testdata {
		wg.Add(1)
		go func(_data m.News) {
			newsService.Delete(&_data)
			elData := m.ElasticNews{
				ID:      _data.ID,
				Created: _data.Created,
			}
			elasticService.Delete(elData)
			wg.Done()
		}(data)
	}
	wg.Wait()

	t.Run("Case 1: Save data", func(t *testing.T) {
		for _, data := range testdata {
			wg.Add(1)
			go func(_data m.News) {
				if e := PostData(t, ts, "/news", _data); e != nil {
					t.Errorf("[ERROR] - Failed to save data %s ", e.Error())
				}
				wg.Done()
			}(data)
		}
		wg.Wait()

		time.Sleep(time.Second * 2)

		for _, data := range testdata {
			res, e := newsService.GetById(data.ID)
			if e != nil || res.ID == 0 {
				t.Errorf("[ERROR] - Failed to get data")
			}

			payload := m.GetPayload{
				Filter: map[string]interface{}{
					"id": data.ID,
				},
			}
			if elRes, e := elasticService.GetBy(payload); e != nil || len(elRes) == 0 {
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
		if e := PostData(t, ts, "/update", _data); e != nil {
			t.Errorf("[ERROR] - Failed to update data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := m.News{ID: -999}
		if e := PostData(t, ts, "/update", _data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func DeleteData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Delete data", func(t *testing.T) {
		_data := testdata[1]
		if e := PostData(t, ts, "/delete", _data); e != nil {
			t.Errorf("[ERROR] - Failed to delete data %s ", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		_data := testdata[1]
		if e := PostData(t, ts, "/delete", _data); e == nil {
			t.Error("[ERROR] - It should be error 'User Not Found'")
		}
	})
}

func GetDataByID(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get Data", func(t *testing.T) {
		_data := testdata[0]
		payload := m.GetPayload{
			Filter: map[string]interface{}{"id": _data.ID},
		}
		if e := GetData(t, ts, "/news", payload, _data.Author); e != nil {
			t.Errorf("[ERROR] - Failed to get data %s", e.Error())
		}
	})
	t.Run("Case 2: Negative Test", func(t *testing.T) {
		payload := m.GetPayload{
			Filter: map[string]interface{}{"id": -999},
		}
		if e := GetData(t, ts, "/news", payload, ""); e == nil {
			t.Error("[ERROR] - It should be error 'Data Not Found'")
		}
	})
}

func GetAllData(t *testing.T) {
	testdata := ListTestData()
	t.Run("Case 1: Get Data", func(t *testing.T) {
		_data := testdata[0]
		payload := m.GetPayload{
			Offset: 0,
			Limit:  10,
		}
		if e := GetData(t, ts, "/news", payload, _data.Author); e != nil {
			t.Errorf("[ERROR] - Failed to get data %s", e.Error())
		}
	})
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
