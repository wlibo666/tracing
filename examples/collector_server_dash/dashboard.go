package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"

	"sourcegraph.com/sourcegraph/appdash"
	"sourcegraph.com/sourcegraph/appdash/traceapp"
)

const (
	IMG_SERVER_ADDR      = ":1200"
	S3_PROXY_SERVER_ADDR = "http://127.0.0-.1:1201/aksk"
	S3_SERVER_ADDR       = "http://127.0.0.1:1202/download"
	COLLETOR_ADDR        = ":1300"

	DASHBOARD_ADDR_IPV4 = "10.58.60.254"
	DASHBOARD_ADDR_PORT = 8800
)

func startDashboard() error {
	// malloc memory and store traces by recving from 10.58.60.254:1300
	store := appdash.NewMemoryStore()
	l, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 1300})
	if err != nil {
		log.Fatalf("listen tcp failed,err:%s\n", err.Error())
		return err
	}
	cs := appdash.NewServer(l, appdash.NewLocalCollector(store))
	go func() {
		cs.Start()
		log.Fatalf("collector store exit\n")
	}()

	// start dashboard
	appdashUrlStr := fmt.Sprintf("http://%s:%d", DASHBOARD_ADDR_IPV4, DASHBOARD_ADDR_PORT)
	fmt.Fprintf(os.Stdout, "dashboard addr:%s\n", appdashUrlStr)
	appdashUrl, err := url.Parse(appdashUrlStr)
	if err != nil {
		log.Fatalf("parse url failed,err:%s\n", err.Error())
		return err
	}
	tapp, err := traceapp.New(nil, appdashUrl)
	if err != nil {
		log.Fatalf("new trace app failed:%s\n", err.Error())
		return err
	}
	// 指定tracer控制台的数据源
	tapp.Store = store
	tapp.Queryer = store

	log.Fatal(http.ListenAndServe("10.58.60.254:8800", tapp))
	return nil
}

func main() {
	startDashboard()
}
