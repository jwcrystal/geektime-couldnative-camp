package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"handlers"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

//$ curl http://localhost:8080
//Hello. CNCamp

//$ curl http://localhost:8080/healthz
//<h1>Hello, cncamp</h1>
//<h2>Header of Request: </h2>
//Accept= */*
//User-Agent= curl/7.79.1
//<h2>Env Variable: </h2>
//Version = 3.0

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	eg, egCtx := errgroup.WithContext(ctx)

	// signal listening
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt) // = syscall.SIGINT
	// web
	engine := handlers.New()

	// start http server
	eg.TryGo(func() error {
		return StartHttpServer(engine)
	})
	eg.TryGo(func() error {
		// block until receiving the signal of egCtx was closed
		<-egCtx.Done()
		log.Println("http server stop")
		return engine.Shutdown(ctx)
	})
	// listen signal: kill -9 or ctrl+c
	eg.TryGo(func() error {
		for {
			select {
			case <-egCtx.Done():
				return egCtx.Err()
			case <-sig:
				cancel() // ctrl+c or kill -9
			}
		}
	})
	if err := eg.Wait(); err != nil {
		fmt.Println("error: ", err)
	} else {
		fmt.Println("done successfully")
	}
}

func StartHttpServer(engine *handlers.Engine) error {
	engine.GET("/", func(w http.ResponseWriter, res *http.Request) {
		fmt.Fprintln(w, "Hello. CNCamp")
		fmt.Fprintln(w, handlers.GetClientIP(res))
		fmt.Fprintln(w, res.Host)

	})

	engine.GET("/healthz", handlers.TestHandler)

	port := ":" + strconv.Itoa(8080)
	log.Println("http server", port, " started")
	err := engine.Run(port)
	return err
}
