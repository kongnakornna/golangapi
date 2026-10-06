package repository

import (
	"context"

	"icmongolang/internal/modules/users"
	usermodels "icmongolang/internal/modules/users/models"

	"gorm.io/gorm"
)

type logRepo struct{ db *gorm.DB }

func CreateLogPgRepository(db *gorm.DB) users.LogRepository {
	return &logRepo{db: db}
}

func (r *logRepo) Create(ctx context.Context, l *usermodels.SdUserLog) (*usermodels.SdUserLog, error) {
	if err := r.db.WithContext(ctx).Create(l).Error; err != nil { return nil, err }
	return l, nil
}

func (r *logRepo) Get(ctx context.Context, id int, logTypeID int) (*usermodels.SdUserLog, error) {
	var l usermodels.SdUserLog
	if err := r.db.WithContext(ctx).Where("id = ? AND log_type_id = ?", id, logTypeID).First(&l).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *logRepo) List(ctx context.Context, f users.LogFilter) ([]*usermodels.SdUserLog, int64, error) {
	q := r.db.WithContext(ctx).Model(&usermodels.SdUserLog{})
	if f.UID != nil       { q = q.Where("uid = ?", *f.UID) }
	if f.LogTypeID != nil { q = q.Where("log_type_id = ?", *f.LogTypeID) }
	if f.IPAddress != nil { q = q.Where("ipaddress = ?", *f.IPAddress) }
	if f.From != nil      { q = q.Where("create >= ?", *f.From) }
	if f.To != nil        { q = q.Where("create <= ?", *f.To) }

	var total int64
	if err := q.Count(&total).Error; err != nil { return nil, 0, err }
	if f.Limit <= 0 { f.Limit = 20 }
	if f.Offset < 0 { f.Offset = 0 }
	var out []*usermodels.SdUserLog
	if err := q.Order("create DESC, id DESC").Limit(f.Limit).Offset(f.Offset).Find(&out).Error; err != nil {
		return nil, 0, err
	}
	return out, total, nil
}

func (r *logRepo) Delete(ctx context.Context, id int, logTypeID int) error {
	res := r.db.WithContext(ctx).Where("id = ? AND log_type_id = ?", id, logTypeID).Delete(&usermodels.SdUserLog{})
	if res.Error != nil { return res.Error }
	if res.RowsAffected == 0 { return gorm.ErrRecordNotFound }
	return nil
}

var _ users.LogRepository = (*logRepo)(nil)