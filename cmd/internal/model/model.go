package model

import "time"

// Subscription is a model for storing subscription
type Subscription struct {
	UID     string    `gorm:"primaryKey" json:"user_id"`
	Service string    `gorm:"not null" json:"service_name"`
	Price   uint      `gorm:"not null" json:"price"`
	Start   time.Time `gorm:"not null" json:"start_date"`
	End     time.Time `gorm:"column:end_date"`
}

/*
 	“service_name”: “Yandex Plus”,
    “price”: 400,
    “user_id”: “60601fee-2bf1-4721-ae6f-7636e79a0cba”,
    “start_date”: “07-2025”

	 - Название сервиса, предоставляющего подписку
    - Стоимость месячной подписки в рублях
    - ID пользователя в формате UUID
    - Дата начала подписки (месяц и год)
    - Опционально дата окончания подписки
*/
