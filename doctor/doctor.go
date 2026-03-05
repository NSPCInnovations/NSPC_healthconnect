package doctor

import (
	"time"

	"gorm.io/gorm"
)

type Doctorregister struct {
	ID uint `gorm:"primaryKey" json:"id"`

	// Unique Identifiers
	DoctorRegNo string `gorm:"uniqueIndex;type:varchar(30);not null" json:"doctor_reg_no"`
	UserID      string `gorm:"uniqueIndex;type:varchar(30);not null" json:"user_id"`

	// Personal Details
	Prefix     string `gorm:"type:varchar(10)" json:"prefix"`
	FirstName  string `gorm:"type:varchar(50);not null" json:"first_name"`
	MiddleName string `gorm:"type:varchar(50)" json:"middle_name"`
	LastName   string `gorm:"type:varchar(50);not null" json:"last_name"`
	Gender     string `gorm:"type:varchar(10)" json:"gender"`

	Email           string `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Mobile          string `gorm:"uniqueIndex;type:varchar(15);not null" json:"mobile"`
	AlternateMobile string `gorm:"type:varchar(15)" json:"alternate_mobile"`

	// Professional Details
	Specialization  string  `gorm:"type:varchar(100);not null" json:"specialization"`
	Description     string  `gorm:"type:text" json:"description"`
	Experience      int     `json:"experience"`
	ConsultationFee float64 `gorm:"type:decimal(10,2)" json:"consultation_fee"`
	Currency        string  `gorm:"type:varchar(10);default:'INR'" json:"currency"`
	SelfPractice    bool    `gorm:"default:false" json:"self_practice"`

	// Location
	Location string `gorm:"type:varchar(150)" json:"location"`
	Address  string `gorm:"type:text" json:"address"`

	// Status & Verification
	RegistrationSource string `gorm:"type:varchar(20)" json:"registration_source"`
	VerificationStatus string `gorm:"type:varchar(20);default:'PENDING'" json:"verification_status"`
	IsActive           bool   `gorm:"default:false" json:"is_active"`

	// Rating (calculated field)
	AverageRating float64 `gorm:"type:decimal(2,1);default:0" json:"average_rating"`

	// Relationships
	Documents     []DoctorDocument     `gorm:"foreignKey:DoctorID" json:"documents,omitempty"`
	Availability  []DoctorAvailability `gorm:"foreignKey:DoctorID" json:"availability,omitempty"`
	Verifications []DoctorVerification `gorm:"foreignKey:DoctorID" json:"verifications,omitempty"`
	Hospitals     []DoctorHospitalMap  `gorm:"foreignKey:DoctorID" json:"hospitals,omitempty"`
	Ratings       []DoctorRating       `gorm:"foreignKey:DoctorID" json:"ratings,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
