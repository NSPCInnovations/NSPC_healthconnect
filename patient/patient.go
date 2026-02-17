package NSPC_HEALTHCONNECT

import "time"

type Patient struct {
	PatientID    uint      `gorm:"primaryKey;autoIncrement" json:"patient_id"`
	FirstName    string    `gorm:"column:first_name;not null" json:"first_name"`
	LastName     string    `gorm:"column:last_name;not null" json:"last_name"`
	MobileNumber string    `gorm:"column:mobile_number;not null;unique" json:"mobile_number"`
	Email        string    `gorm:"column:email" json:"email"`
	Age          int       `gorm:"column:age;not null" json:"age"`
	Gender       string    `gorm:"column:gender;not null" json:"gender"`
	CreatedAt    time.Time `gorm:"autoCreateTime" json:"created_at"`
}
