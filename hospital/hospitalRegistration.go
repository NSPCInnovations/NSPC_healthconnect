package hospital

import (
	"time"

	"gorm.io/gorm"
)

type Hospital struct {
	ID                    uint           `gorm:"primaryKey" json:"id"`
	HospitalName          string         `gorm:"column:hospital_name;not null" json:"hospital_name"`
	Address               string         `gorm:"column:address;not null" json:"address"`
	Type                  string         `gorm:"column:type" json:"type"`
	Contact               string         `gorm:"column:contact;size:15" json:"contact"`
	Status                string         `gorm:"column:status;default:Pending" json:"status"`
	PermissionCertificate string         `gorm:"column:permission_certificate" json:"permission_certificate"`
	CreatedAt             time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt             time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt             gorm.DeletedAt `gorm:"index" json:"-"`
}
