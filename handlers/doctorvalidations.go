package handlers

import (
	"NSPC_healthconnect/doctor"
	"regexp"
)

// ValidateDoctor validates doctor registration fields
func ValidateDoctor(d doctor.Doctor) string {

	// Name validation (alphabets only)
	nameRegex := regexp.MustCompile(`^[A-Za-z ]+$`)

	if d.FirstName == "" {
		return "First name is required"
	}

	if !nameRegex.MatchString(d.FirstName) {
		return "First name should contain only alphabets"
	}

	if d.LastName != "" {
		if !nameRegex.MatchString(d.LastName) {
			return "Last name should contain only alphabets"
		}
	}

	// Mobile validation
	mobileRegex := regexp.MustCompile(`^[0-9]{10}$`)

	if d.Mobile == "" {
		return "Mobile number is required"
	}

	if !mobileRegex.MatchString(d.Mobile) {
		return "Mobile number must be exactly 10 digits"
	}

	// Email validation (only gmail allowed)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@gmail\.com$`)

	if d.Email == "" {
		return "Email is required"
	}

	if !emailRegex.MatchString(d.Email) {
		return "Email must be a valid Gmail address (example@gmail.com)"
	}

	// Specialization validation
	specialRegex := regexp.MustCompile(`^[A-Za-z ]+$`)

	if d.Specialization == "" {
		return "Specialization is required"
	}

	if !specialRegex.MatchString(d.Specialization) {
		return "Specialization should contain only alphabets"
	}

	// Consultation fee validation
	if d.ConsultationFee <= 0 {
		return "Consultation fee must be a positive number"
	}

	// Location validation
	locationRegex := regexp.MustCompile(`^[A-Za-z ]+$`)

	if d.Location == "" {
		return "Location is required"
	}

	if !locationRegex.MatchString(d.Location) {
		return "Location should contain only alphabets"
	}

	return ""
}
