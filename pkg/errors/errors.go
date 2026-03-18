package errors

import (
	"fmt"
	"net/http"

	"github.com/kongnakornna/golangapi/pkg/logger"
)

var log = logger.Default()

// ErrorType คือประเภทของข้อผิดพลาด (enum)
type ErrorType string

// ค่าคงที่สำหรับประเภทข้อผิดพลาดที่กำหนดไว้ล่วงหน้า
const (
	// ErrorTypeValidation ข้อผิดพลาดจากการตรวจสอบความถูกต้อง
	ErrorTypeValidation ErrorType = "VALIDATION_ERROR"
	// ErrorTypeNotFound ไม่พบทรัพยากร
	ErrorTypeNotFound ErrorType = "NOT_FOUND"
	// ErrorTypeUnauthorized ไม่ได้รับอนุญาต (ยืนยันตัวตน)
	ErrorTypeUnauthorized ErrorType = "UNAUTHORIZED"
	// ErrorTypeForbidden ถูกห้ามเข้าถึง
	ErrorTypeForbidden ErrorType = "FORBIDDEN"
	// ErrorTypeInternal ข้อผิดพลาดภายในเซิร์ฟเวอร์
	ErrorTypeInternal ErrorType = "INTERNAL_ERROR"
	// ErrorTypeBadRequest คำขอไม่ถูกต้อง
	ErrorTypeBadRequest ErrorType = "BAD_REQUEST"
	// ErrorTypeConflict ทรัพยากรขัดแย้งกัน
	ErrorTypeConflict ErrorType = "CONFLICT"
)
const (
	// ErrorTypeValidation ข้อผิดพลาดจากการตรวจสอบความถูกต้อง  
	ErrorTypeValidationTh ErrorType = "ข้อผิดพลาดจากการตรวจสอบความถูกต้อง"  
	// ErrorTypeNotFound ไม่พบทรัพยากร
	ErrorTypeNotFoundTh ErrorType = "ไม่พบทรัพยากร"
	// ErrorTypeUnauthorized ไม่ได้รับอนุญาต (ยืนยันตัวตน)
	ErrorTypeUnauthorizedTh ErrorType = "ไม่ได้รับอนุญาต"
	// ErrorTypeForbidden 
	ErrorTypeForbiddenTh ErrorType = "ไม่อนุญาติให้เข้าถึง"
	// ErrorTypeInternal 
	ErrorTypeInternalTh ErrorType = "ข้อผิดพลาดภายในเซิร์ฟเวอร์"
	// ErrorTypeBadRequest 
	ErrorTypeBadRequestTh ErrorType = "คำขอไม่ถูกต้อง"
	// ErrorTypeConflict ทรัพยากรขัดแย้งกัน
	ErrorTypeConflictTh ErrorType = "ระบบขัดแย้งกัน"
)
// Error โครงสร้างข้อผิดพลาดแบบมีโครงสร้าง
type Error struct {
	Type    ErrorType `json:"type"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
}

// Error ใช้อินเทอร์เฟซมาตรฐาน error
func (e *Error) Error() string {
	// ในสภาพแวดล้อมการผลิต จะไม่แสดงรายละเอียดข้อผิดพลาดภายใน
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap ใช้อินเทอร์เฟซ errors.Unwrap
func (e *Error) Unwrap() error {
	return e.Err
}

// StatusCode คืนค่า HTTP status code ที่สอดคล้องกัน
func (e *Error) StatusCode() int {
	switch e.Type {
	case ErrorTypeValidation:
		return http.StatusBadRequest
	case ErrorTypeNotFound:
		return http.StatusNotFound
	case ErrorTypeUnauthorized:
		return http.StatusUnauthorized
	case ErrorTypeForbidden:
		return http.StatusForbidden
	case ErrorTypeInternal:
		return http.StatusInternalServerError
	case ErrorTypeBadRequest:
		return http.StatusBadRequest
	case ErrorTypeConflict:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

// New สร้างข้อผิดพลาดใหม่
func New(errType ErrorType, message string, err error) *Error {
	return &Error{
		Type:    errType,
		Message: message,
		Err:     err,
	}
}

// ValidationError สร้างข้อผิดพลาดประเภทการตรวจสอบความถูกต้อง
func ValidationError(message string, err error) *Error {
	return New(ErrorTypeValidation, message, err)
}

// NotFoundError สร้างข้อผิดพลาดประเภทไม่พบข้อมูล
func NotFoundError(entity string, err error) *Error {
	return New(ErrorTypeNotFound, fmt.Sprintf("%s not found", entity), err)
}

// UnauthorizedError สร้างข้อผิดพลาดประเภทไม่ได้รับอนุญาต
func UnauthorizedError(message string, err error) *Error {
	return New(ErrorTypeUnauthorized, message, err)
}

// ForbiddenError สร้างข้อผิดพลาดประเภทถูกห้ามเข้าถึง
func ForbiddenError(message string, err error) *Error {
	return New(ErrorTypeForbidden, message, err)
}

// InternalError สร้างข้อผิดพลาดภายในเซิร์ฟเวอร์
func InternalError(message string, err error) *Error {
	return New(ErrorTypeInternal, message, err)
}

// BadRequestError สร้างข้อผิดพลาดประเภทคำขอไม่ถูกต้อง
func BadRequestError(message string, err error) *Error {
	return New(ErrorTypeBadRequest, message, err)
}

// ConflictError สร้างข้อผิดพลาดประเภททรัพยากรขัดแย้ง
func ConflictError(message string, err error) *Error {
	return New(ErrorTypeConflict, message, err)
}

// AsError พยายามแปลง error มาตรฐานเป็น Error type ที่กำหนดเอง
func AsError(err error) *Error {
	if err == nil {
		return nil
	}

	// ถ้า err เป็นประเภท *Error อยู่แล้ว ให้คืนค่าโดยตรง
	if e, ok := err.(*Error); ok {
		return e
	}

	// มิฉะนั้นจะห่อหุ้มเป็นข้อผิดพลาดภายใน
	return InternalError("unexpected error", err)
}

// RecoverPanic ใช้สำหรับกู้คืนจาก panic และบันทึกข้อผิดพลาด
func RecoverPanic(source string) {
	if r := recover(); r != nil {
		// ในสภาพแวดล้อมการผลิต จะบันทึกเฉพาะข้อมูลที่จำเป็น
		log.Error("panic recovered",
			"source", source,
			"error", fmt.Sprintf("%v", r))
	}
}

// RecoverPanicWithCallback กู้คืนจาก panic และเรียกใช้ฟังก์ชัน callback
func RecoverPanicWithCallback(source string, callback func(err interface{})) {
	if r := recover(); r != nil {
		log.Error("panic recovered",
			"source", source,
			"error", fmt.Sprintf("%v", r))

		if callback != nil {
			callback(r)
		}
	}
}