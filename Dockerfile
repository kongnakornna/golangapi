# เปลี่ยนจาก 1.22 เป็น 1.26 ให้ตรงกับ go.mod
FROM golang:1.26-alpine

WORKDIR /app

RUN apk add --no-cache git bash curl

RUN go install github.com/air-verse/air@v1.52.3 && \
    go install github.com/roerohan/wait-for-it@v0.2.13

ENV PATH=$PATH:/go/bin

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 5000
CMD ["air", "-c", ".air.toml", "serve"]