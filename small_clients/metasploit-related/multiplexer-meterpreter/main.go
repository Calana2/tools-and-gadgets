package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"github.com/gorilla/mux"
)

var (
  RHOST_1   = "http://192.168.137.5:8000"
  RHOST_2   = "http://192.168.137.5:9000"
	hostProxy = make(map[string]string)
	proxies   = make(map[string]*httputil.ReverseProxy)
)

func init() {
	hostProxy["attacker1.com"] = RHOST_1
	hostProxy["attacker2.com"] = RHOST_2

	for k, v := range hostProxy {
		remote, err := url.Parse(v)
		if err != nil {
			log.Fatal("Unable to parse proxy target")
		}
		proxies[k] = httputil.NewSingleHostReverseProxy(remote)
	}
}

func main() {
	r := mux.NewRouter()
	for host, proxy := range proxies {
		r.Host(host).Handler(proxy)
	}
	log.Fatal(http.ListenAndServe(":80", r))
}
