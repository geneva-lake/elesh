package main

import (
	"log"
	"net/http"

	"github.com/gocraft/web"
)

type AssetsContext struct {
	*Context
}

func (c *AssetsContext) Assets(w web.ResponseWriter, r *web.Request) {
	AssetsFileHandlerFromDir(w, r)
}

func AssetsFileHandlerFromDir(w web.ResponseWriter, r *web.Request) {
	if r.Method != "GET" && r.Method != "HEAD" {
		return
	}
	var file string = r.PathParams["*"] //r.URL.Path
	f, err := assetsDir.Open(file)
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
