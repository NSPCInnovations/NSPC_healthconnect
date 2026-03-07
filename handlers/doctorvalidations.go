package handlers

import "NSPC_healthconnect/doctor"

func ValidateDoctor(d doctor.Doctor) string {

	if d.FirstName == "" {
		return "First name required"
	}

	if d.Mobile == "" {
		return "Mobile required"
	}

	if d.Email == "" {
		return "Email required"
	}

	if d.Specialization == "" {
		return "Specialization required"
	}

	return ""
}