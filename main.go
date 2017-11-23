package main

import (
	"./news"
	"./util"
	"./web"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/robfig/cron"
	"log"
	"net/http"
	"os"
)

func init() {
	util.ShowLog()
}

func maino() {
	fs := http.FileServer(http.Dir("."))
	http.Handle("/", fs)

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	http.ListenAndServe(bind, nil)
}

func main() {
	// CRON
	c := cron.New()
	c.AddFunc("@every 2h", func() { news.ReadNHK() })
	c.Start()

	// MAIN
	r := mux.NewRouter()
	for _, h := range web.Handlers {
		if len(h.Method) > 0 {
			r.HandleFunc(h.Path, h.Fn).Methods(h.Method)
		} else {
			r.HandleFunc(h.Path, h.Fn)
		}
	}

	bind := fmt.Sprintf("%s:%s", os.Getenv("HOST"), os.Getenv("PORT"))
	util.Info.Printf("listening on %s...\n", bind)
	err := http.ListenAndServe(bind, AuthHandler()(r))
	if err != nil {
		panic(err)
	}
}

type authHandler struct {
	handler http.Handler
	logger  *log.Logger
}

func AuthHandler() func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		a := &authHandler{handler: h}
		a.logger = util.Info
		return a
	}
}

func (h authHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	pass := util.AuthenticateRequest(req)

	if !pass {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprint(w, "Access to myapp is Forbidden !!")
		return
	}

	h.handler.ServeHTTP(w, req)
}
