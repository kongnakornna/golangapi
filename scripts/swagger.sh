# ตรวจสอบว่ามีคำสั่ง swag หรือไม่
if ! command -v swag &> /dev/null; then
    echo "ไม่พบคำสั่ง swag กำลังติดตั้ง..."
    go install github.com/swaggo/swag/cmd/swag@latest
    
    # ตรวจสอบว่าการติดตั้งสำเร็จหรือไม่
    if ! command -v swag &> /dev/null; then
        echo "การติดตั้ง swag ล้มเหลว อาจเป็นปัญหาเกี่ยวกับตัวแปรสภาพแวดล้อม PATH"
        
        # พยายามใช้ swag ใน GOPATH โดยตรง
        GOPATH=$(go env GOPATH)
        SWAG_PATH="$GOPATH/bin/swag"
        
        if [ -f "$SWAG_PATH" ]; then
            echo "พบ swag ที่ $SWAG_PATH จะใช้เส้นทางนี้โดยตรง"
            
            # ล้างไดเรกทอรี swag-docs ที่มีอยู่
            rm -rf api/app
            
            # เรียกใช้ swag ด้วยเส้นทางแบบเต็ม
            "$SWAG_PATH" init -g cmd/app/main.go -o api/app
            
            echo "เอกสาร Swagger ถูกสร้างไปยังไดเรกทอรี api/app/ แล้ว"
            exit 0
        else
            echo "ไม่พบ swag ใน $GOPATH/bin โปรดตรวจสอบว่าการติดตั้งสำเร็จและเพิ่ม $GOPATH/bin ลงในตัวแปรสภาพแวดล้อม PATH"
            echo "สามารถลอง: export PATH=\$PATH:\$(go env GOPATH)/bin"
            exit 1
        fi
    fi
    echo "ติดตั้ง swag สำเร็จ"
fi

# ล้างไดเรกทอรี swag-docs ที่มีอยู่
rm -rf api/app

# สร้างเอกสาร Swagger
swag init -g cmd/app/main.go -o api/app

echo "เอกสาร Swagger ถูกสร้างไปยังไดเรกทอรี api/app/ แล้ว"