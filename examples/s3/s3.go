package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	trace "github.com/wlibo666/tracing"
)

const (
	IMG_SERVER_ADDR      = ":1200"
	S3_PROXY_SERVER_ADDR = "http://127.0.0.1:1201/aksk"
	S3_SERVER_ADDR       = "10.58.60.254:1202"
	COLLETOR_ADDR        = ":1300"
)

func printInfo(r *http.Request) {
	log.Printf("time:[%s],method:[%s],url:[%s]\n", time.Now().String(), r.Method, r.URL)
}

func downLoad() {
	time.Sleep(time.Duration(rand.Int63n(time.Now().UnixNano())%1000) * time.Millisecond)
}

func DownLoadResouse(w http.ResponseWriter, r *http.Request) {
	printInfo(r)
	sp := trace.StartSpanFromHttpHeader(trace.GenOpName(r.URL.Path), r.Header)
	defer sp.Finish()

	downLoad()

	return
}

func main() {
	err := trace.TracerInit("./trace.json")
	if err != nil {
		log.Fatalf("trace init failed,err:%s\n", err.Error())
	}

	http.HandleFunc("/download", DownLoadResouse)
	log.Fatal(http.ListenAndServe(":1202", nil))
}
