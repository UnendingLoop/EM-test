package utils

import (
	"em-test/cmd/internal/model"
	"errors"
	"fmt"
	"time"
)

// ErrConvertToNorm - error reflecting problems while converting raw subscription data from model.RawSubscription to model.Subscription.
var ErrConvertToNorm = errors.New("failed to convert raw data to normal")

func ConvertRawSubToNormal(rawSub *model.RawSubscription) (*model.Subscription, error) {
	var normSub model.Subscription
	sid := rawSub.SID
	price := rawSub.Price
	start, err1 := formatTextToTime(rawSub.Start)
	end, err2 := formatTextToTime(rawSub.End)
	normSub.End = end
	if err1 != nil || err2 != nil {
		return nil, fmt.Errorf("%w: %w\n%w", ErrConvertToNorm, err1, err2)
	}
	normSub.SID = sid
	normSub.UID = rawSub.UID
	normSub.Provider = rawSub.Provider

	normSub.Price = *price
	normSub.Start = *start

	return &normSub, nil
}

func ConvertNormalSubToRaw(normSub *model.Subscription) *model.RawSubscription {
	var rawSub model.RawSubscription
	rawSub.SID = normSub.SID
	rawSub.UID = normSub.UID
	rawSub.Provider = normSub.Provider
	rawSub.Price = &normSub.Price
	rawSub.Start = formatTimeToText(&normSub.Start)
	rawSub.End = formatTimeToText(normSub.End)
	return &rawSub
}

func ConvertFilterToNorm(rawfilter *model.RawReportFilter) (*model.ReportFilter, error) {
	normFilter := &model.ReportFilter{}

	// Формат только "01-2006"
	startOfMonth, err := time.Parse("01-2006", rawfilter.Period)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrConvertToNorm, err)
	}

	endOfMonth := startOfMonth.AddDate(0, 1, 0).Add(-time.Nanosecond)

	if rawfilter.UID == "" {
		normFilter.UID = nil
	} else {
		normFilter.UID = &rawfilter.UID
	}
	if rawfilter.Provider == "" {
		normFilter.Provider = nil
	} else {
		normFilter.Provider = &rawfilter.Provider
	}

	normFilter.Start = startOfMonth
	normFilter.End = endOfMonth

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
