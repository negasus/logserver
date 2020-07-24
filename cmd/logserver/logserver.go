package main

import (
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"log"
	"net"
	"os"
	"sync/atomic"
	"time"
)

var counter int64
var version = "undefined"
var body = flag.String("b", "", "response body")
var addr = flag.String("a", ":2000", "listen address")

func main() {

	for _, pair := range os.Environ() {
		log.Printf(pair)
	}

	log.Printf("")

	flag.Parse()

	if addrEnv := os.Getenv("LISTEN_ADDR"); addrEnv != "" {
		*addr = addrEnv
	}
	if bodyEnv := os.Getenv("RESPONSE_BODY"); bodyEnv != "" {
		*body = bodyEnv
	}

	log.Printf("[INFO] logserver.%s listen %s", version, *addr)

	ln, err := net.Listen("tcp4", *addr)
	if err != nil {
		log.Printf("[ERROR] error listen, %v", err)
		os.Exit(1)
	}

	server := &fasthttp.Server{
		Name:    "logserver " + version,
		Handler: handler,
	}

	if err := server.Serve(ln); err != nil {
		log.Printf("[ERROR] error serve, %v", err)
		os.Exit(1)
	}
}

func handler(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
	ctx.Response.Header.Add("Access-Control-Allow-Methods", "*")
	ctx.Response.Header.Add("Access-Control-Allow-Headers", "*")

	c := atomic.AddInt64(&counter, 1)

	fmt.Printf("----------[ %d ]----------\n", c)
	fmt.Printf("%v\n", time.Now())
	fmt.Printf("[%s] %s %s\n\n", ctx.RemoteAddr(), ctx.Method(), ctx.RequestURI())
	ctx.Request.Header.VisitAll(func(key, value []byte) {
		fmt.Printf("%s: %s\n", key, value)
	})

	fmt.Printf("\n%s\n", ctx.PostBody())

	if *body != "" {
		ctx.Response.SetBodyString(*body)
	}
}
