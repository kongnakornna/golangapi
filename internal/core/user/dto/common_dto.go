package dto

// ListResponse การตอบกลับแบบแบ่งหน้ารายการ
type ListResponse struct {
	Data  interface{} `json:"data"`  // ข้อมูลรายการ
	Page  int         `json:"page"`  // หมายเลขหน้าปัจจุบัน
	Size  int         `json:"size"`  // ขนาดต่อหน้า
	Total int64       `json:"total"` // จำนวนเรคคอร์ดทั้งหมด
}
