package service

import (
	"context"
	"em-test/cmd/internal/model"
	"em-test/cmd/internal/repository"
	"errors"
	"fmt"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// SubscriptionService provides methods to business logics and further repo(bd-requeste) calls.
type SubscriptionService struct {
	Repo repository.SubscriptionRepo
}

// ErrConvertToNorm - error reflecting problems while converting raw subscription data from model.RawSubscription to model.Subscription.
var ErrConvertToNorm = errors.New("failed to convert raw data to normal")

// Create - validates input data, checks if such subscription already exists, and if not - creates it in DB via Repository layer.
func (SS *SubscriptionService) Create(ctx context.Context, rawSub *model.RawSubscription) error {
	if rawSub.UID == "" || rawSub.Start == "" || rawSub.Price == "" || rawSub.Provider == "" {
		return fmt.Errorf("Warning on creation: %w", repository.ErrEmptySomeFields)
	}

	newSub, err := convertRawSubToNormal(rawSub)
	if err != nil {
		return fmt.Errorf("Conversion warning: %w", err)
	}

	if err := SS.Repo.CheckIfExists(ctx, newSub); !errors.Is(err, repository.ErrSubNotFound) {
		return fmt.Errorf("Creation warning: %w", err)
	}

	return SS.Repo.CreateSubscription(ctx, newSub)
}

// UpdateBySID - validates data, checks if such SID exists, and if so - updates record in DB via Repository layer.
func (SS *SubscriptionService) UpdateBySID(ctx context.Context, rawSub *model.RawSubscription, sidStr string) error {
	if sidStr == "" {
		return fmt.Errorf("Failed to update subscription: %w", repository.ErrEmptySomeFields)
	}
	if rawSub.UID == "" && rawSub.Provider == "" && rawSub.Price == "" && rawSub.Start == "" && rawSub.End == "" {
		return fmt.Errorf("Failed to update subscription %v: %w", sidStr, repository.ErrEmptyFields)
	}

	rawSub.SID = sidStr
	newSub, err := convertRawSubToNormal(rawSub)
	if err != nil {
		return fmt.Errorf("Convert failure: %w", err)
	}

	dbSub, err := SS.Repo.GetSubscriptionBySID(ctx, newSub.SID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("Failed to update user info: %w", repository.ErrSubNotFound)
		}
		return err
	}

	if rawSub.UID != "" {
		dbSub.UID = newSub.UID
	}
	if rawSub.Provider != "" {
		dbSub.Provider = newSub.Provider
	}
	if rawSub.Price != "" {
		dbSub.Price = newSub.Price
	}
	if rawSub.Start != "" {
		dbSub.Start = newSub.Start
	}
	if rawSub.End != "" {
		dbSub.End = newSub.End
	}
	if err := SS.Repo.UpdateSubscriptionInfo(ctx, dbSub); err != nil {
		return fmt.Errorf("Failed to update subscription info: %w", err)
	}
	return nil
}

// GetBySID - returns an instance of type model.Subscription if there is a record under provided SID in DB
func (SS *SubscriptionService) GetBySID(ctx context.Context, sid uint) (*model.RawSubscription, error) {
	dbSub, err := SS.Repo.GetSubscriptionBySID(ctx, sid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Failed to update user info: %w", repository.ErrSubNotFound)
		}
		return nil, err
	}
	rawSub := convertNormalSubToRaw(dbSub)
	return rawSub, nil
}

// GetList - provides array of all subscription records existing in DB
func (SS *SubscriptionService) GetList(ctx context.Context) ([]*model.RawSubscription, error) {
	dbSubs, err := SS.Repo.GetAllSubscriptions(ctx)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch subscriptions list: %w", err)
	}
	if len(dbSubs) == 0 {
		return nil, nil
	}
	rawSubs := make([]*model.RawSubscription, 0, len(dbSubs))
	for i, v := range dbSubs {
		rawSubs[i] = convertNormalSubToRaw(v)
	}
	return rawSubs, nil
}

// Delete - removes record by SID, returns error if no rows affected
func (SS *SubscriptionService) Delete(ctx context.Context, sid uint) error {
	count, err := SS.Repo.DeleteSubcription(ctx, sid)
	if count == 0 && err == nil {
		return fmt.Errorf("Failed to remove user: %w", repository.ErrSubNotFound)
	}
	return err
}

// Report - provides a total price of subscriptions which meet the search request: period(mandatory, specific month), UID(optional) and Provider(optional)
func (SS *SubscriptionService) Report(ctx context.Context, filter *model.RawReportFilter) (uint, error) {
	normFilter, err := convertFilterToNorm(filter)
	if err != nil {
		return 0, err
	}
	res, err := SS.Repo.ComposeReport(ctx, normFilter)
	if err != nil {
		return 0, fmt.Errorf("Failed to make report: %w", err)
	}
	return res, nil
}

func convertRawSubToNormal(rawSub *model.RawSubscription) (*model.Subscription, error) {
	var normSub model.Subscription
	sid, err0 := strconv.ParseUint(rawSub.SID, 10, 64)
	price, err1 := strconv.ParseUint(rawSub.Price, 10, 64)
	start, err2 := formatTextToTime(rawSub.Start)
	end, err3 := formatTextToTime(rawSub.End)
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, fmt.Errorf("%w: %w\n%w\n%w\n%w", ErrConvertToNorm, err0, err1, err2, err3)
	}
	normSub.SID = uint(sid)
	normSub.UID = rawSub.UID
	normSub.Provider = rawSub.Provider
	normSub.Price = uint(price)
	normSub.Start = *start
	normSub.End = end

	return &normSub, nil
}

func convertNormalSubToRaw(normSub *model.Subscription) *model.RawSubscription {
	var rawSub model.RawSubscription
	rawSub.SID = strconv.Itoa(int(normSub.SID))
	rawSub.UID = normSub.UID
	rawSub.Provider = normSub.Provider
	rawSub.Price = strconv.Itoa(int(normSub.Price))
	rawSub.Start = formatTimeToText(&normSub.Start)
	rawSub.End = formatTimeToText(normSub.End)
	return &rawSub
}

func convertFilterToNorm(rawfilter *model.RawReportFilter) (*model.ReportFilter, error) {
	normFilter := &model.ReportFilter{}
	middleOfMonth, err := formatTextToTime(rawfilter.Period)
	if err != nil {
		return nil, fmt.Errorf("Convert warning: %w: %w", ErrConvertToNorm, err)
	}
	start := middleOfMonth
	end := middleOfMonth.AddDate(0, 0, 10)
	normFilter.UID = &rawfilter.UID
	normFilter.Provider = &rawfilter.Provider
	normFilter.Start = *start
	normFilter.End = end
	return normFilter, nil
}

func formatTextToTime(source string) (*time.Time, error) {
	if source == "" {
		return nil, nil
	}
	format := "01-2006" //Mon Jan 2 15:04:05 MST 2006
	startOfMonth, err := time.Parse(format, source)
	if err != nil {
		return &startOfMonth, fmt.Errorf("failed to convert time: %w", err)
	}
	middleOfMonth := startOfMonth.AddDate(0, 0, 15)
	return &middleOfMonth, nil
}

func formatTimeToText(source *time.Time) string {
	if source == nil {
		return ""
	}
	return source.Format("01-2006")
}
