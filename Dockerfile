FROM golang:1.23.0-alpine as server
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN go build ./cmd/api/main.go
EXPOSE 8080
ENTRYPOINT [ "/app/main" ]