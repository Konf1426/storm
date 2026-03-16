# Build
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o gateway .

# Run
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/gateway /app/gateway
EXPOSE 8080
ENV GATEWAY_ADDR=:8080
CMD ["/app/gateway"]
