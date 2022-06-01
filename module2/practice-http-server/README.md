# Module2 - a simple http server

1. build demo and launch it
```shell
$ go build -o demo

$ ./demo

/*
>>log 
http server :8080  started
2022/06/01 18:47:12 [200] 127.0.0.1 - / - 3.417µs
*/
```

2. Provide 2 API

- /
  - return IP of client
- /healthz
  - return data of header and specific variable ("Version") in env
  - log request on the terminal of server