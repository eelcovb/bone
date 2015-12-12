package main

import (
	"fmt"
	"net/http"

	"github.com/go-zoo/bone"
	"github.com/gorilla/mux"
	"github.com/julienschmidt/httprouter"
)

func main() {
	muxx := bone.New().Prefix("/api")
	boneSub := bone.New()
	gorrilaSub := mux.NewRouter()
	httprouterSub := httprouter.New()

	boneSub.GetFunc("/test/:var", func(rw http.ResponseWriter, req *http.Request) {
		val := muxx.GetValue(req, "var")
		rw.Write([]byte(fmt.Sprintf("Hello %s from bone mux", val)))
	})

	gorrilaSub.HandleFunc("/test", func(rw http.ResponseWriter, req *http.Request) {
		rw.Write([]byte("Hello from gorilla mux"))
	})

	httprouterSub.GET("/test", func(rw http.ResponseWriter, req *http.Request, _ httprouter.Params) {
		rw.Write([]byte("Hello from httprouter mux"))
	})

	muxx.SubRoute("/bone", boneSub)
	muxx.SubRoute("/gorilla", gorrilaSub)
	muxx.SubRoute("/http", httprouterSub)

	http.ListenAndServe(":8080", muxx)
}
