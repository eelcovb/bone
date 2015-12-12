package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/go-zoo/bone"
)

var mux = bone.New()

func main() {

	mux.GetFunc("/", defaultHandler)
	mux.GetFunc("/reg/#var^[a-z]$/#var2^[0-9]$", ShowVar)
	mux.GetFunc("/test", defaultHandler)
	mux.Get("/file/", http.StripPrefix("/file/", http.FileServer(http.Dir("assets"))))

	http.ListenAndServe(":8080", mux)
}

func defaultHandler(rw http.ResponseWriter, req *http.Request) {
	file, _ := ioutil.ReadFile("index.html")
	rw.Write(file)
}

func ShowVar(rw http.ResponseWriter, req *http.Request) {
	value := mux.GetAllValues(req)
	fmt.Println(value)
	//rw.Write([]byte(value))
}
