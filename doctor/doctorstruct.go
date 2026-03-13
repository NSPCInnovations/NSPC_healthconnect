package doctor

import (
	"time"

	"gorm.io/gorm"
)

type Doctor struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DoctorRegNo string `gorm:"uniqueIndex;type:varchar(30);not null" json:"doctor_reg_no"`
	UserID      string `gorm:"uniqueIndex;type:varchar(30);not null" json:"user_id"`

	Prefix     string `gorm:"type:varchar(10)" json:"prefix"`
	FirstName  string `gorm:"type:varchar(50);not null" json:"first_name"`
	MiddleName string `gorm:"type:varchar(50)" json:"middle_name"`
	LastName   string `gorm:"type:varchar(50);not null" json:"last_name"`
	Gender     string `gorm:"type:varchar(10)" json:"gender"`

	Email           string `gorm:"uniqueIndex;type:varchar(100);not null" json:"email"`
	Mobile          string `gorm:"uniqueIndex;type:varchar(15);not null" json:"mobile"`
	AlternateMobile string `gorm:"type:varchar(15)" json:"alternate_mobile"`

	Specialization  string  `gorm:"type:varchar(100);not null" json:"specialization"`
	Description     string  `gorm:"type:text" json:"description"`
	Experience      int     `json:"experience"`
	ConsultationFee float64 `gorm:"type:decimal(10,2)" json:"consultation_fee"`
	Currency        string  `gorm:"type:varchar(10);default:'INR'" json:"currency"`

	SelfPractice bool `gorm:"default:false" json:"self_practice"`

	Location string `gorm:"type:varchar(150)" json:"location"`
	Address  string `gorm:"type:text" json:"address"`

	RegistrationSource string `gorm:"type:varchar(20)" json:"registration_source"`
	VerificationStatus string `gorm:"type:varchar(20);default:'PENDING'" json:"verification_status"`
	IsActive           bool   `gorm:"default:false" json:"is_active"`

	AverageRating float64 `gorm:"type:decimal(2,1);default:0" json:"average_rating"`

	// Relationships
	Documents     []DoctorDocument     `gorm:"foreignKey:DoctorID" json:"documents"`
	Availability  []DoctorAvailability `gorm:"foreignKey:DoctorID" json:"availability"`
	Verifications []DoctorVerification `gorm:"foreignKey:DoctorID" json:"verifications"`
	Hospitals     []DoctorHospitalMap  `gorm:"foreignKey:DoctorID" json:"hospitals"`
	Ratings       []DoctorRating       `gorm:"foreignKey:DoctorID" json:"ratings"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type DoctorDocument struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DoctorID uint `json:"doctor_id"`

	LicenseNumber string    `gorm:"type:varchar(50)" json:"license_number"`
	MedicalBoard  string    `gorm:"type:varchar(100)" json:"medical_board"`
	LicenseExpiry time.Time `json:"license_expiry"`

	DocumentURL string `gorm:"type:text" json:"document_url"`

	IsValid bool `gorm:"default:true" json:"is_valid"`

	CreatedAt time.Time `json:"created_at"`
}

type DoctorVerification struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DoctorID uint `json:"doctor_id"`

	VerificationStatus string `gorm:"type:varchar(20)" json:"verification_status"`

	VerifiedBy string `gorm:"type:varchar(50)" json:"verified_by"`

	VerifiedAt time.Time `json:"verified_at"`

	ApprovalDescription string `gorm:"type:text" json:"approval_description"`

	CreatedAt time.Time `json:"created_at"`
}

type DoctorAvailability struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DoctorID uint `json:"doctor_id"`

	DayOfWeek string `gorm:"type:varchar(15)" json:"day_of_week"`

	TimeSlot string `gorm:"type:varchar(50)" json:"time_slot"`

	SlotType string `gorm:"type:varchar(20)" json:"slot_type"`

	PatientDurationMinutes int `gorm:"default:10" json:"patient_duration_minutes"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
}

type DoctorHospitalMap struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DoctorID uint `json:"doctor_id"`

	HospitalID uint `json:"hospital_id"`

	Specialization string `gorm:"type:varchar(100)" json:"specialization"`

	HospitalLocation string `gorm:"type:varchar(150)" json:"hospital_location"`

	GoogleMapLink string `gorm:"type:text" json:"google_map_link"`

	AvailabilitySummary string `gorm:"type:varchar(100)" json:"availability_summary"`

	IsActive bool `gorm:"default:true" json:"is_active"`

	CreatedAt time.Time `json:"created_at"`
}

type DoctorRating struct {
	ID uint `gorm:"primaryKey" json:"id"`

	DoctorID uint `json:"doctor_id"`

	PatientID string `gorm:"type:varchar(50)" json:"patient_id"`

	Rating float64 `gorm:"type:decimal(2,1)" json:"rating"`

	Comments string `gorm:"type:text" json:"comments"`

	CreatedAt time.Time `json:"created_at"`
}
