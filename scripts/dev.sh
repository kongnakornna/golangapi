set -e

# สร้างเอกสาร Swagger ล่าสุดก่อน
./scripts/swagger.sh

# ตรวจสอบว่าได้ติดตั้ง air หรือไม่
if ! command -v air &> /dev/null; then
    echo "กำลังติดตั้ง air..."
    go install github.com/air-verse/air@latest
fi

# รับค่า GOPATH
GOPATH=$(go env GOPATH)
AIR_BIN="$GOPATH/bin/air"

# ตรวจสอบว่า air ติดตั้งสำเร็จหรือไม่
if [ ! -f "$AIR_BIN" ]; then
    echo "ข้อผิดพลาด: การติดตั้ง air ล้มเหลว"
    exit 1
fi

# ใช้พาธแบบเต็มเพื่อรัน air
echo "กำลังเริ่ม air..."
"$AIR_BIN" -c .air.toml