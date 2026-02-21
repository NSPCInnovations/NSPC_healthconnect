package doctor

import "time"

type Doctor struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DoctorRegNo string `gorm:"unique;not null" json:"doctor_reg_no"`
	UserID      string `gorm:"unique;not null" json:"user_id"`

	FirstName string `gorm:"not null" json:"first_name"`
	LastName  string `gorm:"not null" json:"last_name"`

	Email  string `gorm:"unique;not null" json:"email"`
	Mobile string `gorm:"unique;not null" json:"mobile"`

	Specialization string `gorm:"not null" json:"specialization"`

	VerificationStatus string `gorm:"default:PENDING" json:"verification_status"`
	IsActive           bool   `gorm:"default:false" json:"is_active"`

	// Optional fields
	Prefix           string  `json:"prefix"`
	MiddleName       string  `json:"middle_name"`
	Gender           string  `json:"gender"`
	AlternateMobile  string  `json:"alternate_mobile"`
	Description      string  `json:"description"`
	Experience        int     `json:"experience"`
	ConsultationFee  float64 `json:"consultation_fee"`
	Currency          string  `gorm:"default:INR" json:"currency"`
	SelfPractice      bool    `json:"self_practice"`
	Location          string  `json:"location"`
	Address           string  `json:"address"`
	RegistrationSource string `json:"registration_source"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
