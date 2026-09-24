package usecase

import (
	"context"
	"fmt"

	"icmongolang/config"
	"icmongolang/internal/models"
	"icmongolang/internal/modules/payment"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

type paymentUseCase struct {
	usecase.UseCase[models.Payment]
	pgRepo payment.PaymentPgRepository
}

func CreatePaymentUseCaseI(
	pgRepo payment.PaymentPgRepository,
	cfg *config.Config,
	logger logger.Logger,
) payment.PaymentUseCaseI {
	return &paymentUseCase{
		UseCase: usecase.CreateUseCase[models.Payment](pgRepo, cfg, logger),
		pgRepo:  pgRepo,
	}
}

func (u *paymentUseCase) Count(ctx context.Context) (int64, error) {
	return u.pgRepo.Count(ctx)
}

func (u *paymentUseCase) GetByInvoiceId(ctx context.Context, invoiceId uuid.UUID) (*models.Payment, error) {
	return u.pgRepo.GetByInvoiceId(ctx, invoiceId)
}

func (u *paymentUseCase) Search(ctx context.Context, filter map[string]interface{}, limit, offset int) ([]*models.Payment, error) {
	return u.pgRepo.Search(ctx, filter, limit, offset)
}

func (u *paymentUseCase) GetOutstandingByCustomerId(ctx context.Context, customerId uuid.UUID) ([]*models.OutstandingBalance, error) {
	return u.pgRepo.GetOutstandingByCustomerId(ctx, customerId)
}

func (u *paymentUseCase) GetPaymentHistory(ctx context.Context, customerId uuid.UUID) ([]*models.PaymentHistory, error) {
	return u.pgRepo.GetPaymentHistory(ctx, customerId)
}

func (u *paymentUseCase) GetReceipt(ctx context.Context, id uuid.UUID) (*models.Receipt, error) {
	return u.pgRepo.GetReceipt(ctx, id)
}

func (u *paymentUseCase) GetReceiptByPaymentId(ctx context.Context, paymentId uuid.UUID) (*models.Receipt, error) {
	return u.pgRepo.GetReceiptByPaymentId(ctx, paymentId)
}

func (u *paymentUseCase) CancelReceipt(ctx context.Context, id uuid.UUID, reason string) error {
	return u.pgRepo.CancelReceipt(ctx, id, reason)
}

func (u *paymentUseCase) ProcessRefund(ctx context.Context, id uuid.UUID, amount float64, reason string) (*models.Payment, error) {
	payment, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	if payment.Status != "PAID" && payment.Status != "PARTIAL" {
		return nil, fmt.Errorf("cannot refund payment in status %s", payment.Status)
	}

	if amount <= 0 {
		return nil, fmt.Errorf("refund amount must be greater than zero")
	}

	maxRefundable := payment.Amount - payment.RefundedAmount
	if amount > maxRefundable {
		return nil, fmt.Errorf("refund amount %.2f exceeds refundable amount %.2f", amount, maxRefundable)
	}

	values := make(map[string]interface{})
	values["status"] = "REFUNDED"
	values["refunded_amount"] = payment.RefundedAmount + amount

	return u.pgRepo.Update(ctx, payment, values)
}

func (u *paymentUseCase) CancelPayment(ctx context.Context, id uuid.UUID, reason string) error {
	payment, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		return err
	}

	if payment.Status == "CANCELLED" || payment.Status == "REFUNDED" {
		return fmt.Errorf("cannot cancel payment in status %s", payment.Status)
	}

	values := make(map[string]interface{})
	values["status"] = "CANCELLED"
	values["notes"] = reason

	_, err = u.pgRepo.Update(ctx, payment, values)
	return err
}
