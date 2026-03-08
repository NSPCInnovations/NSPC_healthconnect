package hospital

import (
	"regexp"
)

func ValidateHospital(h Hospital) string {

	// Hospital Name
	if h.HospitalName == "" {
		return "Hospital name is required"
	}

	if len(h.HospitalName) < 3 {
		return "Hospital name must be at least 3 characters"
	}

	// Address
	if h.Address == "" {
		return "Address is required"
	}

	// Type
	if h.Type == "" {
		return "Hospital type is required"
	}

	if h.Type != "Private" && h.Type != "Government" {
		return "Hospital type must be Private or Government"
	}

	// Contact validation
	if h.Contact == "" {
		return "Contact number is required"
	}

	phoneRegex := regexp.MustCompile(`^[0-9]{10}$`)
	if !phoneRegex.MatchString(h.Contact) {
		return "Contact number must be 10 digits"
	}

	// Certificate
	if h.PermissionCertificate == "" {
		return "Permission certificate is required"
	}

	return ""
}
