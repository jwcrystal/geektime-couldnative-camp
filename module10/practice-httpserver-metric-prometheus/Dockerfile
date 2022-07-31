FROM golang:1.18-alpine AS builder

#ENV GO111MODULE=off  \
ENV	CGO_ENABLED=0
#	GOOS=linux    \
#	GOARCH=amd64

WORKDIR /build
COPY . ./
RUN echo "Install dependent modules" && \
    go mod download && go mod verify && \
    cd demo/ && \
    go build -o http-server .

FROM busybox
COPY --from=builder /build/demo/http-server /
EXPOSE 8080
CMD ["/http-server"]
#ENTRYPOINT ["/http-server"]
