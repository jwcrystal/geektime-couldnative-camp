package handlers

import (
	"net/http"
	"os"
	"testing"
)

func TestEnv(t *testing.T) {
	key := "Version"
	_, ok := os.LookupEnv(key)
	if !ok {
		t.Fatalf("%s not set\n", key)
	}
	//for _, e := range os.Environ() {
	//	pair := strings.SplitN(e, "=", 2)
	//	fmt.Println(pair[0])
	//}
}

func TestLogHandler(t *testing.T) {
	header := http.Header{"Content-Type": {"text/plains"}}
	r := &http.Request{
		Method:        "GET",
		Header:        header,
		ContentLength: 10,
		Host:          "localhost:8080",
	}
	LogHandler(200, r)
}
