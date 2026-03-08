
    
    
	b.Amount = b.Amount * 100
	// --- 1. BASIC PATIENT VALIDATIONS ---
	if len(strings.TrimSpace(b.PatientName)) < 3 || b.PatientName == "" {
		utils.SendJSON(w, 400, "failed", "Name too short please enter a valid name", nil)
		return
	}
	if !strings.Contains(b.PatientEmail, "@") || len(strings.TrimSpace(b.PatientEmail)) < 5 {
		utils.SendJSON(w, 400, "failed", "Invalid email", nil)
		return
	}
	if strings.TrimSpace(b.PatientAddress) == "" {
		utils.SendJSON(w, 400, "failed", "Address field is required", nil)
		return
	}
	// gmail format validation
	if !strings.HasSuffix(b.PatientEmail, "@gmail.com") {
		utils.SendJSON(w, 400, "failed", "Only Gmail addresses are accepted like example@gmail.com", nil)
		return
	}
	gender := strings.ToLower(strings.TrimSpace(b.PatientGender))
	if gender != "male" && gender != "female" && gender != "other" {
		utils.SendJSON(w, 400, "failed", "Invalid Gender", nil)
		return
	}
	if b.PatientAge <= 0 || b.PatientAge > 120 || b.PatientAge == 0 {
		utils.SendJSON(w, 400, "failed", "Invalid Age", nil)
		return
	}

	// --- Simple & Solid Mobile Validation ---
	mobile := strings.TrimSpace(b.AccountHolderMobile)

	// 1. Length check
	if len(mobile) != 10 {
		utils.SendJSON(w, 400, "failed", "Mobile must be 10 digits", nil)
		return
	}

	// 2. Pure Numbers-only check (No +, -, or special chars)
	for _, char := range mobile {
		if !unicode.IsDigit(char) {
			utils.SendJSON(w, 400, "failed", "Mobile must contain only numbers (no special characters)", nil)
			return
		}
	}
	b.AccountHolderMobile = mobile

	// 1. Parsing the Appointment Date
	// Format: "2026-02-05" (YYYY-MM-DD)
	appDate, err := time.Parse("2006-01-02", b.AppointmentDate)
	if err != nil {
		utils.SendJSON(w, 400, "failed", "Invalid Date Format. Use YYYY-MM-DD", nil)
		return
	}

	// 2. Today's Date (without time component)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 3. Past Date Check
	if appDate.Before(today) {
		utils.SendJSON(w, 400, "failed", "Past dates are not allowed for booking", nil)
		return
	}

	// 4. Future 14-Day Limit Check
	maxDate := today.AddDate(0, 0, 14)
	if appDate.After(maxDate) {
		utils.SendJSON(w, 400, "failed", "Booking only allowed for the next 14 days", nil)
		return
	}

	// --- 2. STRICT LOCATION VALIDATION (India, Telangana, Hyderabad) ---
	// Normalize strings to ignore case and spaces
	country := strings.ToLower(strings.TrimSpace(b.PatientCountry))
	state := strings.ToLower(strings.TrimSpace(b.PatientState))
	city := strings.ToLower(strings.TrimSpace(b.PatientCity))

	if country != "india" {
		utils.SendJSON(w, 400, "failed", "Service only available in India", nil)
		return
	}
	if state != "telangana" && state != "ts" {
		utils.SendJSON(w, 400, "failed", "Service only available in Telangana", nil)
		return
	}
	if city != "hyderabad" && city != "hyd" {
		utils.SendJSON(w, 400, "failed", "Service only available in Hyderabad city", nil)
		return
	}

	// --- 3. TIME BOUNDARY VALIDATION (9 AM to 9 PM) ---
	b.TimeSlot = strings.ToUpper(strings.TrimSpace(b.TimeSlot))
	layout := "03:04 PM"
	reqTime, err := time.Parse(layout, b.TimeSlot)
	if err != nil {
		utils.SendJSON(w, 400, "failed", "Use format 09:00 AM", nil)
		return
	}
	hStart, _ := time.Parse(layout, "09:00 AM")
	hEnd, _ := time.Parse(layout, "09:00 PM")

	if reqTime.Before(hStart) || reqTime.After(hEnd) {
		utils.SendJSON(w, 400, "failed", "Hospital closed. Choose 9 AM - 9 PM", nil)
		return
	}

	// --- 4. GEOLOCATION & DISTANCE CHECK ---
	fullAddr := fmt.Sprintf("%s, %s, %s, %s", b.PatientAddress, b.PatientCity, b.PatientState, b.PatientPincode)
	lat, lon, actualPin, err := utils.GetCoordinates(fullAddr)
	if err != nil {
		utils.SendJSON(w, 400, "failed", "GPS Location not found", nil)
		return
	}

	if strings.ReplaceAll(actualPin, " ", "") != b.PatientPincode {
		utils.SendJSON(w, 400, "failed", "Pincode mismatch with current location", nil)
		return
	}

	b.Latitude, b.Longitude = lat, lon
	dist := utils.CalculateDistance(17.4126, 78.4735, lat, lon)
	if dist > 50.0 {
		utils.SendJSON(w, 403, "failed", "Service area is within 50km only", nil)
		return
	}

	// --- 5. DB TRANSACTIONS ---
	err = config.DB.Transaction(func(tx *gorm.DB) error {
		// One booking per mobile per day
		var dailyCheck models.Booking
		if err := tx.Where("account_holder_mobile = ? AND appointment_date = ? AND status != ?",
			b.AccountHolderMobile, b.AppointmentDate, "cancelled").First(&dailyCheck).Error; err == nil {
			return fmt.Errorf("YOU_ALREADY_HAVE_A_BOOKING_FOR_THIS_DATE")
		}

		// Slot availability
		var existing models.Booking
		if err := tx.Where("doctor_id = ? AND hospital_id = ? AND appointment_date = ? AND time_slot = ? AND status IN (?)",
			b.DoctorID, b.HospitalID, b.AppointmentDate, b.TimeSlot, []string{"confirmed", "pending"}).First(&existing).Error; err == nil {
			return fmt.Errorf("SLOT_ALREADY_BOOKED")
		}

		b.AppointmentID = fmt.Sprintf("APP-%d", time.Now().Unix())
		b.Status = "pending"
		if strings.ToLower(b.PaymentMethod) == "offline" {
			b.Status = "confirmed"
		}

		return tx.Create(&b).Error
	})
	if err != nil {
		utils.SendJSON(w, 400, "failed", err.Error(), nil)
		return
	}
