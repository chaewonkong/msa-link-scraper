FROM golang:1.22 as builder

WORKDIR /go/src/app

COPY . .


RUN go mod tidy && \
    go mod vendor && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o bin/app .

FROM scratch

WORKDIR /go/src/app

COPY --from=builder /go/src/app/bin/app ./app

ENTRYPOINT ["./app"]