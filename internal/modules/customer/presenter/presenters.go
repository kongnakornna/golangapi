package presenter

import (
	"time"

	"icmongolang/pkg/helpers"

	"github.com/google/uuid"
)

// ─── Customer DTOs ────────────────────────────────────────────────────────

type CustomerCreate struct {
	CustomerCode   string  `json:"customerCode" validate:"required" example:"CUST001"`
	FullName       string  `json:"fullName" validate:"required" example:"สมชาย ใจดี"`
	DisplayName    *string `json:"displayName,omitempty" example:"คุณสมชาย"`
	CustomerType   string  `json:"customerType" example:"INDIVIDUAL"`
	TaxID          *string `json:"taxId,omitempty" example:"1234567890123"`
	Email          *string `json:"email,omitempty" example:"somchai@example.com"`
	PhoneNumber    string  `json:"phoneNumber" validate:"required" example:"0812345678"`
	SecondaryPhone *string `json:"secondaryPhone,omitempty"`
	Address        *string `json:"address,omitempty"`
	Province       *string `json:"province,omitempty"`
	City           *string `json:"city,omitempty"`
	District       *string `json:"district,omitempty"`
	PostalCode     *string `json:"postalCode,omitempty"`
	Country        *string `json:"country,omitempty" example:"Thailand"`
	ContactPerson  *string `json:"contactPerson,omitempty"`
	ContactPhone   *string `json:"contactPhone,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}

type CustomerUpdate struct {
	CustomerCode   *string `json:"customerCode,omitempty"`
	FullName       *string `json:"fullName,omitempty"`
	DisplayName    *string `json:"displayName,omitempty"`
	CustomerType   *string `json:"customerType,omitempty"`
	Status         *string `json:"status,omitempty"`
	TaxID          *string `json:"taxId,omitempty"`
	Email          *string `json:"email,omitempty"`
	PhoneNumber    *string `json:"phoneNumber,omitempty"`
	SecondaryPhone *string `json:"secondaryPhone,omitempty"`
	Address        *string `json:"address,omitempty"`
	Province       *string `json:"province,omitempty"`
	City           *string `json:"city,omitempty"`
	District       *string `json:"district,omitempty"`
	PostalCode     *string `json:"postalCode,omitempty"`
	Country        *string `json:"country,omitempty"`
	ContactPerson  *string `json:"contactPerson,omitempty"`
	ContactPhone   *string `json:"contactPhone,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}

type CustomerResponse struct {
	ID              uuid.UUID          `json:"id,omitempty"`
	CustomerCode    string             `json:"customerCode,omitempty"`
	FullName        string             `json:"fullName,omitempty"`
	DisplayName     *string            `json:"displayName,omitempty"`
	CustomerType    string             `json:"customerType,omitempty"`
	Status          string             `json:"status,omitempty"`
	TaxID           *string            `json:"taxId,omitempty"`
	Email           *string            `json:"email,omitempty"`
	PhoneNumber     string             `json:"phoneNumber,omitempty"`
	SecondaryPhone  *string            `json:"secondaryPhone,omitempty"`
	Address         *string            `json:"address,omitempty"`
	Province        *string            `json:"province,omitempty"`
	City            *string            `json:"city,omitempty"`
	District        *string            `json:"district,omitempty"`
	PostalCode      *string            `json:"postalCode,omitempty"`
	Country         *string            `json:"country,omitempty"`
	ContactPerson   *string            `json:"contactPerson,omitempty"`
	ContactPhone    *string            `json:"contactPhone,omitempty"`
	Notes           *string            `json:"notes,omitempty"`
	LastVisitDate   *time.Time         `json:"lastVisitDate,omitempty"`
	TotalVisitCount int                `json:"totalVisitCount,omitempty"`
	TotalSpent      float64            `json:"totalSpent,omitempty"`
	UserID          uuid.UUID          `json:"userId,omitempty"`
	WhitelabelID    uuid.UUID          `json:"whitelabelId,omitempty"`
	CreatedAt       helpers.LocalTime  `json:"createdAt,omitempty"`
	UpdatedAt       *helpers.LocalTime `json:"updatedAt,omitempty"`
}

type PaginatedCustomersResponse struct {
	Customers  []*CustomerResponse `json:"customers"`
	Total      int64               `json:"total"`
	Page       int                 `json:"page"`
	PerPage    int                 `json:"per_page"`
	TotalPages int                 `json:"total_pages"`
}

// ─── Car DTOs ─────────────────────────────────────────────────────────────

type CarCreate struct {
	CustomerID         uuid.UUID  `json:"customerId" validate:"required"`
	LicensePlate       string     `json:"licensePlate" validate:"required" example:"กข1234"`
	Province           *string    `json:"province,omitempty"`
	Brand              string     `json:"brand" validate:"required" example:"Toyota"`
	Model              string     `json:"model" validate:"required" example:"Camry"`
	SubModel           *string    `json:"subModel,omitempty"`
	Year               *int       `json:"year,omitempty"`
	Color              *string    `json:"color,omitempty"`
	EngineNumber       *string    `json:"engineNumber,omitempty"`
	ChassisNumber      *string    `json:"chassisNumber,omitempty"`
	FuelType           *string    `json:"fuelType,omitempty"`
	TransmissionType   *string    `json:"transmissionType,omitempty"`
	EngineCC           *int       `json:"engineCc,omitempty"`
	SeatingCapacity    *int       `json:"seatingCapacity,omitempty"`
	Mileage            int        `json:"mileage,omitempty"`
	LastServiceDate    *time.Time `json:"lastServiceDate,omitempty"`
	NextServiceMileage *int       `json:"nextServiceMileage,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
}

type CarUpdate struct {
	LicensePlate       *string    `json:"licensePlate,omitempty"`
	Province           *string    `json:"province,omitempty"`
	Brand              *string    `json:"brand,omitempty"`
	Model              *string    `json:"model,omitempty"`
	SubModel           *string    `json:"subModel,omitempty"`
	Year               *int       `json:"year,omitempty"`
	Color              *string    `json:"color,omitempty"`
	EngineNumber       *string    `json:"engineNumber,omitempty"`
	ChassisNumber      *string    `json:"chassisNumber,omitempty"`
	FuelType           *string    `json:"fuelType,omitempty"`
	TransmissionType   *string    `json:"transmissionType,omitempty"`
	EngineCC           *int       `json:"engineCc,omitempty"`
	SeatingCapacity    *int       `json:"seatingCapacity,omitempty"`
	Mileage            *int       `json:"mileage,omitempty"`
	LastServiceDate    *time.Time `json:"lastServiceDate,omitempty"`
	NextServiceMileage *int       `json:"nextServiceMileage,omitempty"`
	Notes              *string    `json:"notes,omitempty"`
	IsActive           *bool      `json:"isActive,omitempty"`
}

type CarResponse struct {
	ID                 uuid.UUID          `json:"id,omitempty"`
	CustomerID         uuid.UUID          `json:"customerId,omitempty"`
	LicensePlate       string             `json:"licensePlate,omitempty"`
	Province           *string            `json:"province,omitempty"`
	Brand              string             `json:"brand,omitempty"`
	Model              string             `json:"model,omitempty"`
	SubModel           *string            `json:"subModel,omitempty"`
	Year               *int               `json:"year,omitempty"`
	Color              *string            `json:"color,omitempty"`
	EngineNumber       *string            `json:"engineNumber,omitempty"`
	ChassisNumber      *string            `json:"chassisNumber,omitempty"`
	FuelType           *string            `json:"fuelType,omitempty"`
	TransmissionType   *string            `json:"transmissionType,omitempty"`
	EngineCC           *int               `json:"engineCc,omitempty"`
	SeatingCapacity    *int               `json:"seatingCapacity,omitempty"`
	Mileage            int                `json:"mileage,omitempty"`
	LastServiceDate    *time.Time         `json:"lastServiceDate,omitempty"`
	NextServiceMileage *int               `json:"nextServiceMileage,omitempty"`
	Notes              *string            `json:"notes,omitempty"`
	IsActive           bool               `json:"isActive,omitempty"`
	UserID             uuid.UUID          `json:"userId,omitempty"`
	WhitelabelID       uuid.UUID          `json:"whitelabelId,omitempty"`
	CreatedAt          helpers.LocalTime  `json:"createdAt,omitempty"`
	UpdatedAt          *helpers.LocalTime `json:"updatedAt,omitempty"`
}

type PaginatedCarsResponse struct {
	Cars       []*CarResponse `json:"cars"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PerPage    int            `json:"per_page"`
	TotalPages int            `json:"total_pages"`
}
