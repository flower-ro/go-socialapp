package v1

import (
	"go-socialapp/internal/pkg/util/idgenerate"
	"gorm.io/gorm"
	"time"
)

type ObjectMetaNotIncrement struct {
	ID uint64 `json:"id,omitempty" gorm:"primary_key;column:id"`

	// CreatedAt is a timestamp representing the server time when this object was
	// created. It is not guaranteed to be set in happens-before order across separate operations.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	//
	// Populated by the system.
	// Read-only.
	// Null for lists.
	CreatedAt time.Time `json:"createDate,omitempty" gorm:"column:create_date"`

	// UpdatedAt is a timestamp representing the server time when this object was updated.
	// Clients may not set this value. It is represented in RFC3339 form and is in UTC.
	//
	// Populated by the system.
	// Read-only.
	// Null for lists.
	UpdatedAt time.Time `json:"updateDate,omitempty" gorm:"column:update_date"`

	// DeletedAt is RFC 3339 date and time at which this resource will be deleted. This
	// field is set by the server when a graceful deletion is requested by the user, and is not
	// directly settable by a client.
	//
	// Populated by the system when a graceful deletion is requested.
	// Read-only.
	// DeletedAt gorm.DeletedAt `json:"-" gorm:"column:deletedAt;index:idx_deletedAt"`

	Creator int `json:"creator,omitempty" gorm:"column:creator"`

	Updater int `json:"updater,omitempty" gorm:"column:updater"`
}

// BeforeCreate run before create database record.
func (obj *ObjectMetaNotIncrement) BeforeCreate(tx *gorm.DB) error {
	id, err := idgenerate.GetIntID()
	obj.ID = id
	return err
}
