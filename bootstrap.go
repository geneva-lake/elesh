package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gocraft/web"
	"github.com/robfig/cron"
)

var currentRoot string
var webDir http.FileSystem
var assetsDir http.FileSystem

var mongoConnection MongoConnection
var influxConnection InfluxConnection

var apiKey = "123456789"

type Context struct {
}

func (c *Context) ErrorMiddleware(rw web.ResponseWriter, req *web.Request, next web.NextMiddlewareFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.Header().Set("Content-Type", "text/html")
			rw.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(rw, "Internal server error")
		}
	}()
	next(rw, req)
}

func InternalError(rwl *web.ResponseWriter, mes string) {
	(*rwl).WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(*rwl, mes)
}

func pleaseAuth(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("Access denied"))
}

func bootstrap() {
	currentRoot, _ := os.Getwd()
	webDir = http.Dir(currentRoot + "\\assets\\web")
	assetsDir = http.Dir(currentRoot + "\\assets")
	dbl, err := os.OpenFile("errors.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("Error opening/creating log file")
		os.Exit(1)
	}
	defer dbl.Close()
	log.SetOutput(dbl)
	var mc MongoConnection
	err = mc.Init()
	if err != nil {
		log.Fatalln("mongo init err ", err)
	}
	defer mc.SessionClose()
	var ic InfluxConnection
	err = ic.Init()
	if err != nil {
		log.Fatalln("influx init err ", err)
	}
	ic.mc = mc

	c := cron.New()
	c.AddFunc("@daily", func() { ic.ApiCount() })
	c.Start()

	rootRouter := web.New(Context{}).
		Middleware((*Context).ErrorMiddleware)
	apiRouter := rootRouter.Subrouter(ApiContext{}, "/api").
		Middleware(InjectMiddlewareApi(&mc, &ic))
	apiRouter.Middleware((*ApiContext).CheckAuth).
		Middleware((*ApiContext).CheckSession).
		Middleware((*ApiContext).CountRequest)
	apiRouter.Get("/:id", (*ApiContext).GetDeviceById)
	apiRouter.Post("/", (*ApiContext).CreateDeviceWithId)
	apiRouter.Get("/all", (*ApiContext).GetDevices)
	apiRouter.Get("/trial", (*ApiContext).CheckTrial)

	loginRouter := rootRouter.Subrouter(LoginContext{}, "/web").
		Middleware(InjectMiddlewareLogin(&mc))
	loginRouter.Get("/", (*LoginContext).Index)
	loginRouter.Post("/", (*LoginContext).Login)
	webRouter := rootRouter.Subrouter(WebContext{}, "/admin").
		Middleware(InjectMiddlewareWeb(&mc, &ic, &authenticated))

	webRouter.Middleware((*WebContext).CheckSession)
	webRouter.Get("/devices/:skip/:limit/:filter/:order", (*WebContext).GetDevices)
	webRouter.Post("/use-count", (*WebContext).GetUseCount)
	webRouter.Get("/site", (*WebContext).GetClass)
	webRouter.Post("/site", (*WebContext).SetClass)

	assetsRouter := rootRouter.Subrouter(AssetsContext{}, "/assets")
	assetsRouter.Get("/:*", (*AssetsContext).Assets)
	http.ListenAndServe("localhost:3000", rootRouter)
}
