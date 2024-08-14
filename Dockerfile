FROM golang:1.22.5-alpine
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . ./
RUN go run ./cmd/migrations/main.go -up
RUN go build ./cmd/api/main.go -o /bin/server
EXPOSE 8080
ENTRYPOINT [ "/bin/server" ]