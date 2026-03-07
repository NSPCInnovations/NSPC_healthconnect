package doctor

import "time"

type Doctor struct {
	ID uint `gorm:"primaryKey"`

	DoctorRegNo string `gorm:"unique;not null" json:"doctor_reg_no"`
	UserID      string `gorm:"unique;not null" json:"user_id"`

	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Email  string `gorm:"unique" json:"email"`
	Mobile string `gorm:"unique" json:"mobile"`

	Specialization  string  `json:"specialization"`
	ConsultationFee float64 `json:"consultation_fee"`

	Location string `json:"location"`

	VerificationStatus string `json:"verification_status"`

	CreatedAt time.Time
	UpdatedAt time.Time
}
type DoctorVerification struct {
	ID uint `gorm:"primaryKey"`

	DoctorID uint

	VerificationStatus  string
	VerifiedBy          string
	ApprovalDescription string

	CreatedAt time.Time
}

type DoctorAvailability struct {
	ID uint `gorm:"primaryKey"`

	DoctorID uint

	DayOfWeek string
	TimeSlot  string

	CreatedAt time.Time
}

type DoctorHospitalMapping struct {
	ID uint `gorm:"primaryKey"`

	DoctorID   uint
	HospitalID uint

	Specialization string
	Location       string

	CreatedAt time.Time
}
