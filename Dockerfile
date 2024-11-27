FROM golang:1.22.2

COPY . ./app
WORKDIR ./app

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o cmd/server/server cmd/server/main.go

CMD ["cmd/server/server"]