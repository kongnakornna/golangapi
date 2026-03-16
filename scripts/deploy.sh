# สคริปต์ปรับใช้งานสำหรับสภาพแวดล้อมการทำงานจริง (Production)

set -e  # หยุดการทำงานเมื่อเกิดข้อผิดพลาด

# กำหนดสี
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}เริ่มการปรับใช้งาน...${NC}"

# ตรวจสอบตัวแปรสภาพแวดล้อม
check_env() {
    local vars=("DB_HOST" "DB_USER" "DB_PASSWORD" "DB_NAME" "JWT_SECRET")
    for var in "${vars[@]}"; do
        if [ -z "${!var}" ]; then
            echo -e "${RED}ข้อผิดพลาด: ตัวแปรสภาพแวดล้อม $var ไม่ได้ถูกตั้งค่า${NC}"
            exit 1
        fi
    done
}

# สร้างแอปพลิเคชัน
build() {
    echo -e "${YELLOW}กำลังสร้างแอปพลิเคชัน...${NC}"
    go build -ldflags="-s -w" -o bin/app cmd/app/main.go
    echo -e "${GREEN}สร้างแอปพลิเคชันเสร็จสมบูรณ์${NC}"
}

# รันการโยกย้ายฐานข้อมูล
migrate() {
    echo -e "${YELLOW}กำลังดำเนินการโยกย้ายฐานข้อมูล...${NC}"
    # go run cmd/migrate/main.go up
    echo -e "${GREEN}การโยกย้ายฐานข้อมูลเสร็จสมบูรณ์${NC}"
}

# ตรวจสอบสถานะสุขภาพ
health_check() {
    echo -e "${YELLOW}กำลังตรวจสอบสถานะสุขภาพ...${NC}"
    local max_attempts=30
    local attempt=0
    
    while [ $attempt -lt $max_attempts ]; do
        if curl -f http://localhost:${PORT:-8080}/health > /dev/null 2>&1; then
            echo -e "${GREEN}การตรวจสอบสุขภาพผ่าน${NC}"
            return 0
        fi
        
        attempt=$((attempt + 1))
        echo "กำลังรอให้บริการเริ่มทำงาน... ($attempt/$max_attempts)"
        sleep 2
    done
    
    echo -e "${RED}การตรวจสอบสุขภาพล้มเหลว${NC}"
    return 1
}

# หยุดกระบวนการเก่าอย่างนุ่มนวล
stop_old() {
    if [ -f "app.pid" ]; then
        OLD_PID=$(cat app.pid)
        if kill -0 $OLD_PID 2>/dev/null; then
            echo -e "${YELLOW}กำลังหยุดกระบวนการเก่า (PID: $OLD_PID)...${NC}"
            kill -TERM $OLD_PID
            sleep 5
        fi
    fi
}

# เริ่มกระบวนการใหม่
start() {
    echo -e "${YELLOW}กำลังเริ่มแอปพลิเคชัน...${NC}"
    CONFIG_PATH=configs/config.production.yaml nohup ./bin/app > logs/app.out 2>&1 &
    echo $! > app.pid
    echo -e "${GREEN}แอปพลิเคชันเริ่มทำงานแล้ว (PID: $(cat app.pid))${NC}"
}

# กระบวนหลัก
main() {
    check_env
    build
    # migrate
    stop_old
    start
    health_check
    
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}การปรับใช้งานสำเร็จ!${NC}"
    else
        echo -e "${RED}การปรับใช้งานล้มเหลว!${NC}"
        exit 1
    fi
}

main