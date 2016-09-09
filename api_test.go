package main

import (
	"fmt"
	_ "log"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/gocraft/web"
)

func callerInfo() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return fmt.Sprintf("%s:%d", file, line)
}

func assertResponse(t *testing.T, rr *httptest.ResponseRecorder, body string, code int) {
	if gotBody := strings.TrimSpace(string(rr.Body.Bytes())); body != gotBody {
		t.Errorf("assertResponse: expected body to be %s but got %s. (caller: %s)", body, gotBody, callerInfo())
	}
	if code != rr.Code {
		t.Errorf("assertResponse: expected code to be %d but got %d. (caller: %s)", code, rr.Code, callerInfo())
	}
}

func assertToken(t *testing.T, rr *httptest.ResponseRecorder, code int, setToken bool, compare string) {
	token := rr.Header().Get("Token")
	if token == "" && !setToken {
		t.Errorf("assertToken: expected token to be not spare but got spare. (caller: %s)", callerInfo())
	}
	if setToken && token != compare {
		t.Errorf("assertToken: expected token to be %s but got %s. (caller: %s)", compare, token, callerInfo())
	}
	if code != rr.Code {
		t.Errorf("assertResponse: expected code to be %d but got %d. (caller: %s)", code, rr.Code, callerInfo())
	}
}

func TestSimpleRequest(t *testing.T) {
	var mc MongoConnection
	mc.InitTest()
	var ic InfluxConnection
	ic.InitTest()
	apiRouter := web.New(ApiContext{}).
		Middleware(InjectMiddlewareApi(&mc, &ic))
	apiRouter.Get("/", (*ApiContext).TestRequest)

	//Test simple request
	req, _ := http.NewRequest("GET", "/", nil)
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "Test", 200)
}

func TestCheckAuth(t *testing.T) {

	apiRouter := web.New(ApiContext{})
	apiRouter.Middleware((*ApiContext).CheckAuth)
	apiRouter.Get("/", (*ApiContext).TestRequest)

	//Check pass authentication
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("ApiKey", "123456789")
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "Test", 200)

	//Check not passed authentication
	req, _ = http.NewRequest("GET", "/", nil)
	rw = httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "Access denied", 401)
}

func TestCheckSession(t *testing.T) {

	var mc MongoConnection
	mc.InitTest()
	mc.DeleteSessions()
	var ic InfluxConnection
	ic.InitTest()
	apiRouter := web.New(ApiContext{})
	apiRouter.Middleware((*ApiContext).CheckAuth).
		Middleware(InjectMiddlewareApi(&mc, &ic)).
		Middleware((*ApiContext).CheckSession)
	apiRouter.Get("/", (*ApiContext).MiddleTest)

	//Test begin new session
	req, _ := http.NewRequest("GET", "/", nil)
	req.Header.Set("ApiKey", "123456789")
	req.Header.Set("Device-Id", "1234")
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	assertToken(t, rw, 200, false, "")

	//Test continue session
	mc.DeleteSessions()
	mc.setSession("1234", "asdf")
	req, _ = http.NewRequest("GET", "/", nil)
	rw = httptest.NewRecorder()
	req.Header.Set("ApiKey", "123456789")
	req.Header.Set("Device-Id", "1234")
	req.Header.Set("Token", "asdf")
	apiRouter.ServeHTTP(rw, req)
	assertToken(t, rw, 200, true, "asdf")

}

func TestCheckTrial(t *testing.T) {

	var mc MongoConnection
	mc.InitTest()
	mc.DeleteDevices()
	var ic InfluxConnection
	ic.InitTest()
	apiRouter := web.New(ApiContext{})
	apiRouter.Middleware((*ApiContext).CheckAuth).
		Middleware(InjectMiddlewareApi(&mc, &ic))
	apiRouter.Get("/trial", (*ApiContext).CheckTrial)

	//check trial time when device is not exist
	req, _ := http.NewRequest("GET", "/trial", nil)
	req.Header.Set("ApiKey", "123456789")
	req.Header.Set("Device-Id", "1234")
	rw := httptest.NewRecorder()
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "True", 200)

	//check trial time when device just had created
	req, _ = http.NewRequest("GET", "/trial", nil)
	rw = httptest.NewRecorder()
	req.Header.Set("ApiKey", "123456789")
	req.Header.Set("Device-Id", "1234")
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "True", 200)

	//check trial time when device created 4 days ago
	mc.DeleteDevices()
	install := time.Now().AddDate(0, 0, -4)
	mc.CreateDevice("1234", install)
	req, _ = http.NewRequest("GET", "/trial", nil)
	rw = httptest.NewRecorder()
	req.Header.Set("ApiKey", "123456789")
	req.Header.Set("Device-Id", "1234")
	apiRouter.ServeHTTP(rw, req)
	assertResponse(t, rw, "False", 200)
}
