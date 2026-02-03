run:
	go run ./...

tidy:
	go mod tidy

deps:
	go get github.com/go-chi/chi/v5
	go get github.com/prometheus/client_golang/prometheus/promhttp
	go get github.com/gorilla/websocket
	go get github.com/nats-io/nats.go

test:
	go test ./...
