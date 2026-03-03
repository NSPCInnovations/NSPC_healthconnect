package appointment

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	AppointmentID   string         `gorm:"uniqueIndex;type:varchar(50)" json:"appointment_id"`
	DoctorID        string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;uniqueIndex:idx_mobile_date_slot;type:varchar(50)" json:"doctor_id"`
	HospitalID      string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;type:varchar(50)" json:"hospital_id"`
	Email           string         `gorm:"type:varchar(100)" json:"email"`
	AppointmentDate string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;uniqueIndex:idx_mobile_date_slot;type:varchar(20)" json:"appointment_date"`
	TimeSlot        string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;uniqueIndex:idx_mobile_date_slot;type:varchar(20)" json:"time_slot"`
	PaymentMethod   string         `gorm:"type:varchar(20)" json:"payment_method"`
	OrderID         string         `gorm:"column:order_id;default:null" json:"order_id"`
	PaymentID       string         `gorm:"type:varchar(100)" json:"payment_id"`
	Amount          int            `gorm:"column:amount" json:"amount"`
	PaymentLink     string         `gorm:"type:text" json:"payment_link"`
	Status          string         `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}
