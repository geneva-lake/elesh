package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gocraft/web"
)

type LoginContext struct {
	*Context
	mc *MongoConnection
}

func InjectMiddlewareLogin(mc *MongoConnection) func(*LoginContext, web.ResponseWriter, *web.Request, web.NextMiddlewareFunc) {
	return func(ctx *LoginContext, w web.ResponseWriter, r *web.Request, next web.NextMiddlewareFunc) {
		ctx.mc = mc
		next(w, r)
	}
}

func (c *LoginContext) Index(w web.ResponseWriter, r *web.Request) {
	FileHandlerFromDir("login.html", w, r)
}

func (c *LoginContext) Login(w web.ResponseWriter, r *web.Request) {
	decoder := json.NewDecoder(r.Body)
	var pasw Password
	err := decoder.Decode(&pasw)
	if err != nil {
		log.Println("marshalling error", err)
		InternalError(&w, "marshalling error")
		return
	}
	password := pasw.Password
	login := pasw.Login
	psw, err := c.mc.getPassword()
	if err != nil {
		log.Println(err)
		InternalError(&w, "database error")
		return
	}
	if psw.Password == password && psw.Login == login {
		cookie := &http.Cookie{
			Name:    "session",
			Value:   "1234",
			Expires: time.Now().Local().Add(time.Hour * 24),
		}
		c.mc.setCookieExp(cookie.Expires)
		http.SetCookie(w, cookie)
		w.WriteHeader(http.StatusOK)
		authenticated = true
		fmt.Fprint(w, "ok")
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		authenticated = false
	}
}

func FileHandlerFromDir(file string, w web.ResponseWriter, r *web.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		return
	}
	f, err := webDir.Open(file)
	if err != nil {
		log.Fatalln("can not open file ", file, " ", err)
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		log.Fatalln("can not open file ", file, " ", err)
	}
	http.ServeContent(w, r.Request, file, fi.ModTime(), f)
}
