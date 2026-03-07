package doctor

import "time"

type Doctor struct {
	ID uint `gorm:"primaryKey"`

	DoctorRegNo string `gorm:"unique" json:"doctor_reg_no"`
	UserID      string `json:"user_id"`

	Prefix     string `json:"prefix"`
	FirstName  string `json:"first_name"`
	MiddleName string `json:"middle_name"`
	LastName   string `json:"last_name"`

	Email  string `gorm:"unique" json:"email"`
	Mobile string `gorm:"unique" json:"mobile"`

	Specialization string `json:"specialization"`
	Experience     int    `json:"experience"`
	Description    string `json:"description"`

	ConsultationFee float64 `json:"consultation_fee"`
	Currency        string  `json:"currency"`

	Location string `json:"location"`

	VerificationStatus string `json:"verification_status"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

type DoctorVerification struct {
	ID uint `gorm:"primaryKey"`

	DoctorID uint

	VerificationStatus string
	VerifiedBy         string
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