package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gocraft/web"
)

var test_auth bool = true

func TestSimpleRequestWeb(t *testing.T) {
	var mc MongoConnection
	mc.InitTest()
	var ic InfluxConnection
	ic.InitTest()
	apiRouter := web.New(WebContext{}).
		Middleware(InjectMiddlewareWeb(&mc, &ic, &test_auth))
	apiRouter.Get("/", (*WebContext).TestRequest)

	//Test simple request
	req, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "Test", 200)
}

func TestGetDevices(t *testing.T) {
	var mc MongoConnection
	mc.InitTest()
	mc.DeleteDevices()
	var ic InfluxConnection
	ic.InitTest()
	apiRouter := web.New(WebContext{})
	apiRouter.Middleware(InjectMiddlewareWeb(&mc, &ic, &test_auth))
	apiRouter.Get("/devices/:skip/:limit/:filter/:order", (*WebContext).GetDevices)
	for i := 0; i < 10; i++ {
		install := time.Now().AddDate(0, 0, i)
		mc.CreateDevice("1234"+strconv.Itoa(i), install)
	}
	skip := 0
	limit := 5
	filter := "device-id"
	order := -1
	url := fmt.Sprintf("/devices/%d/%d/%s/%d", skip, limit, filter, order)
	req, _ := http.NewRequest("GET", url, nil)
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	decoder := json.NewDecoder(rw.Body)
	var results Devices
	err := decoder.Decode(&results)
	if err != nil {
		t.Errorf("json decoding error ", err)
	}
	if results.Total != 10 {
		t.Errorf("Total is not equal")
	}
	devices, _ := mc.getAllDevices(skip, limit, filter, order)
	if (*devices)[0].DeviceId != (*(results.Devices))[0].DeviceId {
		t.Errorf("Order is not correct")
	}

}

func TestGetUseCount(t *testing.T) {
	var mc MongoConnection
	mc.InitTest()
	mc.DeleteCounts()
	var ic InfluxConnection
	ic.InitTest()
	apiRouter := web.New(WebContext{})
	apiRouter.Middleware(InjectMiddlewareWeb(&mc, &ic, &test_auth))
	apiRouter.Post("/use-count", (*WebContext).GetUseCount)
	for i := 0; i < 10; i++ {
		date := time.Now().AddDate(0, 0, i*-1)
		mc.WriteUseCount(i, date)
	}
	var postDate PostDate
	postDate.Begin = time.Now().AddDate(0, 0, -10)
	postDate.End = time.Now()
	jsonM, err := json.Marshal(postDate)
	if err != nil {
		t.Errorf("marshal err ", err)
	}
	req, _ := http.NewRequest("POST", "/use-count", strings.NewReader(string(jsonM)))
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	decoder := json.NewDecoder(rw.Body)
	var dto UsePerDayDto
	err = decoder.Decode(&dto)
	if err != nil {
		t.Errorf("json decoding error ", err)
	}
	if dto.Count[0] != 0 {
		t.Errorf("Count is not correct")
	}
}

func TestCheckSessionWeb(t *testing.T) {
	var mc MongoConnection
	mc.InitTest()
	mc.DeleteSessions()
	var ic InfluxConnection
	ic.InitTest()
	apiRouter := web.New(WebContext{})
	apiRouter.Middleware(InjectMiddlewareWeb(&mc, &ic, &test_auth)).
		Middleware((*WebContext).CheckSession)
	apiRouter.Get("/", (*WebContext).MiddleTest)
	cookie := &http.Cookie{
		Name:    "session",
		Value:   "1234",
		Expires: time.Now().Local().Add(time.Hour * 24),
	}
	mc.setCookieExp(cookie.Expires)

	//Test begin new session
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Add("Cookie", (*cookie).String())
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "Test", 200)
}
