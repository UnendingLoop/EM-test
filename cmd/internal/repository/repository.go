package repository

import "gorm.io/gorm"

type SubscriptionRepo struct {
	DB *gorm.DB
}
