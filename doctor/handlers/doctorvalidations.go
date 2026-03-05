package doctor

import (
	"nspc_healthcare/config"

	"github.com/gin-gonic/gin"
)

func RegisterDoctor(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", 405)
		return
	}

	var d models.Doctor
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		utils.SendJSON(w, 400, "failed", "Invalid JSON", nil)
		return
	}

	// 🔹 VALIDATIONS

	if len(strings.TrimSpace(d.FirstName)) < 3 {
		utils.SendJSON(w, 400, "failed", "First name too short", nil)
		return
	}

	if !strings.Contains(d.Email, "@") {
		utils.SendJSON(w, 400, "failed", "Invalid email format", nil)
		return
	}

	if len(d.Mobile) != 10 {
		utils.SendJSON(w, 400, "failed", "Mobile must be 10 digits", nil)
		return
	}

	for _, ch := range d.Mobile {
		if !unicode.IsDigit(ch) {
			utils.SendJSON(w, 400, "failed", "Mobile must contain only numbers", nil)
			return
		}
	}

	if d.Specialization == "" {
		utils.SendJSON(w, 400, "failed", "Specialization required", nil)
		return
	}

	// 🔹 DUPLICATE CHECK
	var existing models.Doctor
	if err := config.DB.Where("email = ? OR mobile = ?", d.Email, d.Mobile).
		First(&existing).Error; err == nil {
		utils.SendJSON(w, 400, "failed", "Doctor already exists", nil)
		return
	}

	// 🔹 AUTO GENERATED FIELDS
	d.DoctorRegNo = fmt.Sprintf("DOC-%d", time.Now().Unix())
	d.VerificationStatus = "PENDING"
	d.IsActive = false

	// 🔹 TRANSACTION
	err := config.DB.Transaction(func(tx *gorm.DB) error {
		return tx.Create(&d).Error
	})

	if err != nil {
		utils.SendJSON(w, 500, "failed", "Registration Failed", nil)
		return
	}

	utils.SendJSON(w, 201, "success", "Doctor Registered. Awaiting Verification", d)
}

func VerifyDoctor(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPut {
		http.Error(w, "Only PUT allowed", 405)
		return
	}

	doctorID := r.URL.Query().Get("doctor_id")
	if doctorID == "" {
		utils.SendJSON(w, 400, "failed", "Doctor ID required", nil)
		return
	}

	var doctor models.Doctor
	if err := config.DB.First(&doctor, doctorID).Error; err != nil {
		utils.SendJSON(w, 404, "failed", "Doctor not found", nil)
		return
	}

	if doctor.VerificationStatus == "APPROVED" {
		utils.SendJSON(w, 400, "failed", "Doctor already verified", nil)
		return
	}

	err := config.DB.Transaction(func(tx *gorm.DB) error {

		doctor.VerificationStatus = "APPROVED"
		doctor.IsActive = true

		return tx.Save(&doctor).Error
	})

	if err != nil {
		utils.SendJSON(w, 500, "failed", "Verification failed", nil)
		return
	}

	utils.SendJSON(w, 200, "success", "Doctor Verified & Activated", doctor)
}

func AddDoctorAvailability(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", 405)
		return
	}

	var a models.DoctorAvailability
	if err := json.NewDecoder(r.Body).Decode(&a); err != nil {
		utils.SendJSON(w, 400, "failed", "Invalid JSON", nil)
		return
	}

	if a.DoctorID == 0 {
		utils.SendJSON(w, 400, "failed", "Doctor ID required", nil)
		return
	}

	if a.Day == "" || a.TimeSlot == "" {
		utils.SendJSON(w, 400, "failed", "Day and Time Slot required", nil)
		return
	}

	// Duplicate slot check
	var existing models.DoctorAvailability
	err := config.DB.Where(
		"doctor_id = ? AND day = ? AND time_slot = ?",
		a.DoctorID,
		a.Day,
		a.TimeSlot,
	).First(&existing).Error

	if err == nil {
		utils.SendJSON(w, 400, "failed", "Slot already exists", nil)
		return
	}

	if err := config.DB.Create(&a).Error; err != nil {
		utils.SendJSON(w, 500, "failed", "Failed to add availability", nil)
		return
	}

	utils.SendJSON(w, 201, "success", "Availability Added Successfully", a)
}

func MapDoctorToHospital(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", 405)
		return
	}

	var m models.DoctorHospitalMap
	if err := json.NewDecoder(r.Body).Decode(&m); err != nil {
		utils.SendJSON(w, 400, "failed", "Invalid JSON", nil)
		return
	}

	if m.DoctorID == 0 || m.HospitalID == 0 {
		utils.SendJSON(w, 400, "failed", "DoctorID and HospitalID required", nil)
		return
	}

	var existing models.DoctorHospitalMap
	err := config.DB.Where(
		"doctor_id = ? AND hospital_id = ?",
		m.DoctorID,
		m.HospitalID,
	).First(&existing).Error

	if err == nil {
		utils.SendJSON(w, 400, "failed", "Doctor already mapped", nil)
		return
	}

	if err := config.DB.Create(&m).Error; err != nil {
		utils.SendJSON(w, 500, "failed", "Mapping failed", nil)
		return
	}

	utils.SendJSON(w, 201, "success", "Doctor mapped successfully", m)
}

func ListActiveDoctors(w http.ResponseWriter, r *http.Request) {

	var doctors []models.Doctor

	err := config.DB.
		Where("verification_status = ? AND is_active = ?", "APPROVED", true).
		Find(&doctors).Error

	if err != nil {
		utils.SendJSON(w, 500, "failed", "Failed to fetch doctors", nil)
		return
	}

	utils.SendJSON(w, 200, "success", "Doctors fetched successfully", doctors)
}