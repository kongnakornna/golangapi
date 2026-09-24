FROM golang:1.26-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ai .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /ai /usr/local/bin/ai
EXPOSE 8001
ENTRYPOINT ["/usr/local/bin/ai"]
