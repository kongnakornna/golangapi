# ขั้นตอนการ build
FROM golang:1.21-alpine AS builder

# ติดตั้ง dependencies ที่จำเป็นสำหรับการ build
RUN apk add --no-cache git

WORKDIR /app

# คัดลอกไฟล์ go mod
COPY go.mod go.sum ./
RUN go mod download

# คัดลอกซอร์สโค้ด
COPY . .

# Build แอปพลิเคชัน
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o main cmd/app/main.go

# ขั้นตอนการรัน
FROM alpine:latest

# ติดตั้ง dependencies ที่จำเป็นสำหรับการรัน
RUN apk --no-cache add ca-certificates tzdata

# ตั้งค่าโซนเวลา
ENV TZ=Asia/Shanghai

WORKDIR /root/

# คัดลอกไบนารีจากขั้นตอนการ build
COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

# สร้างไดเรกทอรีสำหรับเก็บ logs
RUN mkdir -p logs

# เปิดพอร์ต
EXPOSE 8080

# ตรวจสอบสุขภาพ
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# รันแอปพลิเคชัน
CMD ["./main"]