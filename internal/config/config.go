package config

import (
	"flag"
	"net/url"
)

var (
	ServerAddr     = flag.String("a", ":8080", "server address")
	ResultBaseAddr = flag.String("b", "http://localhost:8080/", "result base address of shortened URL")
)

func Init() {
	flag.Parse()
	_, err := url.ParseRequestURI(*ResultBaseAddr)
	if err != nil {
		panic(err)
	}
}
