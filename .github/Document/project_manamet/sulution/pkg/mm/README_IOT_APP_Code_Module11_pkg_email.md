# Module 11: pkg/email (Email Sending)

## สำหรับโฟลเดอร์ `internal/pkg/email/`

ไฟล์ที่เกี่ยวข้อง:
- `internal/pkg/email/sender.go`
- `internal/pkg/email/gomail_sender.go`
- `internal/pkg/email/templates/verification.html`
- `internal/pkg/email/templates/reset_password.html`

---

## หลักการ (Concept)

### คืออะไร?
Email package คือส่วนที่รับผิดชอบการส่งอีเมลจากแอปพลิเคชันไปยังผู้ใช้ เช่น การยืนยันอีเมล (verification), การรีเซ็ตรหัสผ่าน (password reset), การแจ้งเตือน (alerts), และรายงานอัตโนมัติ (reports) โดยแยก logic การส่งออกจาก business logic

### มีกี่แบบ?
1. **SMTP sender** – ส่งผ่าน SMTP server (Gmail, Outlook, SendGrid)
2. **Template-based** – ใช้ HTML template สำหรับอีเมลสวยงาม
3. **Queue-based** – ส่งแบบ async ผ่าน worker (ใช้ร่วมกับ Redis queue)
4. **Attachment support** – แนบไฟล์ (PDF, Excel)
5. **Bulk email** – ส่งอีเมลจำนวนมาก (พร้อม rate limiting)

### ใช้อย่างไร / นำไปใช้กรณีไหน
- ใช้ส่ง verification link เมื่อผู้ใช้สมัครสมาชิก
- ใช้ส่ง reset password link
- ใช้ส่งรายงานประจำวัน/สัปดาห์
- ใช้ส่งการแจ้งเตือนเมื่อเซนเซอร์เกินค่า threshold

### ทำไมต้องใช้
- อีเมลเป็นช่องทางติดต่อหลักสำหรับการยืนยันตัวตน
- แยก template ออกจากโค้ด (maintainability)
- รองรับ HTML ทำให้อีเมลดูสวยงาม

### ประโยชน์ที่ได้รับ
- ผู้ใช้ได้รับข้อมูลสำคัญทางอีเมล
- ลดการใช้รหัส OTP ทาง SMS (ประหยัดต้นทุน)
- สร้างความน่าเชื่อถือให้ระบบ

### ข้อควรระวัง
- SMTP password ห้าม hardcode ใช้ environment variables
- ระวัง rate limit ของ SMTP provider (Gmail ~500/วัน)
- ต้องมี fallback mechanism (ถ้าส่งไม่สำเร็จ ให้ retry)

### ข้อดี
- ส่งได้หลาย provider, template ยืดหยุ่น

### ข้อเสีย
- ช้ากว่า SMS (ไม่เหมาะกับ OTP ที่ต้องการความเร็ว)
- อาจถูกตีว่าเป็น spam ถ้าไม่ตั้งค่า DKIM/SPF

### ข้อห้าม
- ห้ามส่งอีเมลหาคนที่ไม่ยินยอม (GDPR violation)
- ห้ามส่ง sensitive data ใน plain text (ใช้ HTTPS link แทน)

---

## โค้ดที่รันได้จริง

### 1. Interface Sender – `sender.go`

```go
// Package email provides email sending capabilities with templates.
// ----------------------------------------------------------------
// แพ็คเกจ email ให้บริการส่งอีเมลพร้อมเทมเพลต
package email

// Sender defines the interface for email sending.
// ----------------------------------------------------------------
// Sender กำหนด interface สำหรับการส่งอีเมล
type Sender interface {
	// Send sends an email to a single recipient.
	// ----------------------------------------------------------------
	// Send ส่งอีเมลไปยังผู้รับเดียว
	Send(to, subject, body string) error
	
	// SendWithAttachment sends email with file attachment.
	// ----------------------------------------------------------------
	// SendWithAttachment ส่งอีเมลพร้อมไฟล์แนบ
	SendWithAttachment(to, subject, body, attachmentPath string) error
}

// TemplateEngine defines interface for rendering email templates.
// ----------------------------------------------------------------
// TemplateEngine กำหนด interface สำหรับการ render เทมเพลตอีเมล
type TemplateEngine interface {
	// Render renders a template with given data.
	// ----------------------------------------------------------------
	// Render render เทมเพลตด้วยข้อมูลที่กำหนด
	Render(templateName string, data interface{}) (string, error)
}
```

### 2. Gomail Sender Implementation – `gomail_sender.go`

```go
package email

import (
	"crypto/tls"
	"fmt"

	"gopkg.in/gomail.v2"
)

// GomailSender implements Sender using gomail library.
// ----------------------------------------------------------------
// GomailSender อิมพลีเมนต์ Sender ด้วยไลบรารี gomail
type GomailSender struct {
	host     string
	port     int
	username string
	password string
	from     string
}

// NewGomailSender creates a new Gomail sender.
// ----------------------------------------------------------------
// NewGomailSender สร้าง Gomail sender ใหม่
func NewGomailSender(host string, port int, username, password, from string) *GomailSender {
	return &GomailSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
	}
}

// Send sends an email with HTML body.
// ----------------------------------------------------------------
// Send ส่งอีเมลแบบ HTML
func (s *GomailSender) Send(to, subject, body string) error {
	return s.SendWithAttachment(to, subject, body, "")
}

// SendWithAttachment sends email with optional file attachment.
// ----------------------------------------------------------------
// SendWithAttachment ส่งอีเมลพร้อมไฟล์แนบ (ถ้ามี)
func (s *GomailSender) SendWithAttachment(to, subject, body, attachmentPath string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Attach file if path provided
	// แนบไฟล์ถ้ามีการระบุ path
	if attachmentPath != "" {
		m.Attach(attachmentPath)
	}

	dialer := gomail.NewDialer(s.host, s.port, s.username, s.password)
	dialer.TLSConfig = &tls.Config{InsecureSkipVerify: false}
	
	if err := dialer.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}
	return nil
}
```

### 3. HTML Templates – `templates/verification.html`

```html
<!-- 
  Verification Email Template
  เทมเพลตอีเมลยืนยันตัวตน
-->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .header { background: #4CAF50; color: white; padding: 10px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { padding: 20px; }
        .button { display: inline-block; background: #4CAF50; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
        .footer { font-size: 12px; color: #777; text-align: center; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>ยืนยันอีเมลของคุณ | Verify Your Email</h2>
        </div>
        <div class="content">
            <p>สวัสดีคุณ {{.Name}},</p>
            <p>กรุณาคลิกที่ปุ่มด้านล่างเพื่อยืนยันอีเมลของคุณ:</p>
            <p style="text-align: center;">
                <a href="{{.Link}}" class="button">ยืนยันอีเมล | Verify Email</a>
            </p>
            <p>หากคุณไม่ได้สมัครสมาชิก กรุณาละเว้นอีเมลนี้</p>
            <p>ลิงก์นี้จะหมดอายุใน 24 ชั่วโมง</p>
            <hr>
            <p>Hello {{.Name}},</p>
            <p>Please click the button below to verify your email:</p>
            <p style="text-align: center;">
                <a href="{{.Link}}" class="button">Verify Email</a>
            </p>
            <p>If you did not sign up, please ignore this email.</p>
            <p>This link expires in 24 hours.</p>
        </div>
        <div class="footer">
            &copy; 2025 CMON System. All rights reserved.
        </div>
    </div>
</body>
</html>
```

### 4. Password Reset Template – `templates/reset_password.html`

```html
<!-- 
  Password Reset Email Template
  เทมเพลตอีเมลรีเซ็ตรหัสผ่าน
-->
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <style>
        body { font-family: Arial, sans-serif; line-height: 1.6; color: #333; }
        .container { max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #ddd; border-radius: 5px; }
        .header { background: #f44336; color: white; padding: 10px; text-align: center; border-radius: 5px 5px 0 0; }
        .content { padding: 20px; }
        .button { display: inline-block; background: #f44336; color: white; padding: 10px 20px; text-decoration: none; border-radius: 5px; }
        .footer { font-size: 12px; color: #777; text-align: center; margin-top: 20px; }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h2>รีเซ็ตรหัสผ่านของคุณ | Reset Your Password</h2>
        </div>
        <div class="content">
            <p>สวัสดีคุณ {{.Name}},</p>
            <p>เราได้รับคำขอรีเซ็ตรหัสผ่านของคุณ คลิกปุ่มด้านล่างเพื่อตั้งรหัสผ่านใหม่:</p>
            <p style="text-align: center;">
                <a href="{{.Link}}" class="button">รีเซ็ตรหัสผ่าน | Reset Password</a>
            </p>
            <p>หากคุณไม่ได้ขอรีเซ็ตรหัสผ่าน กรุณาละเว้นอีเมลนี้</p>
            <p>ลิงก์นี้จะหมดอายุใน 1 ชั่วโมง</p>
            <hr>
            <p>Hello {{.Name}},</p>
            <p>We received a request to reset your password. Click the button below to set a new password:</p>
            <p style="text-align: center;">
                <a href="{{.Link}}" class="button">Reset Password</a>
            </p>
            <p>If you did not request a password reset, please ignore this email.</p>
            <p>This link expires in 1 hour.</p>
        </div>
        <div class="footer">
            &copy; 2025 CMON System. All rights reserved.
        </div>
    </div>
</body>
</html>
```

### 5. Template Renderer (Go HTML template) – `template_renderer.go`

```go
package email

import (
	"bytes"
	"html/template"
	"path/filepath"
)

// GoTemplateEngine implements TemplateEngine using Go's html/template.
// ----------------------------------------------------------------
// GoTemplateEngine อิมพลีเมนต์ TemplateEngine ด้วย html/template ของ Go
type GoTemplateEngine struct {
	templatesDir string
}

// NewGoTemplateEngine creates a new template engine.
// ----------------------------------------------------------------
// NewGoTemplateEngine สร้าง template engine ใหม่
func NewGoTemplateEngine(templatesDir string) *GoTemplateEngine {
	return &GoTemplateEngine{templatesDir: templatesDir}
}

// Render renders a template file with given data.
// ----------------------------------------------------------------
// Render render ไฟล์เทมเพลตด้วยข้อมูลที่กำหนด
func (e *GoTemplateEngine) Render(templateName string, data interface{}) (string, error) {
	tmplPath := filepath.Join(e.templatesDir, templateName+".html")
	tmpl, err := template.ParseFiles(tmplPath)
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return "", err
	}
	return buf.String(), nil
}
```

### 6. ตัวอย่างการใช้งานใน handler (auth_handler.go)

```go
// ส่งอีเมลยืนยันหลังจาก register
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
    // ... create user
    // ส่งอีเมลยืนยันแบบ async ผ่าน worker
    verificationLink := fmt.Sprintf("https://yourapp.com/verify?token=%s", token)
    data := map[string]string{
        "Name": user.FullName,
        "Link": verificationLink,
    }
    body, _ := h.templateEngine.Render("verification", data)
    
    // Enqueue to email worker (using Redis queue)
    task := &queue.Task{
        Type: "email",
        Payload: map[string]interface{}{
            "to":      user.Email,
            "subject": "Verify your email",
            "body":    body,
        },
    }
    h.queue.Enqueue(ctx, task)
}
```

---

## วิธีใช้งาน module นี้

1. ติดตั้ง dependency:
   ```bash
   go get gopkg.in/gomail.v2
   ```
2. วางไฟล์ตามโครงสร้าง
3. สร้าง sender instance ใน `main.go`:
   ```go
   emailSender := email.NewGomailSender(
       os.Getenv("SMTP_HOST"),
       587,
       os.Getenv("SMTP_USER"),
       os.Getenv("SMTP_PASSWORD"),
       os.Getenv("SMTP_FROM"),
   )
   templateEngine := email.NewGoTemplateEngine("./internal/pkg/email/templates")
   ```
4. ใช้ร่วมกับ worker (Module 6) เพื่อส่งแบบ async

---

## ตารางสรุป Components

| Component | หน้าที่ | ตัวอย่าง |
|-----------|--------|----------|
| `Sender` | ส่งอีเมลผ่าน SMTP | `GomailSender` |
| `TemplateEngine` | Render HTML template | `GoTemplateEngine` |
| Templates | ไฟล์ HTML สำหรับอีเมล | `verification.html` |

---

## แบบฝึกหัดท้าย module (3 ข้อ)

1. เพิ่ม method `SendBatch` ใน `Sender` ที่รับ slice ของ recipients และส่งอีเมลเดียวกันไปหลายคน (BCC)
2. สร้าง template `alert.html` สำหรับการแจ้งเตือนเซนเซอร์เกิน threshold พร้อมตารางข้อมูล
3. Implement rate limiter ใน `GomailSender` ที่จำกัดจำนวนอีเมลต่อชั่วโมง (ป้องกันถูก SMTP provider block)

---

## แหล่งอ้างอิง

- [Gomail documentation](https://github.com/go-gomail/gomail)
- [Go html/template](https://pkg.go.dev/html/template)
- [SMTP best practices](https://mailgun.com/blog/email-best-practices/)

---

**หมายเหตุ:** module นี้ครบถ้วนสำหรับ `pkg/email` หากต้องการ module เพิ่มเติม (เช่น `pkg/jwt`, `pkg/hash`) โปรดแจ้ง