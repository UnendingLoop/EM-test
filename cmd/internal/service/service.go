package service

import (
	"context"
	"em-test/cmd/internal/model"
	"em-test/cmd/internal/repository"
	"em-test/cmd/internal/utils"
	"errors"
	"fmt"
	"log"
	"strconv"
	"time"

	"gorm.io/gorm"
)

// SubscriptionService provides methods to business logics and further repo(bd-requeste) calls.
type SubscriptionService struct {
	Repo repository.SubscriptionRepo
}

func CreateService(db *gorm.DB) *SubscriptionService {
	return &SubscriptionService{Repo: *repository.CreateRepo(db)}
}

// CreateSubscription - validates input data, checks if such subscription already exists, and if not - creates it in DB via Repository layer.
func (ss *SubscriptionService) CreateSubscription(ctx context.Context, rawSub *model.RawSubscription) error {
	if rawSub.UID == "" || rawSub.Start == "" || rawSub.Provider == "" {
		return fmt.Errorf("Warning on creation: %w", repository.ErrEmptySomeFields)
	}

	newSub, err := utils.ConvertRawSubToNormal(rawSub)
	if err != nil {
		return fmt.Errorf("Convert failure: %w", err)
	}

	err = ss.Repo.CheckIfExists(ctx, newSub)
	if err != nil {
		if errors.Is(err, repository.ErrSubExists) {
			return fmt.Errorf("Failed to create subscription: %w", err)
		}
		if !errors.Is(err, repository.ErrSubNotFound) {
			log.Printf("DB problem while CheckIfExists attempt: %v", err)
			return fmt.Errorf("Creation failed: %w", err)
		}
	}

	err = ss.Repo.CreateSubscription(ctx, newSub)
	if err != nil { //проблема с подключением к базе
		log.Printf("[%v] DB problem while CreateSubscription attempt: %v\nInput data: %v\n", time.Now().Format("2006-01-02 15:04:05"), err, rawSub)
		return err
	}
	rawSub.SID = newSub.SID
	return nil
}

// UpdateBySID - validates data, checks if such SID exists, and if so - updates record in DB via Repository layer.
func (ss *SubscriptionService) UpdateBySID(ctx context.Context, rawSub *model.RawSubscription, sidStr string) error {
	if sidStr == "" {
		return fmt.Errorf("Failed to update subscription: %w", repository.ErrEmptySomeFields)
	}
	if rawSub.UID == "" && rawSub.Provider == "" && rawSub.Price == nil && rawSub.Start == "" && rawSub.End == "" {
		return fmt.Errorf("Failed to update subscription %v: %w", sidStr, repository.ErrEmptyAllFields)
	}
	sid, err := strconv.ParseUint(sidStr, 10, 64)
	if err != nil {
		return fmt.Errorf("Convert failure: %w, %w", utils.ErrConvertToNorm, err)
	}
	rawSub.SID = &sid
	newSub, err := utils.ConvertRawSubToNormal(rawSub)
	if err != nil {
		return fmt.Errorf("Convert failure: %w", err)
	}

	dbSub, err := ss.Repo.GetSubscriptionBySID(ctx, *newSub.SID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { //подписка не найдена
			return fmt.Errorf("Failed to update user info: %w", repository.ErrSubNotFound)
		}
		//проблема с подключением к базе
		log.Printf("[%v] DB problem while GetSubscriptionBySID attempt: %v\nInput data: %v\n", time.Now().Format("2006-01-02 15:04:05"), err, sidStr)
		return err
	}

	if rawSub.UID != "" {
		dbSub.UID = newSub.UID
	}
	if rawSub.Provider != "" {
		dbSub.Provider = newSub.Provider
	}
	if rawSub.Price != nil {
		dbSub.Price = newSub.Price
	}
	if rawSub.Start != "" {
		dbSub.Start = newSub.Start
	}
	if rawSub.End != "" {
		dbSub.End = newSub.End
	}
	if err := ss.Repo.UpdateSubscriptionInfo(ctx, dbSub); err != nil { //проблема с подключением к базе
		log.Printf("[%v] DB problem while UpdateSubscriptionInfo attempt: %v\nInput data: %v\n", time.Now().Format("2006-01-02 15:04:05"), err, rawSub)
		return fmt.Errorf("Failed to update subscription info: %w", err)
	}
	return nil
}

// GetBySID - returns an instance of type model.Subscription if there is a record under provided SID in DB
func (ss *SubscriptionService) GetBySID(ctx context.Context, sid uint64) (*model.RawSubscription, error) {
	dbSub, err := ss.Repo.GetSubscriptionBySID(ctx, sid)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("Failed to get subscription info: %w", repository.ErrSubNotFound)
		}
		//проблема с подключением к базе
		log.Printf("[%v]DB problem while GetSubscriptionBySID attempt: %v\nInput data: %v\n", time.Now().Format("2006-01-02 15:04:05"), err, sid)
		return nil, err
	}
	rawSub := utils.ConvertNormalSubToRaw(dbSub)
	return rawSub, nil
}

// GetList - provides array of all subscription records existing in DB
func (ss *SubscriptionService) GetList(ctx context.Context) ([]*model.RawSubscription, error) {
	dbSubs, err := ss.Repo.GetAllSubscriptions(ctx)
	if err != nil {
		//проблема с подключением к базе
		log.Printf("[%v]DB problem while GetAllSubscriptions attempt: %v\n", time.Now().Format("2006-01-02 15:04:05"), err)
		return nil, err
	}
	if len(dbSubs) == 0 {
		return nil, nil
	}
	rawSubs := make([]*model.RawSubscription, len(dbSubs))
	for i, v := range dbSubs {
		rawSubs[i] = utils.ConvertNormalSubToRaw(v)
	}
	return rawSubs, nil
}

// DeleteSubscription - removes record by SID, returns error if no rows affected
func (ss *SubscriptionService) DeleteSubscription(ctx context.Context, sid uint) error {
	count, err := ss.Repo.DeleteSubcription(ctx, sid)
	if count == 0 && err == nil {
		return fmt.Errorf("Failed to remove susbcription: %w", repository.ErrSubNotFound)
	}
	if count != 0 && err == nil {
		return nil
	}

	//проблема с подключением к базе
	log.Printf("[%v] DB problem while DeleteSubcription attempt: %v\nInput data: %v\n", time.Now().Format("2006-01-02 15:04:05"), err, sid)
	return err
}

// Report - provides a total price of subscriptions which meet the search request: period(mandatory, specific month), UID(optional) and Provider(optional)
func (ss *SubscriptionService) Report(ctx context.Context, filter *model.RawReportFilter) (uint, error) {
	normFilter, err := utils.ConvertFilterToNorm(filter)
	if err != nil {
		return 0, err
	}

	res, err := ss.Repo.ComposeReport(ctx, normFilter)
	if err != nil {
		//проблема с подключением к базе
		log.Printf("[%v] DB problem while ComposeReport attempt: %v\nInput data: %v\n", time.Now().Format("2006-01-02 15:04:05"), err, normFilter)
		return 0, fmt.Errorf("Failed to make report: %w", err)
	}

	return res, nil
}
