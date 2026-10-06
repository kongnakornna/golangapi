// Package pdpa เป็น entry point ของโมดูล PDPA
package pdpa

import (
	"context"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"icmongolang/config"
	"icmongolang/internal/modules/pdpa/application/command"
	"icmongolang/internal/modules/pdpa/application/port"
	"icmongolang/internal/modules/pdpa/application/query"
	domainrepo "icmongolang/internal/modules/pdpa/domain/repository"
	"icmongolang/internal/modules/pdpa/domain/service"
	"icmongolang/internal/modules/pdpa/infrastructure/persistence/postgres"
	"icmongolang/internal/modules/pdpa/infrastructure/scheduler"
	pdpahandlers "icmongolang/internal/modules/pdpa/interfaces/handlers"
	pdpahttp "icmongolang/internal/modules/pdpa/interfaces/http"
	"icmongolang/pkg/kafka"
	"icmongolang/pkg/logger"
)

type Module struct {
	db         *gorm.DB
	rdb        *redis.Client
	producer   *kafka.Producer
	dispatcher *scheduler.OutboxDispatcherJob

	ConsentRepo domainrepo.ConsentRepository
	DSARRepo    domainrepo.DSARRepository
	AccountRepo domainrepo.UserAccountStatusRepository
	AuditRepo   domainrepo.AuditRepository
	PolicyRepo  domainrepo.PolicyRepository
	OutboxRepo  domainrepo.OutboxRepository

	GrantConsent      *command.GrantConsentCommand
	RevokeConsent     *command.RevokeConsentCommand
	SubmitDSAR        *command.SubmitDSARRequestCommand
	VerifyOTP         *command.VerifyDSAROTPCommand
	ProcessDSAR       *command.ProcessDSARCommand
	CompleteDSAR      *command.CompleteDSARCommand
	RejectDSAR        *command.RejectDSARCommand
	SuspendAccount    *command.SuspendAccountCommand
	TerminateAccount  *command.TerminateAccountCommand
	ConfirmDeletion   *command.ConfirmDeletionCommand
	ImmediateDeletion *command.ImmediateDeletionCommand
	AutoDeleteExpired *command.AutoDeleteExpiredCommand
	PublishPolicy     *command.PublishPrivacyPolicyCommand

	GetDSARStatus   *query.GetDSARStatusQuery
	ListConsentLog  *query.ListConsentLogByUserQuery
	ListDSAR        *query.ListDSARByUserQuery
	GetAdminReport  *query.GetAdminReportQuery
	ListAuditTrails *query.ListAuditTrailsQuery
	GetActivePolicy *query.GetActivePrivacyPolicyQuery

	Handlers *HandlerSet
}

type HandlerSet struct {
	ConsentGranted    *pdpahandlers.ConsentGrantedHandler
	ConsentRevoked    *pdpahandlers.ConsentRevokedHandler
	DSARSubmitted     *pdpahandlers.DSARSubmittedHandler
	DSARCompleted     *pdpahandlers.DSARCompletedHandler
	AccountSuspended  *pdpahandlers.AccountSuspendedHandler
	AccountTerminated *pdpahandlers.AccountTerminatedHandler
	DataDeletion      *pdpahandlers.DataDeletionRequestedHandler
	AuditTrail        *pdpahandlers.AuditTrailHandler
	DLQ               *pdpahandlers.DLQHandler
}

func NewModule(db *gorm.DB, rdb *redis.Client, cfg *config.Config, log logger.Logger) (*Module, error) {
	consentRepo := postgres.NewConsentPgRepository(db)
	dsarRepo := postgres.NewDSARPgRepository(db)
	accountRepo := postgres.NewAccountStatusPgRepository(db)
	auditRepo := postgres.NewAuditPgRepository(db)
	policyRepo := postgres.CreatePolicyPgRepository(db)
	outboxRepo := postgres.CreateOutboxPgRepository(db)

	deletionPolicy := service.NewDeletionPolicyService()

	producer, err := kafka.NewProducer(cfg.Kafka.Brokers, cfg.Kafka.Topic)
	if err != nil {
		return nil, err
	}

	dispatcher := scheduler.NewOutboxDispatcherJob(
		outboxRepo, producer, port.SystemClock{}, 5*time.Second,
	)

	grantCmd := command.NewGrantConsentCommand(consentRepo, nil, producer, cfg, log)
	revokeCmd := command.NewRevokeConsentCommand(consentRepo, nil, producer, cfg, log)
	submitDSAR := command.NewSubmitDSARRequestCommand(dsarRepo, producer, cfg, log)
	verifyOTP := command.NewVerifyDSAROTPCommand(dsarRepo, producer, cfg, log)
	processDSAR := command.NewProcessDSARCommand(dsarRepo, producer, cfg, log)
	completeDSAR := command.NewCompleteDSARCommand(dsarRepo, outboxRepo, producer, cfg, log)
	rejectDSAR := command.NewRejectDSARCommand(dsarRepo, producer, cfg, log)
	suspendAcct := command.NewSuspendAccountCommand(accountRepo, producer, cfg, log)
	terminateAcct := command.NewTerminateAccountCommand(accountRepo, producer, cfg, log)
	confirmDel := command.NewConfirmDeletionCommand(accountRepo, deletionPolicy, producer, cfg, log)
	immediateDel := command.NewImmediateDeletionCommand(accountRepo, consentRepo, dsarRepo, deletionPolicy, producer, cfg, log)
	autoDelete := command.NewAutoDeleteExpiredCommand(accountRepo, consentRepo, producer, cfg, log)
	publishPolicy := command.NewPublishPrivacyPolicyCommand(policyRepo, producer, cfg, log)

	getDSAR := query.NewGetDSARStatusQuery(dsarRepo, log)
	listConsent := query.NewListConsentLogByUserQuery(consentRepo, log)
	listDSAR := query.NewListDSARByUserQuery(dsarRepo, log)
	getReport := query.NewGetAdminReportQuery(auditRepo, log)
	listAudit := query.NewListAuditTrailsQuery(auditRepo, log)
	getPolicy := query.NewGetActivePrivacyPolicyQuery(policyRepo, log)

	handlers := &HandlerSet{
		ConsentGranted:    pdpahandlers.NewConsentGrantedHandler(auditRepo, log),
		ConsentRevoked:    pdpahandlers.NewConsentRevokedHandler(auditRepo, log),
		DSARSubmitted:     pdpahandlers.NewDSARSubmittedHandler(auditRepo, log),
		DSARCompleted:     pdpahandlers.NewDSARCompletedHandler(auditRepo, log),
		AccountSuspended:  pdpahandlers.NewAccountSuspendedHandler(auditRepo, log),
		AccountTerminated: pdpahandlers.NewAccountTerminatedHandler(auditRepo, log),
		DataDeletion:      pdpahandlers.NewDataDeletionRequestedHandler(auditRepo, log),
		AuditTrail:        pdpahandlers.NewAuditTrailHandler(auditRepo, log),
		DLQ:               pdpahandlers.NewDLQHandler(auditRepo, log),
	}

	return &Module{
		db: db, rdb: rdb, producer: producer, dispatcher: dispatcher,
		ConsentRepo: consentRepo, DSARRepo: dsarRepo, AccountRepo: accountRepo,
		AuditRepo: auditRepo, PolicyRepo: policyRepo, OutboxRepo: outboxRepo,
		GrantConsent: grantCmd, RevokeConsent: revokeCmd, SubmitDSAR: submitDSAR,
		VerifyOTP: verifyOTP, ProcessDSAR: processDSAR, CompleteDSAR: completeDSAR,
		RejectDSAR: rejectDSAR, SuspendAccount: suspendAcct, TerminateAccount: terminateAcct,
		ConfirmDeletion: confirmDel, ImmediateDeletion: immediateDel,
		AutoDeleteExpired: autoDelete, PublishPolicy: publishPolicy,
		GetDSARStatus: getDSAR, ListConsentLog: listConsent, ListDSAR: listDSAR,
		GetAdminReport: getReport, ListAuditTrails: listAudit, GetActivePolicy: getPolicy,
		Handlers: handlers,
	}, nil
}

func (m *Module) RegisterRoutes(r chi.Router, handlers *pdpahttp.Handlers, auth, adminOnly pdpahttp.AuthMiddleware) {
	pdpahttp.RegisterRoutes(r, handlers, auth, adminOnly)
}

func (m *Module) StartBackground(ctx context.Context) {
	go m.dispatcher.Start(ctx)
}

func (m *Module) Shutdown() error {
	return m.producer.Close()
}
