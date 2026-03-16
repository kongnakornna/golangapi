set -e

# การกำหนดค่าเริ่มต้น
IMAGE_NAME="go-rest-starter"
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "latest")
DOCKERFILE="deploy/docker/Dockerfile"

# วิเคราะห์พารามิเตอร์จากบรรทัดคำสั่ง
while [[ $# -gt 0 ]]; do
  case $1 in
    -n|--name)
      IMAGE_NAME="$2"
      shift 2
      ;;
    -v|--version)
      VERSION="$2"
      shift 2
      ;;
    -f|--file)
      DOCKERFILE="$2"
      shift 2
      ;;
    -h|--help)
      echo "การใช้งาน: $0 [ตัวเลือก]"
      echo "ตัวเลือก:"
      echo "  -n, --name IMAGE_NAME    ชื่อ Docker image (ค่าเริ่มต้น: go-rest-starter)"
      echo "  -v, --version VERSION    แท็กเวอร์ชันของ image (ค่าเริ่มต้น: git describe)"
      echo "  -f, --file DOCKERFILE    ตำแหน่งของ Dockerfile (ค่าเริ่มต้น: deploy/docker/Dockerfile)"
      echo "  -h, --help               แสดงข้อความช่วยเหลือนี้"
      exit 0
      ;;
    *)
      echo "ไม่รู้จักตัวเลือก $1"
      exit 1
      ;;
  esac
done

echo "กำลังสร้าง Docker image..."
echo "Image: ${IMAGE_NAME}:${VERSION}"
echo "Dockerfile: ${DOCKERFILE}"

# สร้างเอกสาร Swagger ล่าสุด
echo "กำลังสร้างเอกสาร Swagger..."
./scripts/swagger.sh

# สร้าง Docker image
echo "กำลังสร้าง Docker image..."
docker build -t ${IMAGE_NAME}:${VERSION} -f ${DOCKERFILE} .

# ตั้งแท็ก latest พร้อมกัน
if [[ "${VERSION}" != "latest" ]]; then
    docker tag ${IMAGE_NAME}:${VERSION} ${IMAGE_NAME}:latest
    echo "ตั้งแท็กเป็น ${IMAGE_NAME}:latest แล้ว"
fi

echo "การสร้าง Docker เสร็จสมบูรณ์!"
echo "Images ที่สร้างแล้ว:"
docker images ${IMAGE_NAME}

echo ""
echo "วิธีเรียกใช้ container:"
echo "docker run -d -p 7001:7001 --name go-rest-starter-container \\"
echo "  -e APP_DATABASE_HOST=host.docker.internal \\"
echo "  -e APP_REDIS_HOST=host.docker.internal \\"
echo "  ${IMAGE_NAME}:${VERSION}"