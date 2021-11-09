package main

import (
	"bytes"
	"compress/gzip"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync/atomic"
	"syscall"
	"time"
)

var counter int64
var out io.Writer = os.Stdout
var version = "undefined"

var opts struct {
	listenAddress string
	responseBody  string
	contentType   string
	responseCode  int
}

func main() {
	flag.StringVar(&opts.listenAddress, "a", ":2000", "listen address")
	flag.StringVar(&opts.responseBody, "b", "", "response body")
	flag.StringVar(&opts.contentType, "t", "text/plain; charset=utf-8", "content type header value")
	flag.IntVar(&opts.responseCode, "c", http.StatusOK, "response status code")

	flag.Parse()

	if addrEnv := os.Getenv("LISTEN_ADDR"); addrEnv != "" {
		opts.listenAddress = addrEnv
	}
	if bodyEnv := os.Getenv("RESPONSE_BODY"); bodyEnv != "" {
		opts.responseBody = bodyEnv
	}
	if contentTypeEnv := os.Getenv("CONTENT_TYPE"); contentTypeEnv != "" {
		opts.contentType = contentTypeEnv
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

	opts.responseBody = strings.Replace(opts.responseBody, "\\t", "\t", -1)
	opts.responseBody = strings.Replace(opts.responseBody, "\\n", "\n", -1)

	fmt.Printf("----- Options -----\nListen addr:\n%s\nResonse body:\n%s\nResponse code:\n%d\n----------\n", opts.listenAddress, opts.responseBody, opts.responseCode)

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

	rw.Header().Set("Content-Type", opts.contentType)
	rw.WriteHeader(opts.responseCode)

	if opts.responseBody != "" {
		body := []byte(opts.responseBody)
		if strings.HasPrefix(opts.responseBody, "file://") {
			var errReadBody error
			body, errReadBody = os.ReadFile(strings.TrimPrefix(opts.responseBody, "file://"))
			if errReadBody != nil {
				log.Printf("error read response body file %s, %v", opts.responseBody, errReadBody)
				return
			}
		}

		_, errWrite := rw.Write(body)
		if errWrite != nil {
			fmt.Printf("[ERROR] error write response, %d", errWrite)
		}
	}

	fmt.Printf("\n")
}
