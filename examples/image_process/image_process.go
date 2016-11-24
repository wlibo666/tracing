package main

import (
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/opentracing/opentracing-go"
	trace "github.com/wlibo666/tracing"
)

const (
	IMG_SERVER_ADDR      = ":1200"
	S3_PROXY_SERVER_ADDR = "http://127.0.0.1:1201/aksk"
	S3_SERVER_ADDR       = "http://127.0.0.1:1202/download"
	COLLETOR_ADDR        = ":1300"
)

func printInfo(r *http.Request) {
	log.Printf("time:[%s],method:[%s],url:[%s]\n", time.Now().String(), r.Method, r.URL)
}

func ProcessImage() {
	time.Sleep(time.Duration(rand.Int63n(time.Now().UnixNano())%1000) * time.Millisecond)
}

func FuncProcessImg(w http.ResponseWriter, r *http.Request) {
	printInfo(r)

	// start root span
	sp := opentracing.StartSpan(trace.GenOpName(r.URL.Path))
	defer sp.Finish()

	// get ak/sk from s3 proxy
	s3ProxyReq, err := http.NewRequest("GET", S3_PROXY_SERVER_ADDR, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = trace.InjectSpanToHttpHeader(sp, s3ProxyReq.Header)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	_, err = http.DefaultClient.Do(s3ProxyReq)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// download image from s3
	s3Req, err := http.NewRequest("GET", S3_SERVER_ADDR, nil)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	err = trace.InjectSpanToHttpHeader(sp, s3Req.Header)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	_, err = http.DefaultClient.Do(s3Req)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	// process
	imgProcessSpan := trace.StartChildSpan(trace.GenOpName("processimage"), sp)
	defer imgProcessSpan.Finish()
	ProcessImage()

	w.Write([]byte("OK\n"))
	return
}

func main() {
	err := trace.TracerInit("./trace.json")
	if err != nil {
		log.Fatalf("trace init failed,err:%s\n", err.Error())
	}

	http.HandleFunc("/img", FuncProcessImg)

	log.Fatal(http.ListenAndServe(IMG_SERVER_ADDR, nil))
}
