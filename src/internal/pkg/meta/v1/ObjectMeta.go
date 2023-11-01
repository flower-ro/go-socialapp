package v1

import (
	"time"

	"github.com/marmotedu/component-base/pkg/json"
)

type Extend map[string]interface{}

// String returns the string format of Extend.
func (ext Extend) String() string {
	data, _ := json.Marshal(ext)
	return string(data)
}

// Merge merge extend fields from extendShadow.
func (ext Extend) Merge(extendShadow string) Extend {
	var extend Extend

	// always trust the extendShadow in the database
	_ = json.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}

	return ext
}

// ObjectMeta is metadata that all persisted resources must have, which includes all objects
// ObjectMeta is also used by gorm.
type ObjectMeta struct {
	// ID is the unique in time and space value for this object. It is typically generated by
	// the storage on successful creation of a resource and is not allowed to change on PUT
	// operations.
	//
	// Populated by the system.
	// Read-only.
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`

	// InstanceID defines a string type resource identifier,
	// use prefixed to distinguish resource types, easy to remember, Url-friendly.
	//	InstanceID string `json:"instanceID,omitempty" gorm:"unique;column:instance_id;type:varchar(32);not null"`

	// Name defines the space within each name must be unique.
	// Not all objects are required to be scoped to a username - the value of this field for
	// those objects will be empty.
	//
	// Must be a DNS_LABEL.
	// Cannot be updated.
	// Username string `json:"username,omitempty" gorm:"column:username" validate:"omitempty"`

	// Required: true
	// Name must be unique. Is required when creating resources.
	// Name is primarily intended for creation idempotence and configuration
	// definition.
	// It will be generated automated only if Name is not specified.
	// Cannot be updated.
	//	Name string `json:"name,omitempty" gorm:"column:name;type:varchar(64);not null" validate:"name"`

	// Extend store the fields that need to be added, but do not want to add a new table column, will not be stored in db.
	//Extend Extend `json:"extend,omitempty" gorm:"-" validate:"omitempty"`
	//
	//// ExtendShadow is the shadow of Extend. DO NOT modify directly.
	//ExtendShadow string `json:"-" gorm:"column:extend_shadow" validate:"omitempty"`

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
//func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) error {
//	obj.ExtendShadow = obj.Extend.String()
//
//	return nil
//}
//
//// BeforeUpdate run before update database record.
//func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) error {
//	obj.ExtendShadow = obj.Extend.String()
//
//	return nil
//}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
//func (obj *ObjectMeta) AfterFind(tx *gorm.DB) error {
//	if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
//		return err
//	}
//
//	return nil
//}
