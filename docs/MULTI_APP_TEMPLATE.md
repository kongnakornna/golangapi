# เทมเพลตสำหรับหลายแอปพลิเคชัน

เทมเพลตนี้ใช้สำหรับเพิ่มแอปพลิเคชันที่สอง (เช่น `app2`) ภายใน repository เดียวกัน โดยใช้โมดูลโดเมน `internal/core` และโครงสร้างพื้นฐาน `internal/platform` ร่วมกัน


## แม่แบบโครงสร้างไดเร็กทอรี

```text
cmd/
app2/
main.go
internal/
apps/
app2/
bootstrap/
app.go
injection/
dependencies.go
repositories.go
services.go
handlers.go
router/
router.go
v1/
public_routes.go
protected_routes.go
swagger.go

```

## ขั้นตอนการสร้างแอปพลิเคชันใหม่

1. คัดลอกไดเร็กทอรี assembly ของแอปพลิเคชันที่มีอยู่เป็นแม่แบบ:

- `internal/apps/app` → `internal/apps/app2`

2. เพิ่มไฟล์ entry: `cmd/app2/main.go`

3. ปรับเส้นทางการกำหนดค่าตามต้องการ (แนะนำให้ใช้ตัวแปรสภาพแวดล้อมอิสระ)

4. ปรับแต่ง routes และ module injection ตามต้องการ (เก็บเฉพาะโมดูลที่แอปพลิเคชันต้องการ)

## เทมเพลตเริ่มต้น (`cmd/app2/main.go`)
```go
package main

import (
"context"
"os"
"os/signal"
"syscall"
"time"

app2bootstrap "github.com/kongnakornna/golangapi/internal/apps/app2/bootstrap"
"github.com/kongnakornna/golangapi/pkg/logger"
)

func main() {
var appLogger logger.Logger = logger.Default()

application, err := app2bootstrap.New()

if err != nil {
appLogger.Error("Failed to create application", "error", err)
os.Exit(1)

}
appLogger = application.Logger()

serverErrCh := application.StartServer()

signalCh := สร้าง (chan os.Signal, 1)

signal.Notify(signalCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

เลือก {
กรณี err := <-serverErrCh:

appLogger.Error("ข้อผิดพลาดของเซิร์ฟเวอร์", "error", err)

กรณี sig := <-signalCh:

appLogger.Info("ได้รับสัญญาณระบบ เริ่มการปิดระบบอย่างนุ่มนวล", "signal", sig.String())

}

ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

defer cancel()

if err := application.Shutdown(ctx); err != nil {

appLogger.Error("Application shutdown failed", "error", err)

os.Exit(1)

}
}
```

## เทมเพลตแอสเซมบลี (`internal/apps/app2/bootstrap/app.go`)

> คุณสามารถคัดลอก `internal/apps/app/bootstrap/app.go` โดยตรงและปรับเส้นทางการกำหนดค่าตามต้องการได้:

- การกำหนดค่าเริ่มต้น: อ่าน `CONFIG_PATH` (ใช้ร่วมกับแอปพลิเคชันที่มีอยู่แล้ว)

- การกำหนดค่าอิสระ: ใช้ `APP2_CONFIG_PATH` (แนะนำ)

ตัวอย่างที่แนะนำ:

```go

func getConfigPath() string {

configPath := os.Getenv("APP2_CONFIG_PATH")

if configPath == "" {

configPath = "configs/config.yaml"

}

return configPath

}
```

## เทมเพลตเส้นทาง (`internal/apps/app2/router/router.go`)

- นำโมดูล `internal/core` กลับมาใช้ใหม่

- ประกอบเฉพาะโมดูลที่จำเป็นเท่านั้น

ไฮไลท์ตัวอย่าง:

```go

api.Setup(router, api.RouterConfig{

UserHandler: deps.Handlers.UserHandler,

AuthHandler: deps.Handlers.AuthHandler,

HealthHandler: deps.Handlers.HealthHandler,

JWTSecret: deps.Config.JWT.Secret,

Logger: appLogger,

})
```

## การตัดแต่งโมดูลและการพึ่งพา

- โมดูลที่ไม่จำเป็นสามารถลบออกจากการฉีดและการกำหนดเส้นทางได้:

- `internal/apps/app2/bootstrap/injection/*`

- `internal/apps/app2/router/*`

- โมดูลโดเมนจะอยู่ใน `internal/core` เสมอและสามารถนำกลับมาใช้ใหม่ได้โดยหลายแอปพลิเคชัน

## ข้อแนะนำ

- รักษาโครงสร้างที่สอดคล้องกันสำหรับจุดเริ่มต้นของแอปพลิเคชันและไฟล์การกำหนดค่าเพื่ออำนวยความสะดวกในการทำงานร่วมกันและการขยายทีม

- แอปพลิเคชันต่างๆ สามารถใช้พอร์ตและไฟล์การกำหนดค่าที่แตกต่างกันเพื่อหลีกเลี่ยงความขัดแย้ง