package report

import "time"

type CompanyInfo struct {
	Name    string
	Address string
	Phone   string
	TaxID   string
	LogoURL string
}

type QuotationData struct {
	Company     CompanyInfo
	QuotationNo string
	Date        time.Time
	ExpiryDate  time.Time
	CustomerName  string
	CustomerAddr  string
	CustomerPhone string
	JobNo         string
	LicensePlate  string
	CarModel      string
	Items         []QuotationItem
	Subtotal      float64
	Discount      float64
	TaxRate       float64
	TaxAmount     float64
	GrandTotal    float64
	AmountWords   string
	Remark        string
	CreatedBy     string
}

type QuotationItem struct {
	LineNo      int
	Code        string
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type PartPickingData struct {
	Company     CompanyInfo
	PickingNo   string
	JobNo       string
	RequestDate string
	RequestedBy string
	Mechanic    string
	LicensePlate string
	Items       []PartPickingItem
}

type PartPickingItem struct {
	LineNo      int
	PartCode    string
	PartName    string
	QtyRequest  int
	QtyPicked   int
	Location    string
}

type ReceiptData struct {
	Company       CompanyInfo
	ReceiptNo     string
	Date          time.Time
	CustomerName  string
	PaymentMethod string
	Items         []ReceiptItem
	Amount        float64
	AmountWords   string
}

type ReceiptItem struct {
	Description string
	Total       float64
}

type DeliverySheetData struct {
	Company       CompanyInfo
	DeliveryNo    string
	Date          time.Time
	CustomerName  string
	CustomerAddr  string
	JobNo         string
	LicensePlate  string
	Items         []DeliveryItem
	DeliveredBy   string
	ReceivedBy    string
	Remark        string
}

type DeliveryItem struct {
	LineNo      int
	Description string
	Quantity    int
	Unit        string
	Note        string
}

type CreditNoteData struct {
	Company       CompanyInfo
	CreditNoteNo  string
	Date          time.Time
	CustomerName  string
	CustomerAddr  string
	InvoiceNo     string
	Reason        string
	Items         []CreditNoteItem
	Subtotal      float64
	TaxAmount     float64
	GrandTotal    float64
	AmountWords   string
}

type CreditNoteItem struct {
	LineNo      int
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type DebitNoteData struct {
	Company      CompanyInfo
	DebitNoteNo  string
	Date         time.Time
	CustomerName string
	CustomerAddr string
	InvoiceNo    string
	Reason       string
	Items        []DebitNoteItem
	Subtotal     float64
	TaxAmount    float64
	GrandTotal   float64
	AmountWords  string
}

type DebitNoteItem struct {
	LineNo      int
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type PurchaseOrderData struct {
	Company    CompanyInfo
	PONo       string
	Date       time.Time
	Supplier   string
	SupplierAddr string
	DeliveryDate string
	Items      []PurchaseOrderItem
	Subtotal   float64
	TaxAmount  float64
	GrandTotal float64
	AmountWords string
	Remark     string
}

type PurchaseOrderItem struct {
	LineNo      int
	PartCode    string
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type JobCardData struct {
	Company      CompanyInfo
	JobNo        string
	Date         time.Time
	CustomerName string
	CustomerAddr string
	LicensePlate string
	CarBrand     string
	CarModel     string
	CarYear      int
	Mileage      int
	Symptoms     []string
	Diagnosis    string
	Services     []JobCardService
	Parts        []JobCardPart
	Mechanic     string
	Status       string
	Remark       string
}

type JobCardService struct {
	LineNo      int
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type JobCardPart struct {
	LineNo      int
	PartCode    string
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type InvoiceData struct {
	Company       CompanyInfo
	InvoiceNo     string
	Date          time.Time
	DueDate       time.Time
	CustomerName  string
	CustomerAddr  string
	CustomerPhone string
	JobNo         string
	LicensePlate  string
	Items         []InvoiceItem
	Subtotal      float64
	Discount      float64
	TaxAmount     float64
	GrandTotal    float64
	AmountWords   string
	PaymentStatus string
	QRCodeURL     string
	Remark        string
}

type InvoiceItem struct {
	LineNo      int
	Description string
	Quantity    int
	UnitPrice   float64
	TotalPrice  float64
}

type DailySalesData struct {
	Company   CompanyInfo
	Date      time.Time
	Sales     []DailySaleRow
	TotalRev  float64
	TotalCost float64
	TotalProfit float64
	Summary   DailySalesSummary
}

type DailySaleRow struct {
	InvoiceNo  string
	Customer   string
	Amount     float64
	Cost       float64
	Profit     float64
	PaymentMethod string
	Time       string
}

type DailySalesSummary struct {
	InvoiceCount int
	CashCount    int
	CashAmount   float64
	TransferCount int
	TransferAmount float64
	CardCount    int
	CardAmount   float64
	CreditCount  int
	CreditAmount float64
}

type InventorySummaryData struct {
	Company    CompanyInfo
	Date       time.Time
	Items      []InventoryItem
	TotalValue float64
	TotalQty   int
}

type InventoryItem struct {
	Code        string
	Name        string
	Category    string
	QtyOnHand   int
	MinStock    int
	MaxStock    int
	UnitPrice   float64
	TotalValue  float64
	IsLowStock  bool
}

type CustomerListData struct {
	Company CompanyInfo
	Date    time.Time
	Customers []CustomerListItem
	TotalCount int
}

type CustomerListItem struct {
	No          int
	Code        string
	Name        string
	Phone       string
	Email       string
	CarCount    int
	LastVisit   string
	TotalSpent  float64
	Status      string
}
