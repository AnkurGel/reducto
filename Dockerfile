FROM golang:1.16-alpine
WORKDIR /go/src/reducto
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o keygen cmd/reducto-keygen/main.go
RUN CGO_ENABLED=0 go build -o server cmd/reducto-server/main.go
