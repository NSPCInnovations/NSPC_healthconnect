package models

import (
	"time"

	"gorm.io/gorm"
)

type Booking struct {
	ID                  uint           `gorm:"primaryKey" json:"id"`
	AppointmentID       string         `gorm:"uniqueIndex;type:varchar(50)" json:"appointment_id"`
	PatientName         string         `gorm:"type:varchar(100)" json:"patient_name"`
	AccountHolderMobile string         `gorm:"uniqueIndex:idx_mobile_date_slot;type:varchar(15)" json:"account_holder_mobile"`
	PatientAge          int            `json:"patient_age"`
	PatientGender       string         `gorm:"type:varchar(10)" json:"patient_gender"`
	PatientEmail        string         `gorm:"type:varchar(100)" json:"patient_email"`
	PatientAddress      string         `gorm:"type:text" json:"patient_address"`
	PatientCountry      string         `json:"patient_country"`
	PatientState        string         `json:"patient_state"`
	PatientCity         string         `json:"patient_city"`
	PatientPincode      string         `gorm:"type:varchar(10)" json:"patient_pincode"`
	Latitude            float64        `json:"latitude"`
	Longitude           float64        `json:"longitude"`
	DoctorID            string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;uniqueIndex:idx_mobile_date_slot;type:varchar(50)" json:"doctor_id"`
	HospitalID          string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;type:varchar(50)" json:"hospital_id"`
	AppointmentDate     string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;uniqueIndex:idx_mobile_date_slot;type:varchar(20)" json:"appointment_date"`
	TimeSlot            string         `gorm:"uniqueIndex:idx_doc_hosp_date_slot;uniqueIndex:idx_mobile_date_slot;type:varchar(20)" json:"time_slot"`
	PaymentMethod       string         `gorm:"type:varchar(20)" json:"payment_method"`
	OrderID             string         `gorm:"column:order_id;default:null" json:"order_id"`
	PaymentID           string         `gorm:"type:varchar(100)" json:"payment_id"`
	Amount              int            `gorm:"column:amount" json:"amount"`
	PaymentLink         string         `gorm:"type:text" json:"payment_link"`
	Status              string         `gorm:"type:varchar(20);default:'pending'" json:"status"`
	CreatedAt           time.Time      `json:"created_at"`
	UpdatedAt           time.Time      `json:"updated_at"`
	DeletedAt           gorm.DeletedAt `gorm:"index" json:"-"`
}
