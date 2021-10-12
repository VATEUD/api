FROM golang:1.17-alpine

WORKDIR ./

RUN go mod download

RUN go build -o /api /cmd/api/main.go

EXPOSE 3000
ENTRYPOINT [ "/api/main" ]