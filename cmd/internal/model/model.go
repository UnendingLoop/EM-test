package model

import "time"

// Subscription is a model for storing subscription
type Subscription struct {
	SID      uint       `gorm:"primaryKey" json:"subscription_id"`
	UID      string     `gorm:"not null" json:"user_id"`
	Provider string     `gorm:"not null" json:"service_name"`
	Price    uint       `gorm:"not null" json:"price"`
	Start    time.Time  `gorm:"not null" json:"start_date"`
	End      *time.Time `gorm:"column:end_date" json:"end_date"`
}

// RawSubscription - a model used in handler for basic json-decoding. Converted to model.Subscription in Service-layer.
type RawSubscription struct {
	SID      string `json:"subscription_id"`
	UID      string `json:"user_id"`
	Provider string `json:"service_name"`
	Price    string `json:"price"`
	Start    string `json:"start_date"`
	End      string `json:"end_date"`
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
