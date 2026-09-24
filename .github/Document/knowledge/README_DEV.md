## การทำ Live-reload ใน Go (เมื่อ code เปลี่ยนแปลง)

Live-reload ช่วยให้เซิร์ฟเวอร์รีสตาร์ทอัตโนมัติทุกครั้งที่คุณบันทึกไฟล์ ไม่ต้องกด Ctrl+C แล้วรันใหม่เอง

สำหรับโปรเจกต์ Go ที่ใช้ `chi`, `viper`, `asynq` ฯลฯ มีเครื่องมือยอดนิยมดังนี้

### 1. ใช้ **Air** (แนะนำ)

Air เป็นเครื่องมือที่ใช้กันแพร่หลาย รองรับ Go module และทำงานข้าม OS

#### ติดตั้ง Air

```bash
go install github.com/air-verse/air@latest
```

หรือใช้ `go get` (แล้วแต่เวอร์ชัน Go)

#### สร้างไฟล์คอนฟิก `.air.toml` ที่ root โปรเจกต์

```bash
air init
```

#### ปรับแต่ง `.air.toml` ให้เหมาะกับโครงสร้างโปรเจกต์คุณ

ตัวอย่างคร่าวๆ (แก้ให้ตรงกับโปรเจกต์ `icmongolang`):

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  args_bin = []
  bin = "./tmp/main"
  cmd = "go build -o ./tmp/main ./cmd/api"   # <-- เปลี่ยน path ให้ตรงกับ main ของคุณ
  delay = 1000
  exclude_dir = ["assets", "tmp", "vendor", "testdata", "docs", "migrations"]
  exclude_file = []
  exclude_regex = ["_test.go"]
  exclude_unchanged = false
  follow_symlink = false
  full_bin = ""
  include_dir = []
  include_ext = ["go", "tpl", "tmpl", "html"]
  include_file = []
  kill_delay = "0s"
  log = "build-errors.log"
  send_interrupt = false
  stop_on_error = true

[color]
  app = ""
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = false

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = true
  keep_scroll = true
```

สำคัญ: แก้ `cmd` ให้ชี้ไปที่ `main.go` จริงของโปรเจกต์ เช่น ถ้าคุณมี `cmd/server/main.go` ก็ตั้งเป็น `"go build -o ./tmp/main ./cmd/server"`

#### รันด้วย Air

```bash
air
```

แค่นี้ Air จะ watch ไฟล์ `.go`, `.tpl`, `.html` ฯลฯ และ rebuild + restart ทุกครั้งที่มีการเปลี่ยนแปลง

---

### 2. ใช้ **fresh** (เก่า แต่ใช้งานง่าย)

ติดตั้ง:

```bash
go install github.com/pilu/fresh@latest
```

รัน:

```bash
fresh
```

ไม่ต้องตั้งค่าอะไรเพิ่ม (แต่จะไม่ละเอียดเท่า Air)

---

### 3. ใช้ **CompileDaemon** (ถ้าชอบ approach แบบ minimal)

```bash
go install github.com/githubnemo/CompileDaemon@latest
```

รัน:

```bash
CompileDaemon -command="./your-binary"
```

หรือรวม build และ run:

```bash
CompileDaemon -build="go build -o ./tmp/app ./cmd/api" -command="./tmp/app"
```

---

### 4. ใช้ `nodemon` หรือ `reflex` (ผ่าน tools ภายนอก)

ถ้าติดตั้ง Node.js ไว้แล้ว:

```bash
nodemon --exec go run ./cmd/api --signal SIGTERM
```

หรือใช้ `reflex` (Linux/macOS):

```bash
reflex -r '\.go$' -- sh -c 'go run ./cmd/api'
```

---

### 5. สำหรับโปรเจกต์นี้ (icmongolang) – ตัวอย่างที่ใช้ Air

สมมติว่า `main.go` อยู่ที่ `cmd/api/main.go`

1. สร้าง `.air.toml` (หรือใช้ default ก็ได้)
2. แก้ `cmd` เป็น `"go build -o ./tmp/main ./cmd/api"`
3. รัน `air` ใน root โปรเจกต์

```bash
cd D:\git\icmongolang
air
```

ทุกครั้งที่คุณแก้ไขไฟล์ `.go` Air จะแสดง log rebuild และ restart server อัตโนมัติ

---

### ข้อควรระวัง

- ถ้าโปรเจกต์ใช้ `os.Exit()` หรือ `log.Fatal()` ในโค้ด อาจทำให้ Air หยุดทำงาน ควรเปลี่ยนเป็น return error แล้วให้ main จัดการ
- ถ้าใช้ `asynq` worker, ควรแยก binary หรือใช้ Air แยก process อีกตัว
- กรณีใช้ vendor (`go mod vendor`) ให้ Air ทำงานได้ปกติ ไม่ต้องตั้งค่าพิเศษ

---

### สรุป

**แนะนำ Air** เพราะ:
- ใช้งานง่าย
- รองรับ Windows / macOS / Linux
- ตั้งค่าได้ละเอียด
- มี community ใช้เยอะ
ถ้าต้องการ minimal จริงๆ ใช้ `CompileDaemon` หรือ `reflex` ก็ได้
ต้องการให้ช่วยตั้งค่า `.air.toml` ให้ตรงกับโครงสร้างโปรเจกต์ของคุณไหม? แจ้ง path ของ `main.go` และพอร์ตที่ใช้ มาได้เลย 