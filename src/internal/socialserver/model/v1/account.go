package v1

import (
	metav1 "go-socialapp/internal/pkg/meta/v1"
)

type Account struct {
	// May add TypeMeta in the future.
	// metav1.TypeMeta `json:",inline"`

	// Standard object's metadata.
	metav1.ObjectMetaNotIncrement

	AccountId string `gorm:"column:account_id" json:"accountId,omitempty"`

	// Required: true
	AccountName string `gorm:"column:account_name" json:"accountName,omitempty" `

	// Required: true
	Device string `gorm:"column:device" json:"device,omitempty" `

	// Required: true
	PhoneNumber string `gorm:"column:phone_number" json:"phoneNumber,omitempty"`

	// Required: true
	Country string `gorm:"column:country" json:"country,omitempty"`

	// Required: true
	Status string `gorm:"column:status" json:"status,omitempty"`
}

// TableName maps to mysql table name.
func (u *Account) TableName() string {
	return "account"
}
