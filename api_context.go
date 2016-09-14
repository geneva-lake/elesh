package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gocraft/web"
)

type ApiContext struct {
	*Context
	token string
	mc    *MongoConnection
	ic    *InfluxConnection
}

type PostDevice struct {
	DeviceId string `json:"device-id"`
}

func InjectMiddlewareApi(mc *MongoConnection, ic *InfluxConnection) func(*ApiContext,
	web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {
	return func(ctx *ApiContext, w web.ResponseWriter, r *web.Request,
		next web.NextMiddlewareFunc) {
		ctx.mc = mc
		ctx.ic = ic
		next(w, r)
	}
}

func (c *ApiContext) TestRequest(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Test")
}

func (c *ApiContext) MiddleTest(rw web.ResponseWriter, req *web.Request) {
	rw.Header().Add("Token", c.token)
}

//Checking api key in request
func (c *ApiContext) CheckAuth(w web.ResponseWriter, r *web.Request,
	next web.NextMiddlewareFunc) {
	auth := r.Header.Get("ApiKey")
	if auth == "" {
		pleaseAuth(w)
		return
	} else if auth != apiKey {
		pleaseAuth(w)
		return
	}
	next(w, r)
}

//Checking session token
func (c *ApiContext) CheckSession(rw web.ResponseWriter, r *web.Request,
	next web.NextMiddlewareFunc) {
	id := r.Header.Get("Device-Id")
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "No device id")
		return
	}
	token := r.Header.Get("Token")
	if token == "" {
		c.mc.DeleteSession(id)
		token = c.GenerateToken()
		deviceSession, err := c.mc.setSession(id, token)
		if err != nil {
			log.Println(err)
			InternalError(&rw, "database error")
			return
		}
		c.token = deviceSession.Token
		next(rw, r)
	} else {
		deviceSession, err := c.mc.getSession(id)
		if err != nil {
			log.Println("ses err ", err)
			token = c.GenerateToken()
			deviceSession, err = c.mc.setSession(id, token)
			if err != nil {
				log.Println(err)
				InternalError(&rw, "database error")
				return
			}
		}
		c.token = deviceSession.Token
		next(rw, r)
	}
}

//Check trial time of installed device
func (c *ApiContext) CheckTrial(rw web.ResponseWriter, req *web.Request) {
	deviceId := req.Header.Get("Device-Id")
	if deviceId == "" {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "No device id")
		return
	}
	rw.Header().Add("Token", c.token)
	device, err := c.mc.getDeviceById(deviceId)
	if err != nil {
		if err.Error() == "not found" {
			c.mc.CreateDevice(deviceId, time.Now())
			fmt.Fprint(rw, "True")
			return
		} else {
			InternalError(&rw, "database error")
			return
		}
	}
	var trialTime time.Duration = time.Time.Sub(device.InstallDate, time.Now())
	days, _ := time.ParseDuration("-72h")
	if trialTime > days {
		fmt.Fprint(rw, "True")
	} else {
		fmt.Fprint(rw, "False")
	}
}

func (c *ApiContext) GetDeviceById(rw web.ResponseWriter, req *web.Request) {
	deviceId := req.PathParams["id"]
	device, err := c.mc.getDeviceById(deviceId)
	if err != nil {
		log.Println("contr err ", err)
		InternalError(&rw, "")
		return
	}
	jsonM, err := json.Marshal(device)
	if err != nil {
		log.Println(err)
		InternalError(&rw, "")
		return
	}
	fmt.Fprint(rw, string(jsonM))
}

func (c *ApiContext) GetDevices(rw web.ResponseWriter, req *web.Request) {
	devices, err := c.mc.getAllDevices(0, 10, "InstallDate", -1)
	if err != nil {
		log.Println(err)
		InternalError(&rw, "")
		return
	}
	jsonM, err := json.Marshal(devices)
	if err != nil {
		log.Println(err)
		InternalError(&rw, "")
		return
	}
	fmt.Fprint(rw, string(jsonM))
}

func (c *ApiContext) CreateDeviceWithId(rw web.ResponseWriter, req *web.Request) {
	decoder := json.NewDecoder(req.Body)
	var postDevice PostDevice
	err := decoder.Decode(&postDevice)
	if err != nil {
		log.Println("read err", err)
		InternalError(&rw, "")
		return
	}
	c.mc.CreateDevice(postDevice.DeviceId, time.Now())
	fmt.Fprint(rw, "Created")
}

//Registering and counting api requests
func (c *ApiContext) CountRequest(rw web.ResponseWriter, r *web.Request,
	next web.NextMiddlewareFunc) {
	id := r.Header.Get("Device-Id")
	if id == "" {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "No device id")
		return
	}
	token := r.Header.Get("Token")
	err := c.ic.Write(id, token)
	if err != nil {
		log.Println("count err ", err)
	}
	next(rw, r)
}

func (c *ApiContext) GenerateToken() string {
	var token []byte
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	rnd.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 32; i++ {
		token = append(token, byte(rnd.Int63n(74)+48))
	}
	return string(token)
}
