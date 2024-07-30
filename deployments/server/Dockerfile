FROM golang:1.22.1 AS builder

WORKDIR /

COPY go.mod .
COPY go.sum .
RUN  go mod download

COPY . .

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X 'github.com/ivan-bokov/pow-ddos/internal/version.Version=v1.0.0'" -o word-of-wisdom-server cmd/server/main.go

FROM scratch

COPY --from=builder /word-of-wisdom-server /

COPY data /data

ENTRYPOINT ["/word-of-wisdom-server"]