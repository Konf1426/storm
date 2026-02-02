# Build
FROM golang:1.24-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o messages .

# Run
FROM alpine:3.20
WORKDIR /app
COPY --from=build /app/messages /app/messages
EXPOSE 8081
ENV MESSAGES_ADDR=:8081
CMD ["/app/messages"]
