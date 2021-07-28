package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync/atomic"
	"syscall"
	"time"
)

var counter int64

var version = "undefined"

var opts struct {
	listenAddress string
	responseBody  string
	responseCode  int
}

func main() {
	flag.StringVar(&opts.listenAddress, "a", ":2000", "listen address")
	flag.StringVar(&opts.responseBody, "b", "", "response body")
	flag.IntVar(&opts.responseCode, "c", http.StatusOK, "response status code")

	flag.Parse()

	if addrEnv := os.Getenv("LISTEN_ADDR"); addrEnv != "" {
		opts.listenAddress = addrEnv
	}
	if bodyEnv := os.Getenv("RESPONSE_BODY"); bodyEnv != "" {
		opts.responseBody = bodyEnv
	}
	if responseCodeEnv := os.Getenv("RESPONSE_CODE"); responseCodeEnv != "" {
		responseCodeEnvInt, errInt := strconv.Atoi(responseCodeEnv)
		if errInt != nil {
			fmt.Printf("response code must be an interger")
			return
		}
		opts.responseCode = responseCodeEnvInt
	}

	err := run()
	if err != nil {
		fmt.Printf("error run logserver, %s\n", err)
	}

	fmt.Printf("\ndone\n")
}

func run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	for _, pair := range os.Environ() {
		fmt.Printf("%s\n", pair)
	}

	fmt.Printf("logserver.%s listen %s\n", version, opts.listenAddress)

	ln, errListen := net.Listen("tcp4", opts.listenAddress)
	if errListen != nil {
		return fmt.Errorf("error listen address, %w", errListen)
	}
	defer ln.Close()

	server := &http.Server{
		Handler: http.HandlerFunc(handler),
	}

	go func() {
		<-ctx.Done()
		errShutdown := server.Shutdown(ctx)
		if errShutdown != nil && !errors.Is(errShutdown, context.Canceled) && !errors.Is(errShutdown, http.ErrServerClosed) {
			fmt.Printf("error shutdown server, %s", errShutdown)
		}
	}()

	errServe := server.Serve(ln)
	if errServe != nil && !errors.Is(errServe, http.ErrServerClosed) {
		return fmt.Errorf("error serve, %w", errServe)
	}

	return nil
}

func handler(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Add("Access-Control-Allow-Origin", "*")
	rw.Header().Add("Access-Control-Allow-Methods", "*")
	rw.Header().Add("Access-Control-Allow-Headers", "*")

	c := atomic.AddInt64(&counter, 1)

	fmt.Printf("___________[ %d ]___________\n", c)
	fmt.Printf("|  %v\n", time.Now())
	fmt.Printf("|  [%s] %s %s\n|  \n", req.RemoteAddr, req.Method, req.RequestURI)

	for key, values := range req.Header {
		fmt.Printf("|  %s: %v\n", key, values)
	}

	defer req.Body.Close()

	requestBody, err := io.ReadAll(req.Body)
	if err != nil {
		response(rw, fmt.Errorf("error read request body, %v", err))
		return
	}

	if req.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, errGzipReader := gzip.NewReader(bytes.NewReader(requestBody))
		if errGzipReader != nil {
			response(rw, fmt.Errorf("error init gzip reader, %v", errGzipReader))
			return
		}

		decodedBody, errRead := io.ReadAll(gzipReader)
		if errRead != nil {
			response(rw, fmt.Errorf("error decode gzip request body, %v", errRead))
			return
		}
		requestBody = decodedBody
	}

	if len(requestBody) != 0 {
		fmt.Printf("\n%s\n", requestBody)
	}

	response(rw, nil)
}

func response(rw http.ResponseWriter, err error) {
	if err != nil {
		fmt.Printf("[ERROR] %s\n", err)
	}

	rw.WriteHeader(opts.responseCode)

	if opts.responseBody != "" {
		_, errWrite := rw.Write([]byte(opts.responseBody))
		if errWrite != nil {
			fmt.Printf("[ERROR] error write response, %d", errWrite)
		}
	}

	fmt.Printf("\n")
}
