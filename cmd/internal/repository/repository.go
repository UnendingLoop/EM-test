package repository

import (
	"context"
	"database/sql"
	"em-test/cmd/internal/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// SubscriptionRepo - structure provides access to DB-requests
type SubscriptionRepo struct {
	DB *gorm.DB
}

var ErrSubNotFound = errors.New("subscription not found")
var ErrSubExists = errors.New("subscription already exists")

var ErrEmptyFields = errors.New("all fields are empty")
var ErrEmptySomeFields = errors.New("mandatory fields are empty")

// CreateSubscription -
func (sr SubscriptionRepo) CreateSubscription(ctx context.Context, newSub *model.Subscription) error {
	return sr.DB.WithContext(ctx).Create(&newSub).Error
}

// GetSubscriptionBySID -
func (sr SubscriptionRepo) GetSubscriptionBySID(ctx context.Context, sid uint) (*model.Subscription, error) {
	var dbSub model.Subscription
	err := sr.DB.WithContext(ctx).First(&dbSub, sid).Error
	return &dbSub, err
}

// GetAllSubscriptions -
func (sr SubscriptionRepo) GetAllSubscriptions(ctx context.Context) ([]*model.Subscription, error) {
	var dbSubs []*model.Subscription
	err := sr.DB.WithContext(ctx).Find(&dbSubs).Error
	return dbSubs, err
}

// UpdateSubscriptionInfo -
func (sr SubscriptionRepo) UpdateSubscriptionInfo(ctx context.Context, newSub *model.Subscription) error {
	return sr.DB.WithContext(ctx).Save(&newSub).Error
}

// DeleteSubcription -
func (sr SubscriptionRepo) DeleteSubcription(ctx context.Context, uid uint) (int64, error) {
	res := sr.DB.WithContext(ctx).Delete(&model.Subscription{}, uid)
	return res.RowsAffected, res.Error
}

// ComposeReport provides a total summ of subscription prices that meet requirements of filterSub
func (sr SubscriptionRepo) ComposeReport(ctx context.Context, filterSub *model.ReportFilter) (uint, error) {
	var total sql.NullInt64

	query := sr.DB.WithContext(ctx).Model(&model.Subscription{}).
		Select("SUM(price) as total").
		Where("start_date <= ?", filterSub.End).
		Where("end_date IS NULL OR end_date >= ?", filterSub.Start)

	if filterSub.UID != nil {
		query = query.Where("user_id = ?", filterSub.UID)
	}
	if filterSub.Provider != nil {
		query = query.Where("provider = ?", filterSub.Provider)
	}

	err := query.Scan(&total).Error
	if err != nil {
		return 0, err
	}

	if total.Valid {
		return uint(total.Int64), nil
	}
	return 0, nil
}

// CheckIfExists - checks if subscription data already exists in DB, returns informative error in both cases
func (sr SubscriptionRepo) CheckIfExists(ctx context.Context, candidate *model.Subscription) error {
	var res int64
	err := sr.DB.WithContext(ctx).
		Model(&model.Subscription{}).
		Where("user_id = ? AND service_name = ? AND price =?  AND start_date = ?", candidate.UID, candidate.Provider, candidate.Price, candidate.Start).
		Count(&res).Error
	if err != nil {
		return fmt.Errorf("Failed request: %w", err)
	}
	if res == 0 {
		return ErrSubNotFound
	}
	return ErrSubExists
}
