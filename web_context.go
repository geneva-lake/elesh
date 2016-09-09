package main

import (
	"encoding/json"
	"fmt"
	"log"

	"strconv"
	"time"

	"github.com/gocraft/web"
)

var authenticated bool = false

type WebContext struct {
	*Context
	mc            *MongoConnection
	ic            *InfluxConnection
	filters       map[string]string
	authenticated bool
}

//Injecting databases to WebContext
func InjectMiddlewareWeb(mc *MongoConnection, ic *InfluxConnection, auth *bool) func(*WebContext, web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {
	return func(ctx *WebContext, w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		ctx.mc = mc
		ctx.ic = ic
		ctx.authenticated = *auth
		next(w, r)
	}
}

func (c *WebContext) TestRequest(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Test")
}

func (c *WebContext) CheckSession(w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
	if c.authenticated {
		_, err := r.Cookie("session")
		if err != nil {
			log.Println("cookieSession err ", err)
			pleaseAuth(w)
			return
		} else {
			expire, err := c.mc.getCookieExp()
			if err != nil {
				log.Println("mongo web err ", err)
				InternalError(&w, "databse error")
				return
			}
			if (*expire).Sub(time.Now().Local()) < 0 {
				pleaseAuth(w)
				return
			}
			next(w, r)
		}
	} else {
		pleaseAuth(w)
	}
}

//Get list of installed devices
func (c *WebContext) GetDevices(rw web.ResponseWriter, req *web.Request) {
	skip, _ := strconv.ParseInt(req.PathParams["skip"], 10, 16)
	limit, _ := strconv.ParseInt(req.PathParams["limit"], 10, 16)
	order, _ := strconv.ParseInt(req.PathParams["order"], 10, 16)
	filter := req.PathParams["filter"]
	total, err := c.mc.getDevicesCount()
	if err != nil {
		log.Println("web err ", err)
		InternalError(&rw, "databse error")
		return
	}
	devices, err := c.mc.getAllDevices(int(skip), int(limit), filter, int(order))
	jDev := Devices{}
	if err != nil {
		log.Println("web err ", err)
		InternalError(&rw, "databse error")
		return
	}
	jDev.Devices = devices
	jDev.Total = total
	jsonM, err := json.Marshal(jDev)
	if err != nil {
		log.Println("marshaling error ", err)
		InternalError(&rw, "marshaling error")
		return
	}
	fmt.Fprint(rw, string(jsonM))
}

//Get info about some class
func (c *WebContext) GetClass(rw web.ResponseWriter, req *web.Request) {
	site, err := c.mc.getSite()
	if err != nil {
		if err.Error() == "not found" {
			var siteNew Site
			siteNew.Date = time.Now()
			siteNew.Url = "name"
			err = c.mc.setSite(siteNew)
			fmt.Println("set site")
			if err != nil {
				log.Println("set siteerr ", err)
				InternalError(&rw, "databse error")
			}
			jsonM, err := json.Marshal(siteNew)
			if err != nil {
				log.Println("marshaling error ", err)
				InternalError(&rw, "marshaling error")
				return
			}
			fmt.Fprint(rw, string(jsonM))
		} else {
			log.Println(err)
			InternalError(&rw, "databse error")
			return
		}
	}
	jsonM, err := json.Marshal(site)
	if err != nil {
		log.Println("marshaling error ", err)
		InternalError(&rw, "marshaling error")
		return
	}
	fmt.Fprint(rw, string(jsonM))
}

func (c *WebContext) SetClass(rw web.ResponseWriter, req *web.Request) {
	decoder := json.NewDecoder(req.Body)
	var site Site
	err := decoder.Decode(&site)
	if err != nil {
		log.Println("read error ", err)
		InternalError(&rw, "read error")
		return
	}
	err = c.mc.setSite(site)
	if err != nil {
		log.Println(err)
		InternalError(&rw, "databse error")
	}
}

//Getting count of requesting of api
func (c *WebContext) GetUseCount(rw web.ResponseWriter, req *web.Request) {
	decoder := json.NewDecoder(req.Body)
	var postDate PostDate
	err := decoder.Decode(&postDate)
	if err != nil {
		log.Println("read error ", err)
		InternalError(&rw, "read error")
		return
	}
	counts, err := c.mc.getUseCount(postDate.Begin, postDate.End)
	if err != nil {
		log.Println(err)
		InternalError(&rw, "databse error")
		return
	}
	var dates []int = make([]int, len(*counts))
	var myCount []int = make([]int, len(*counts))
	for i, count := range *counts {
		dates[i] = count.Date.Day()
		myCount[i] = count.Count
	}
	var dto UsePerDayDto
	dto.Count = myCount
	dto.Date = dates
	jsonM, err := json.Marshal(dto)
	if err != nil {
		log.Println(err)
		InternalError(&rw, "marshaling error")
		return
	}
	fmt.Fprint(rw, string(jsonM))
}

func (c *WebContext) MiddleTest(rw web.ResponseWriter, req *web.Request) {
	fmt.Fprint(rw, "Test")
}
