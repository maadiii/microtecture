package models

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// User is user model
type User struct {
	Id           uuid.UUID  `gorm:"type:uuid;primary_key" json:"id,omitempty"`
	CreatedAt    time.Time  `json:"createdAt,omitempty"`
	UpdatedAt    time.Time  `json:"createdAt,omitempty"`
	DeletedAt    *time.Time `sql:"index" json:"deletedId,omitempty"`
	Password     []byte     `gorm:"not null" json:"password,omitempty"`
	MobileNumber string     `gorm:"type:varchar(11);unique_index" json:"mobileNumber,omitempty"`
	FirstName    string     `gorm:"type:varchar(64);index" json:"firstName,omitempty"`
	LastName     string     `gorm:"type:varchar(64);index" json:"lastName,omitempty"`

	GroupID uuid.UUID `gorm:"type:uuid;not null" json:"groupId,omitempty"`
	Group   Group
}

// Group holder of users group
type Group struct {
	Id          uuid.UUID  `gorm:"type:uuid;priamry_key" json:"id,omitempty"`
	DeletedAt   *time.Time `sql:"index" json:"deletedAt,omitempty"`
	Name        string     `gorm:"type:varchar(64);not null;unique;" json:"name,omitempty"`
	Description string     `gorm:"type:varchar(256);not null" json:"description,omitempty"`

	Roles []Role `gorm:"many2many:groups_roles" json:"roles,omitempty"`
	Users []User
}

// Roles is roles of users
type Role struct {
	Id        uuid.UUID  `gorm:"type:uuid;priamry_key" json:"id,omitempty"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt,omitempty"`
	FaName    string     `gorm:"not null; type:varchar(64);unique" json:"faName,omitempty"`
	EnName    string     `gorm:"not null; type:varchar(64);unique" json:"enName,omitempty"`

	Groups []Group `gorm:"many2many:groups_roles" json:"groups,omitempty"`
}

func (self *User) BeforeCreate(scope *gorm.Scope) (err error) {
	hashed, err := bcrypt.GenerateFromPassword(self.Password, bcrypt.DefaultCost)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to password hashing: %v", err))
	}

	err = scope.SetColumn("Password", hashed)
	if err != nil {
		return errors.New(err.Error())
	}

	return scope.SetColumn("ID", uuid.New())
}
