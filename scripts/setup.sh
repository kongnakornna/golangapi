# ตั้งค่าสี
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m' # ไม่มีสี

# ตรวจสอบว่าติดตั้ง Docker แล้วหรือไม่
if ! command -v docker &> /dev/null; then
    echo -e "${RED}ไม่ได้ติดตั้ง Docker กรุณาติดตั้ง Docker ก่อน${NC}"
    exit 1
fi

# ตรวจสอบว่าติดตั้ง Docker Compose แล้วหรือไม่
if ! command -v docker-compose &> /dev/null; then
    echo -e "${RED}ไม่ได้ติดตั้ง Docker Compose กรุณาติดตั้ง Docker Compose ก่อน${NC}"
    exit 1
fi

# กำลังสร้างสภาพแวดล้อมสำหรับการพัฒนา
echo -e "${GREEN}กำลังสร้างสภาพแวดล้อมสำหรับการพัฒนา...${NC}"

# กำลังสร้างไฟล์ docker-compose.yml
cat > docker-compose.yml << EOF
version: '3'

services:
  postgres:
    image: postgres:17-alpine
    container_name: go-rest-postgres
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: postgres
      POSTGRES_DB: restapi
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U postgres"]
      interval: 5s
      timeout: 5s
EOF