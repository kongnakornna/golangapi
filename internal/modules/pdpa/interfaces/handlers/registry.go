package handlers

import (
	"icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/pkg/logger"
)

type Registry struct {
	ConsentGranted        *ConsentGrantedHandler
	ConsentRevoked        *ConsentRevokedHandler
	DSARSubmitted         *DSARSubmittedHandler
	DSARCompleted         *DSARCompletedHandler
	AuditTrail            *AuditTrailHandler
	AccountSuspended      *AccountSuspendedHandler
	AccountTerminated     *AccountTerminatedHandler
	DataDeletionRequested *DataDeletionRequestedHandler
	DLQ                   *DLQHandler
}

func NewRegistry(auditRepo repository.AuditRepository, log logger.Logger) *Registry {
	return &Registry{
		ConsentGranted:        NewConsentGrantedHandler(auditRepo, log),
		ConsentRevoked:        NewConsentRevokedHandler(auditRepo, log),
		DSARSubmitted:         NewDSARSubmittedHandler(auditRepo, log),
		DSARCompleted:         NewDSARCompletedHandler(auditRepo, log),
		AuditTrail:            NewAuditTrailHandler(auditRepo, log),
		AccountSuspended:      NewAccountSuspendedHandler(auditRepo, log),
		AccountTerminated:     NewAccountTerminatedHandler(auditRepo, log),
		DataDeletionRequested: NewDataDeletionRequestedHandler(auditRepo, log),
		DLQ:                   NewDLQHandler(auditRepo, log),
	}
}

func (r *Registry) RegisterWith(register func(topic string, h HandlerFunc)) {
	register("consent.granted", r.ConsentGranted.Handle)
	register("consent.revoked", r.ConsentRevoked.Handle)
	register("dsar.submitted", r.DSARSubmitted.Handle)
	register("dsar.completed", r.DSARCompleted.Handle)
	register("dsar.processing", r.AuditTrail.Handle)
	register("dsar.otp_verified", r.AuditTrail.Handle)
	register("dsar.rejected", r.AuditTrail.Handle)
	register("account.suspended", r.AccountSuspended.Handle)
	register("account.terminated", r.AccountTerminated.Handle)
	register("data.deletion_requested", r.DataDeletionRequested.Handle)
}
