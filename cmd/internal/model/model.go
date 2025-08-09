package model

import "time"

// Subscription is a model for storing subscription
type Subscription struct {
	SID      *uint64    `gorm:"column:subscription_id;primaryKey" json:"subscription_id"`
	UID      string     `gorm:"column:user_id;not null" json:"user_id"`
	Provider string     `gorm:"column:service_name;not null" json:"service_name"`
	Price    uint       `gorm:"column:price;not null" json:"price"`
	Start    time.Time  `gorm:"column:start_date;not null" json:"start_date"`
	End      *time.Time `gorm:"column:end_date" json:"end_date"`
}

// RawSubscription - a model used in handler for basic json-decoding. Converted to model.Subscription in Service-layer.
type RawSubscription struct {
	SID      *uint64 `json:"subscription_id" example:"20"`
	UID      string  `json:"user_id" example:"60601fee-2bf1-4721-ae6f-7636e79a0cba"`
	Provider string  `json:"service_name" example:"Yandex Plus"`
	Price    *uint   `json:"price" example:"400"`
	Start    string  `json:"start_date" example:"07-2025"`
	End      string  `json:"end_date,omitempty" example:"12-2025"`
}

// RawReportFilter - a model used for composing report - used only for storing raw data
type RawReportFilter struct {
	Period   string //mandatory, конкретный месяц в формате "07-2024"
	UID      string //optional
	Provider string //optional
}

// ReportFilter - a model used for composing report - used in Repository for query
type ReportFilter struct {
	Start    time.Time //mandatory
	End      time.Time //mandatory
	UID      *string   //optional
	Provider *string   //optional
}

// Report used for responding with subscription total price
type Report struct {
	Total uint `json:"total"`
}
