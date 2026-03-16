set -e

# รับข้อมูลเวอร์ชัน
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date +%Y-%m-%d_%H:%M:%S)
COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# พารามิเตอร์การ build
LDFLAGS="-s -w"
OUTPUT_DIR="build"

# สร้างไดเรกทอรีเอาต์พุต
mkdir -p ${OUTPUT_DIR}

echo "กำลัง build Go REST Starter..."
echo "เวอร์ชัน: ${VERSION}"
echo "เวลาที่ build: ${BUILD_TIME}"
echo "Commit: ${COMMIT}"

# สร้างเอกสาร Swagger ล่าสุด
./scripts/swagger.sh

# build ไบนารีสำหรับแพลตฟอร์มต่างๆ
echo "กำลัง build สำหรับ Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ${OUTPUT_DIR}/app-linux-amd64 cmd/app/main.go

echo "กำลัง build สำหรับ macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ${OUTPUT_DIR}/app-darwin-amd64 cmd/app/main.go

echo "กำลัง build สำหรับ macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="${LDFLAGS}" -o ${OUTPUT_DIR}/app-darwin-arm64 cmd/app/main.go

echo "กำลัง build สำหรับ Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="${LDFLAGS}" -o ${OUTPUT_DIR}/app-windows-amd64.exe cmd/app/main.go

# build สำหรับแพลตฟอร์มปัจจุบัน
echo "กำลัง build สำหรับแพลตฟอร์มปัจจุบัน..."
go build -ldflags="${LDFLAGS}" -o ${OUTPUT_DIR}/app cmd/app/main.go

echo "การ build เสร็จสมบูรณ์! ไฟล์อยู่ใน ${OUTPUT_DIR}/"
ls -la ${OUTPUT_DIR}/
