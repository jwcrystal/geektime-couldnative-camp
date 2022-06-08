package main

import (
	"fmt"
	"handlers"
	"log"
	"net/http"
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
	engine := handlers.New()

	engine.GET("/", func(w http.ResponseWriter, res *http.Request) {
		fmt.Fprintln(w, "Hello. CNCamp")
		fmt.Fprintln(w, handlers.GetClientIP(res))
		fmt.Fprintln(w, res.Host)

	})

	engine.GET("/healthz", handlers.TestHandler)

	port := ":" + strconv.Itoa(8080)
	fmt.Println("http server", port, " started")
	err := engine.Run(port)
	if err != nil {
		log.Fatal("ListenAndServe: ", err.Error())
		return
	}
}
