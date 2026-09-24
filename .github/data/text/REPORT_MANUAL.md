# 📄 คู่มือระบบรายงาน PDF — Report System Manual

## สารบัญ
1. [ภาพรวมระบบ](#1-ภาพรวมระบบ)
2. [รายการเทมเพลตทั้งหมด](#2-รายการเทมเพลตทั้งหมด)
3. [โครงสร้างแพ็กเกจ](#3-โครงสร้างแพ็กเกจ)
4. [การเรียกใช้งาน API](#4-การเรียกใช้งาน-api)
5. [การเพิ่มเทมเพลตใหม่](#5-การเพิ่มเทมเพลตใหม่)
6. [การปรับแต่ง](#6-การปรับแต่ง)
7. [Troubleshooting](#7-troubleshooting)

---

## 1. ภาพรวมระบบ

ระบบรายงานใช้ **chromedp** (Headless Chrome) แปลง HTML → PDF

```
Request → Handler → UseCase → report.GeneratePDF(template, data) → HTML Template → Chrome → PDF bytes → Response
```

### ข้อกำหนด
- ต้องติดตั้ง **Google Chrome / Chromium** บนเซิร์ฟเวอร์
- chromedp จะหา Chrome จาก PATH โดยอัตโนมัติ
- รองรับ Windows / Linux / macOS

---

## 2. รายการเทมเพลตทั้งหมด

| # | ชื่อไฟล์ | ประเภทเอกสาร | Go DTO |
|---|---------|-------------|--------|
| 1 | `quotation.html` | ใบเสนอราคา | `report.QuotationData` |
| 2 | `invoice.html` | ใบแจ้งหนี้ (มี QR Code) | `report.InvoiceData` |
| 3 | `purchase_order.html` | ใบสั่งซื้อ | `report.PurchaseOrderData` |
| 4 | `part_picking.html` | เอกสารเบิกอะไหล่ | `report.PartPickingData` |
| 5 | `receipt.html` | ใบเสร็จรับเงิน (80mm) | `report.ReceiptData` |
| 6 | `credit_note.html` | ใบลดหนี้ | `report.CreditNoteData` |
| 7 | `debit_note.html` | ใบเพิ่มหนี้ | `report.DebitNoteData` |
| 8 | `delivery_sheet.html` | ใบส่งของ | `report.DeliverySheetData` |
| 9 | `job_card.html` | ใบงานซ่อม (ฉบับเต็ม) | `report.JobCardData` |
| 10 | `daily_sales.html` | รายงานยอดขายรายวัน | `report.DailySalesData` |
| 11 | `inventory_summary.html` | รายงานสินค้าคงคลัง | `report.InventorySummaryData` |
| 12 | `customer_list.html` | รายชื่อลูกค้า | `report.CustomerListData` |

---

## 3. โครงสร้างแพ็กเกจ

```
pkg/report/
├── generator.go       # ฟังก์ชัน GenerateHTML() / GeneratePDF()
├── models.go          # DTO structs ทั้ง 12 รายงาน
├── bahtthai.go        # แปลงตัวเลข → ตัวอักษรไทย (บาท)
└── templates/
    ├── base.html      # CSS/Layout ส่วนกลาง (ฟอนต์ Sarabun, A4)
    ├── quotation.html
    ├── invoice.html
    ├── purchase_order.html
    ├── part_picking.html
    ├── receipt.html
    ├── credit_note.html
    ├── debit_note.html
    ├── delivery_sheet.html
    ├── job_card.html
    ├── daily_sales.html
    ├── inventory_summary.html
    └── customer_list.html
```

---

## 4. การเรียกใช้งาน API

### 4.1 สร้าง PDF
```go
import "icmongolang/pkg/report"

data := report.PurchaseOrderData{
    Company: report.CompanyInfo{
        Name:    "ร้านอู่ซ่อมรถ",
        Address: "123 ถนน...",
        Phone:   "02-xxx-xxxx",
        TaxID:   "เลขประจำตัวผู้เสียภาษี",
    },
    PONo:       "PO-2026-0001",
    Date:       time.Now(),
    Supplier:   "บริษัท อะไหล่ไทย จำกัด",
    Subtotal:   1000.00,
    TaxAmount:  70.00,
    GrandTotal: 1070.00,
    AmountWords: report.BahtThai(1070.00),
}

pdfBytes, err := report.GeneratePDF(ctx, report.TplPurchaseOrder, data)
if err != nil {
    // handle error
}
w.Header().Set("Content-Type", "application/pdf")
w.Write(pdfBytes)
```

### 4.2 สร้าง HTML (Preview)
```go
html, err := report.GenerateHTML(report.TplQuotation, data)
```

---

## 5. การเพิ่มเทมเพลตใหม่

### 5.1 สร้างไฟล์ HTML
สร้างไฟล์ `pkg/report/templates/new_report.html`:

```html
{{define "content"}}
<div class="header">
  <div class="company-info">
    <h1>{{.Company.Name}}</h1>
  </div>
  <div class="doc-title">
    <h2>รายงานใหม่</h2>
  </div>
</div>
<table>...</table>
{{end}}
```

### 5.2 เพิ่ม Constant ใน generator.go
```go
const TplNewReport TemplateType = "new_report"
```

### 5.3 ลงทะเบียน Template ใน init()
```go
tplNewReport = template.Must(template.New("new_report").Funcs(funcMap).ParseFiles(append(layouts, basePath+"new_report.html")...))
templates["new_report"] = tplNewReport
```

### 5.4 เพิ่ม Case ใน getTemplate()
```go
case TplNewReport:
    return tplNewReport
```

ใช้ `report.GeneratePDF(ctx, report.TplNewReport, data)` ได้ทันที

---

## 6. การปรับแต่ง

### 6.1 ฟอนต์
- ใช้ Google Fonts `Sarabun` (รองรับภาษาไทย)
- แก้ไขใน `base.html` ที่ `@import url(...)`
- หรือเปลี่ยนเป็นฟอนต์อื่นโดยแก้ CSS

### 6.2 ขนาดกระดาษ
แก้ไขใน `generator.go` ที่ `page.PrintToPDF()`:
```go
WithPaperWidth(210.0).   // A4 width (mm)
WithPaperHeight(297.0).  // A4 height (mm)
WithMarginTop(15.0).     // ขอบบน
WithMarginBottom(15.0).  // ขอบล่าง
WithMarginLeft(15.0).    // ขอบซ้าย
WithMarginRight(15.0).   // ขอบขวา
```

สำหรับสลิป 80mm: `WithPaperWidth(80.0).WithPaperHeight(297.0)`

### 6.3 โลโก้บริษัท
แก้ไขใน `base.html` หรือเพิ่ม `<img>` ในแต่ละ template

### 6.4 ฟังก์ชัน Template (FuncMap)
| ฟังก์ชัน | รายละเอียด |
|----------|-----------|
| `formatMoney` | แปลงตัวเลขเป็น "1,234.50" |
| `formatDate` | แปลง time.Time เป็น "02/01/2006" |
| `formatDateTime` | แปลงเป็น "02/01/2006 15:04" |
| `bahtthai` | แปลง float64 เป็นตัวอักษรไทย |
| `add` | บวกเลข |

---

## 7. Troubleshooting

| ปัญหา | สาเหตุ | วิธีแก้ไข |
|-------|-------|----------|
| PDF ไม่ออก / คืน error | ไม่ได้ติดตั้ง Chrome | ติดตั้ง Chrome/Chromium และเพิ่มใน PATH |
| ภาษาไทยแสดงเป็นสี่เหลี่ยม | ฟอนต์ไม่รองรับ | ใช้ฟอนต์ Sarabun หรือ TH Sarabun New |
| PDF หน้าเว้นว่าง | HTML template error | เรียก `GenerateHTML` ก่อน debug |
| chromedp timeout | Chrome ไม่ตอบสนอง | เพิ่ม time context, ตรวจสอบว่า Chrome version สูงกว่า 80 |
| `page.PrintToPDF` error | API เปลี่ยน | ตรวจสอบ cdproto version ใน go.sum |

### วิธีตรวจสอบ Chrome
```bash
# Windows
where chrome
# หรือ
"C:\Program Files\Google\Chrome\Application\chrome.exe" --version

# Linux
which google-chrome
google-chrome --version
```

### Debug HTML ก่อนสร้าง PDF
```go
html, err := report.GenerateHTML(report.TplQuotation, data)
if err != nil {
    log.Fatal(err)
}
fmt.Println(html) // เปิดใน browser เพื่อดู layout
```

DEMO