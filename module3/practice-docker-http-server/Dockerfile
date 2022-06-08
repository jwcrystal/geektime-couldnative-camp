# syntax=jwang/http-server:1
FROM golang:1.18-alpine

WORKDIR /practice-http-server

copy . ./
run echo "Install dependent modules" && \
    go mod download && go mod verify

RUN echo "Copy files" && \
   cd demo/ && \
   go build -o /http-server

EXPOSE 8080

CMD ["/http-server"]

