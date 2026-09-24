package usecase

import (
	"context"
	"fmt"
	"time"

	"icmongolang/config"
	"icmongolang/internal/modules/purchaseorder"
	pomodels "icmongolang/internal/modules/purchaseorder/models"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"
	"icmongolang/pkg/report"

	"github.com/google/uuid"
)

type purchaseOrderUseCase struct {
	usecase.UseCase[pomodels.PurchaseOrderHeader]
	pgRepo purchaseorder.PurchaseOrderPgRepository
	logger logger.Logger
}

func CreatePurchaseOrderUseCaseI(
	repo purchaseorder.PurchaseOrderPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) purchaseorder.PurchaseOrderUseCaseI {
	return &purchaseOrderUseCase{
		UseCase: usecase.CreateUseCase[pomodels.PurchaseOrderHeader](repo, cfg, logger),
		pgRepo:  repo,
		logger:  logger,
	}
}

func (u *purchaseOrderUseCase) CreateWithDetails(ctx context.Context, header *pomodels.PurchaseOrderHeader, details []*pomodels.PurchaseOrderDetail) (*pomodels.PurchaseOrderHeader, error) {
	u.logger.Infof("Creating purchase order with %d items", len(details))
	created, err := u.pgRepo.CreateWithDetails(ctx, header, details)
	if err != nil {
		u.logger.Errorf("Failed to create purchase order: %v", err)
		return nil, err
	}
	return created, nil
}

func (u *purchaseOrderUseCase) CreateFromQuotation(ctx context.Context, quotationID uuid.UUID, userID uuid.UUID, whitelabelID uuid.UUID) (*pomodels.PurchaseOrderHeader, error) {
	u.logger.Infof("Creating purchase order from quotation: %s", quotationID)
	// TODO: Implement logic when quotation module exists
	// 1. Fetch quotation by ID
	// 2. Verify quotation is approved
	// 3. Create PO header + details from quotation data
	// 4. Calculate totals
	return nil, fmt.Errorf("สร้างใบสั่งซื้อจาก Quotation ยังไม่พร้อมใช้งาน: %w", fmt.Errorf("quotation module not yet implemented"))
}

func (u *purchaseOrderUseCase) Send(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error) {
	u.logger.Infof("Sending purchase order: %s", id)
	po, err := u.pgRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	if po.Status != "DRAFT" {
		return nil, fmt.Errorf("ไม่สามารถส่งใบสั่งซื้อได้ สถานะปัจจุบัน: %s", po.Status)
	}

	if len(po.Details) == 0 {
		return nil, fmt.Errorf("ไม่สามารถส่งใบสั่งซื้อได้ เนื่องจากไม่มีรายการสั่งซื้อ")
	}

	now := timeNow()
	po.Status = "SENT"
	po.SentAt = &now

	if err := u.pgRepo.UpdateStatus(ctx, id, "SENT", map[string]interface{}{
		"sent_at": now,
	}); err != nil {
		return nil, err
	}

	u.saveStatusHistory(ctx, po.ID, "DRAFT", "SENT", po.UserID, "ส่งใบสั่งซื้อให้ Supplier")

	po, err = u.pgRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}
	return po, nil
}

func (u *purchaseOrderUseCase) Confirm(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error) {
	u.logger.Infof("Confirming purchase order: %s", id)
	po, err := u.pgRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	if po.Status != "SENT" {
		return nil, fmt.Errorf("ไม่สามารถยืนยันใบสั่งซื้อได้ สถานะปัจจุบัน: %s", po.Status)
	}

	now := timeNow()
	if err := u.pgRepo.UpdateStatus(ctx, id, "CONFIRMED", map[string]interface{}{
		"confirmed_at": now,
	}); err != nil {
		return nil, err
	}

	u.saveStatusHistory(ctx, po.ID, "SENT", "CONFIRMED", po.UserID, "Supplier ยืนยันใบสั่งซื้อ")

	return u.pgRepo.GetByIDWithDetails(ctx, id)
}

func (u *purchaseOrderUseCase) Receive(ctx context.Context, id uuid.UUID, request *purchaseorder.ReceiveRequest) (*pomodels.PurchaseOrderHeader, error) {
	u.logger.Infof("Receiving goods for purchase order: %s", id)
	po, err := u.pgRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	if po.Status != "SENT" && po.Status != "CONFIRMED" && po.Status != "SHIPPED" {
		return nil, fmt.Errorf("ไม่สามารถรับสินค้าได้ สถานะปัจจุบัน: %s", po.Status)
	}

	detailMap := make(map[uuid.UUID]*pomodels.PurchaseOrderDetail)
	ptrDetails := make([]*pomodels.PurchaseOrderDetail, len(po.Details))
	for i := range po.Details {
		ptrDetails[i] = &po.Details[i]
		detailMap[po.Details[i].ID] = ptrDetails[i]
	}

	for _, item := range request.Items {
		detail, ok := detailMap[item.DetailID]
		if !ok {
			return nil, fmt.Errorf("ไม่พบรายการสั่งซื้อ: %s", item.DetailID)
		}

		newReceived := detail.QuantityReceived + item.ReceivedQuantity
		if newReceived > detail.QuantityOrdered {
			return nil, fmt.Errorf("จำนวนที่รับเกินจำนวนที่สั่ง: %s", item.DetailID)
		}

		if err := u.pgRepo.IncrementReceivedQuantity(ctx, item.DetailID, item.ReceivedQuantity); err != nil {
			return nil, err
		}
		detail.QuantityReceived = newReceived
	}

	allReceived := true
	for _, d := range po.Details {
		if d.QuantityReceived < d.QuantityOrdered {
			allReceived = false
			break
		}
	}

	newStatus := "SHIPPED"
	if allReceived {
		newStatus = "RECEIVED"
		now := timeNow()
		if err := u.pgRepo.UpdateStatus(ctx, id, "RECEIVED", map[string]interface{}{
			"actual_delivery_date": now,
		}); err != nil {
			return nil, err
		}
	} else {
		if err := u.pgRepo.UpdateStatus(ctx, id, newStatus, nil); err != nil {
			return nil, err
		}
	}

	u.saveStatusHistory(ctx, po.ID, po.Status, newStatus, po.UserID, "รับสินค้า")

	return u.pgRepo.GetByIDWithDetails(ctx, id)
}

func (u *purchaseOrderUseCase) Cancel(ctx context.Context, id uuid.UUID, reason string, userID uuid.UUID) (*pomodels.PurchaseOrderHeader, error) {
	u.logger.Infof("Cancelling purchase order: %s", id)
	po, err := u.pgRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, err
	}

	if po.Status == "RECEIVED" || po.Status == "CANCELLED" {
		return nil, fmt.Errorf("ไม่สามารถยกเลิกใบสั่งซื้อได้ สถานะปัจจุบัน: %s", po.Status)
	}

	fromStatus := po.Status
	if err := u.pgRepo.UpdateStatus(ctx, id, "CANCELLED", nil); err != nil {
		return nil, err
	}

	u.saveStatusHistory(ctx, po.ID, fromStatus, "CANCELLED", userID, reason)

	return u.pgRepo.GetByIDWithDetails(ctx, id)
}

func (u *purchaseOrderUseCase) GetByIDWithDetails(ctx context.Context, id uuid.UUID) (*pomodels.PurchaseOrderHeader, error) {
	u.logger.Infof("Getting purchase order: %s", id)
	po, err := u.pgRepo.GetByIDWithDetails(ctx, id)
	if err != nil {
		u.logger.Errorf("Failed to get purchase order %s: %v", id, err)
		return nil, err
	}
	return po, nil
}

func (u *purchaseOrderUseCase) List(ctx context.Context, limit, offset int, supplierID *uuid.UUID, status string, startDate, endDate string) ([]*pomodels.PurchaseOrderHeader, int64, error) {
	u.logger.Infof("Listing purchase orders: limit=%d, offset=%d", limit, offset)
	items, err := u.pgRepo.List(ctx, limit, offset, supplierID, status, startDate, endDate)
	if err != nil {
		u.logger.Errorf("Failed to list purchase orders: %v", err)
		return nil, 0, err
	}
	total, err := u.pgRepo.CountByFilter(ctx, supplierID, status, startDate, endDate)
	if err != nil {
		u.logger.Errorf("Failed to count purchase orders: %v", err)
		return nil, 0, err
	}
	return items, total, nil
}

func (u *purchaseOrderUseCase) GetSuggestions(ctx context.Context, jobID uuid.UUID) ([]*purchaseorder.SuggestionItem, error) {
	u.logger.Infof("Getting suggestions for job: %s", jobID)
	// TODO: Implement when job/quotation modules exist
	return nil, fmt.Errorf("แนะนำอะไหล่จาก Job ยังไม่พร้อมใช้งาน: %w", fmt.Errorf("job module not yet implemented"))
}

func (u *purchaseOrderUseCase) GetStatusHistory(ctx context.Context, poID uuid.UUID) ([]*pomodels.PurchaseOrderStatusHistory, error) {
	u.logger.Infof("Getting status history for purchase order: %s", poID)
	history, err := u.pgRepo.GetStatusHistory(ctx, poID)
	if err != nil {
		u.logger.Errorf("Failed to get status history: %v", err)
		return nil, err
	}
	return history, nil
}

func (u *purchaseOrderUseCase) GetPDF(ctx context.Context, id uuid.UUID) ([]byte, error) {
	u.logger.Infof("Generating PDF for purchase order: %s", id)
	po, err := u.GetByIDWithDetails(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get purchase order for PDF: %w", err)
	}
	var deliveryDate string
	if po.ExpectedDeliveryDate != nil {
		deliveryDate = po.ExpectedDeliveryDate.Format("02/01/2006")
	}
	var notes string
	if po.Notes != nil {
		notes = *po.Notes
	}
	data := report.PurchaseOrderData{
		Company: report.CompanyInfo{
			Name:    "ICMON Auto Repair",
			Address: "123/4 ถนนสุขุมวิท แขวงคลองเตย เขตคลองเตย กรุงเทพฯ 10110",
			Phone:   "02-123-4567",
			TaxID:   "0123456789012",
		},
		PONo:         po.PONo,
		Date:         po.CreatedAt,
		Supplier:     po.SupplierID.String(),
		DeliveryDate: deliveryDate,
		Subtotal:     po.Subtotal,
		TaxAmount:    po.TaxAmount,
		GrandTotal:   po.Total,
		AmountWords:  report.BahtThai(po.Total),
		Remark:       notes,
	}
	for i, d := range po.Details {
		var desc string
		if d.Note != nil {
			desc = *d.Note
		}
		data.Items = append(data.Items, report.PurchaseOrderItem{
			LineNo:      i + 1,
			PartCode:    d.PartID.String(),
			Description: desc,
			Quantity:    d.QuantityOrdered,
			UnitPrice:   d.UnitPrice,
			TotalPrice:  d.TotalPrice,
		})
	}
	return report.GeneratePDF(ctx, report.TplPurchaseOrder, data)
}

func (u *purchaseOrderUseCase) CountByFilter(ctx context.Context, supplierID *uuid.UUID, status string, startDate, endDate string) (int64, error) {
	return u.pgRepo.CountByFilter(ctx, supplierID, status, startDate, endDate)
}

func (u *purchaseOrderUseCase) saveStatusHistory(ctx context.Context, poHeaderID uuid.UUID, fromStatus, toStatus string, changedBy uuid.UUID, reason string) {
	history := &pomodels.PurchaseOrderStatusHistory{
		PoHeaderID:  poHeaderID,
		ToStatus:    toStatus,
		ChangedBy:   changedBy,
		ChangedAt:   timeNow(),
	}
	if fromStatus != "" {
		history.FromStatus = &fromStatus
	}
	if reason != "" {
		history.Reason = &reason
	}
	if err := u.pgRepo.SaveStatusHistory(ctx, history); err != nil {
		u.logger.Errorf("Failed to save status history: %v", err)
	}
}

func timeNow() time.Time {
	return time.Now().UTC()
}
