FROM golang:1.21.0-alpine3.18 as builder

WORKDIR /usr/local/src

COPY . .

RUN go mod download
RUN go build -o cmd/bin/testtask cmd/main.go

FROM alpine:latest

COPY --from=builder /usr/local/src/cmd/bin/testtask /
COPY configs/configs.yaml configs/configs.yaml

ENTRYPOINT [ "/testtask" ] 
