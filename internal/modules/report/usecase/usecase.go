package usecase

import (
	"context"
	"fmt"

	"icmongolang/internal/modules/customer"
	"icmongolang/internal/modules/quotation"
	"icmongolang/internal/modules/report"
	reportpkg "icmongolang/pkg/report"

	"github.com/google/uuid"
)

type ReportUseCaseI interface {
	BuildInvoice(ctx context.Context, quotationID uuid.UUID) (*reportpkg.InvoiceData, error)
}

type reportUseCase struct {
	customerUC  customer.CustomerUseCaseI
	quotationUC quotation.QuotationUseCaseI
}

func NewReportUseCase(customerUC customer.CustomerUseCaseI, quotationUC quotation.QuotationUseCaseI) ReportUseCaseI {
	return &reportUseCase{customerUC: customerUC, quotationUC: quotationUC}
}

// BuildInvoice assembles an invoice document from a stored quotation.
// Line items come from the quotation's parts and services; part/service
// names are not yet modelled, so their references are used as descriptions.
func (u *reportUseCase) BuildInvoice(ctx context.Context, quotationID uuid.UUID) (*reportpkg.InvoiceData, error) {
	q, err := u.quotationUC.Get(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("load quotation: %w", err)
	}
	if q == nil {
		return nil, fmt.Errorf("quotation %s not found", quotationID)
	}

	parts, err := u.quotationUC.GetQuotationParts(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("load quotation parts: %w", err)
	}
	services, err := u.quotationUC.GetQuotationServices(ctx, quotationID)
	if err != nil {
		return nil, fmt.Errorf("load quotation services: %w", err)
	}

	customerName := ""
	customerAddr := ""
	customer, err := u.customerUC.Get(ctx, q.CustomerID)
	if err == nil && customer != nil {
		customerName = customer.FullName
		if customer.Address != nil {
			customerAddr = *customer.Address
		}
	}

	items := make([]reportpkg.InvoiceItem, 0, len(parts)+len(services))
	line := 0
	for _, p := range parts {
		line++
		desc := "Part ref: " + p.PartID.String()
		if p.Note != nil && *p.Note != "" {
			desc = *p.Note
		}
		items = append(items, reportpkg.InvoiceItem{
			LineNo:      line,
			Description: desc,
			Quantity:    p.Quantity,
			UnitPrice:   p.UnitPrice,
			TotalPrice:  p.UnitPrice * float64(p.Quantity),
		})
	}
	for _, s := range services {
		line++
		desc := "Service ref: " + s.ServiceID.String()
		if s.Note != nil && *s.Note != "" {
			desc = *s.Note
		}
		items = append(items, reportpkg.InvoiceItem{
			LineNo:      line,
			Description: desc,
			Quantity:    s.Quantity,
			UnitPrice:   s.UnitPrice,
			TotalPrice:  s.UnitPrice * float64(s.Quantity),
		})
	}

	amountWords := ""
	if q.AmountInWordsTh != nil {
		amountWords = *q.AmountInWordsTh
	}

	return &reportpkg.InvoiceData{
		Company:       report.CompanyInfo(),
		InvoiceNo:     q.QuotationNo,
		Date:          q.QuotationDate,
		DueDate:       q.QuotationDate.AddDate(0, 0, 30),
		CustomerName:  customerName,
		CustomerAddr:  customerAddr,
		Items:         items,
		Subtotal:      q.Subtotal,
		Discount:      q.DiscountValue,
		TaxAmount:     q.TaxAmount,
		GrandTotal:    q.Total,
		AmountWords:   amountWords,
		PaymentStatus: q.Status,
	}, nil
}
