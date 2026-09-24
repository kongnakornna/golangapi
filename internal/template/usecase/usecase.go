package usecase

import (
	"context"
	"fmt"

	"icmongolang/config"
	"icmongolang/internal/template"
	"icmongolang/internal/template/models"
	"icmongolang/internal/usecase"
	"icmongolang/pkg/logger"

	"github.com/google/uuid"
)

type templateUseCase struct {
	usecase.UseCase[models.Template]
	pgRepo    template.PgRepository
	redisRepo template.RedisRepository
	logger    logger.Logger
}

func NewTemplateUseCase(
	pgRepo template.PgRepository,
	redisRepo template.RedisRepository,
	cfg *config.Config,
	log logger.Logger,
) template.UseCaseI {
	return &templateUseCase{
		UseCase:   usecase.CreateUseCase[models.Template](pgRepo, cfg, log),
		pgRepo:    pgRepo,
		redisRepo: redisRepo,
		logger:    log,
	}
}

func (u *templateUseCase) Create(ctx context.Context, exp *models.Template) (*models.Template, error) {
	u.logger.Infof("Creating template: %s", exp.Name)
	created, err := u.pgRepo.Create(ctx, exp)
	if err != nil {
		u.logger.Errorf("Failed to create template: %v", err)
		return nil, err
	}
	return created, nil
}

func (u *templateUseCase) Get(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	u.logger.Infof("Getting template: %s", id)
	obj, err := u.pgRepo.Get(ctx, id)
	if err != nil {
		u.logger.Errorf("Failed to get template %s: %v", id, err)
		return nil, err
	}
	return obj, nil
}

func (u *templateUseCase) GetMulti(ctx context.Context, limit, offset int) ([]*models.Template, error) {
	u.logger.Infof("Listing templates: limit=%d, offset=%d", limit, offset)
	objs, err := u.pgRepo.GetMulti(ctx, limit, offset)
	if err != nil {
		u.logger.Errorf("Failed to list templates: %v", err)
		return nil, err
	}
	return objs, nil
}

func (u *templateUseCase) Update(ctx context.Context, id uuid.UUID, values map[string]interface{}) (*models.Template, error) {
	u.logger.Infof("Updating template %s", id)
	updated, err := u.pgRepo.Update(ctx, &models.Template{ID: id}, values)
	if err != nil {
		u.logger.Errorf("Failed to update template %s: %v", id, err)
		return nil, err
	}
	return updated, nil
}

func (u *templateUseCase) Delete(ctx context.Context, id uuid.UUID) (*models.Template, error) {
	u.logger.Infof("Deleting template: %s", id)
	deleted, err := u.pgRepo.Delete(ctx, id)
	if err != nil {
		u.logger.Errorf("Failed to delete template %s: %v", id, err)
		return nil, err
	}
	return deleted, nil
}

// GenerateRedisKey generates a Redis cache key for a template
func (u *templateUseCase) GenerateRedisKey(id uuid.UUID) string {
	return fmt.Sprintf("%s:%s", models.Template{}.TableName(), id.String())
}
