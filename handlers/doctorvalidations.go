package handlers

import (
	"NSPC_healthconnect/doctor"
	"regexp"
)

// ValidateDoctor validates doctor registration fields
func ValidateDoctor(d doctor.Doctor) string {

	// Doctor Registration Number validation
	if d.DoctorRegNo == "" {
		return "Doctor registration number is required"
	}

	// UserID validation
	if d.UserID == "" {
		return "User ID is required"
	}

	// Name validation (alphabets only)
	nameRegex := regexp.MustCompile(`^[A-Za-z ]+$`)

	if d.FirstName == "" {
		return "First name is required"
	}

	if !nameRegex.MatchString(d.FirstName) {
		return "First name should contain only alphabets"
	}

	if d.MiddleName != "" && !nameRegex.MatchString(d.MiddleName) {
		return "Middle name should contain only alphabets"
	}

	if d.LastName == "" {
		return "Last name is required"
	}

	if !nameRegex.MatchString(d.LastName) {
		return "Last name should contain only alphabets"
	}

	// Gender validation
	if d.Gender == "" {
		return "Gender is required"
	}

	// Mobile validation
	mobileRegex := regexp.MustCompile(`^[0-9]{10}$`)

	if d.Mobile == "" {
		return "Mobile number is required"
	}

	if !mobileRegex.MatchString(d.Mobile) {
		return "Mobile number must be exactly 10 digits"
	}

	// Alternate mobile validation
	if d.AlternateMobile != "" && !mobileRegex.MatchString(d.AlternateMobile) {
		return "Alternate mobile number must be 10 digits"
	}

	// Email validation (gmail only)
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

	// Experience validation
	if d.Experience < 0 {
		return "Experience cannot be negative"
	}

	// Consultation fee validation
	if d.ConsultationFee <= 0 {
		return "Consultation fee must be greater than zero"
	}

	// Location validation
	locationRegex := regexp.MustCompile(`^[A-Za-z ]+$`)

	if d.Location == "" {
		return "Location is required"
	}

	if !locationRegex.MatchString(d.Location) {
		return "Location should contain only alphabets"
	}

	// Address validation
	if d.Address == "" {
		return "Address is required"
	}

	// Registration Source validation
	if d.RegistrationSource == "" {
		return "Registration source is required"
	}

	return ""
}
