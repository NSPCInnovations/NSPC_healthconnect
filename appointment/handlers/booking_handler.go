package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	//"unicode"
	"NSPC_HEALTHCONNECT/appointment"
	"NSPC_HEALTHCONNECT/appointment/helper"
	"NSPC_HEALTHCONNECT/appointment/services"
	"NSPC_HEALTHCONNECT/database"
	"log"

	"gorm.io/gorm"
)

func CreateBooking(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", 405)
		return
	}
	var b appointment.Booking
	if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
		helper.SendJSON(w, 400, "failed", "Invalid JSON", nil)
		return
	}

	if b.Amount <= 0 {
		b.Amount = 500
	}
	b.Amount = b.Amount * 100

	//email validation
	if b.Email != "" {
		if !strings.Contains(b.Email, "@") || !strings.Contains(b.Email, ".") {
			helper.SendJSON(w, 400, "failed", "Invalid Email Address", nil)
			return
		}
	}
	// gmail format validation
	if !strings.HasSuffix(b.Email, "@gmail.com") {
		helper.SendJSON(w, 400, "failed", "Only Gmail addresses are accepted like example@gmail.com", nil)
		return
	}
	// --- 1. DATE VALIDATION (No past dates, max 14 days in future) ---
	// 1. Parsing the Appointment Date
	// Format: "2026-02-05" (YYYY-MM-DD)
	appDate, err := time.Parse("2006-01-02", b.AppointmentDate)
	if err != nil {
		helper.SendJSON(w, 400, "failed", "Invalid Date Format. Use YYYY-MM-DD", nil)
		return
	}

	// 2. Today's Date (without time component)
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// 3. Past Date Check
	if appDate.Before(today) {
		helper.SendJSON(w, 400, "failed", "Past dates are not allowed for booking", nil)
		return
	}

	// 4. Future 14-Day Limit Check
	maxDate := today.AddDate(0, 0, 14)
	if appDate.After(maxDate) {
		helper.SendJSON(w, 400, "failed", "Booking only allowed for the next 14 days", nil)
		return
	}
	// --- 3. TIME BOUNDARY VALIDATION (9 AM to 9 PM) ---
	b.TimeSlot = strings.ToUpper(strings.TrimSpace(b.TimeSlot))
	layout := "03:04 PM"
	reqTime, err := time.Parse(layout, b.TimeSlot)
	if err != nil {
		helper.SendJSON(w, 400, "failed", "Use format 09:00 AM", nil)
		return
	}
	hStart, _ := time.Parse(layout, "09:00 AM")
	hEnd, _ := time.Parse(layout, "09:00 PM")

	if reqTime.Before(hStart) || reqTime.After(hEnd) {
		helper.SendJSON(w, 400, "failed", "Hospital closed. Choose 9 AM - 9 PM", nil)
		return
	}

	// --- 5. DB TRANSACTIONS ---
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		// One booking per mobile per day
		var dailyCheck appointment.Booking
		if err := tx.Where("appointment_date = ? AND status != ?",
			b.AppointmentDate, "cancelled").First(&dailyCheck).Error; err == nil {
			return fmt.Errorf("YOU_ALREADY_HAVE_A_BOOKING_FOR_THIS_DATE")
		}

		// Slot availability
		var existing appointment.Booking
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
		helper.SendJSON(w, 400, "failed", err.Error(), nil)
		return
	}
	// --- 4. ONLINE PAYMENT LOGIC ---
	if strings.ToLower(b.PaymentMethod) == "online" {
		// Ikkada check chey amount correct ga undo ledo
		fmt.Println("Final Amount for Razorpay Order:", b.Amount)

		// Ikkada b.Amount already 100 aipoindi (Step 1 lo chesam)
		oID, errOrder := services.CreateRazorpayOrder(int64(b.Amount), b.AppointmentID)
		if errOrder != nil {
			fmt.Println("RAZORPAY REAL ERROR:", errOrder)
			helper.SendJSON(w, 500, "failed", "Razorpay Order Failed", nil)
			return
		}
		pLink, _ := services.CreatePaymentLink(oID, float64(b.Amount)/100.0)

		database.DB.Model(&b).Updates(map[string]interface{}{
			"order_id":     oID,
			"payment_link": pLink,
		})
		response := map[string]interface{}{
			"appointment_id": b.AppointmentID,
			"amount":         b.Amount / 100,
			"payment_link":   pLink,
		}
		fmt.Println("-------------------------------------------")
		fmt.Println("GENERATED PAYMENT LINK:", pLink)
		fmt.Println("-------------------------------------------")
		helper.SendJSON(w, 201, "success", "Payment Initiated", response)
		return
	}
	// Case 2: Offline Payment (Direct Confirmation)
	if strings.ToLower(b.PaymentMethod) == "offline" && b.Status == "confirmed" {
		b.Amount = b.Amount / 100
		newAmount := b.Amount
		err := database.DB.Model(&appointment.Booking{}).Where("id = ?", b.ID).Updates(map[string]interface{}{
			"amount": newAmount,
			"status": "confirmed",
		}).Error

		if err != nil {
			fmt.Println("DB Update Error (Offline):", err)
		}
		// 2. Local struct update for PDF
		b.Amount = newAmount
		b.Status = "confirmed"
		pdfPath, _ := services.GenerateBookingPDF(b)
		if b.Email != "" {
			go services.SendEmailWithPDF(b.Email, b.DoctorID, b.HospitalID, pdfPath, b.AppointmentID)
		}
		log.Printf("SMS SENT to %s: Your Appointment %s please pay at hospital reception.\n", b.AppointmentID)

		helper.SendJSON(w, 201, "success", "Booking Successful", b)
		return
	}
}

// 2. GET AVAILABLE SLOTS FOR A DOCTOR & HOSPITAL ON A GIVEN DATE

func GetAvailableSlots(w http.ResponseWriter, r *http.Request) {
	// 1. Fetch Query Parameters
	date := r.URL.Query().Get("date")
	doctorID := r.URL.Query().Get("doctor_id")
	hospitalID := r.URL.Query().Get("hospital_id")

	if date == "" || doctorID == "" || hospitalID == "" {
		helper.SendJSON(w, http.StatusBadRequest, "failed", "Date, DoctorID, and HospitalID are required", nil)
		return
	}
	var req appointment.Booking
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		helper.SendJSON(w, http.StatusBadRequest, "failed", "Invalid request body", nil)
		return
	}
	now := time.Now()
	layout := "2006-01-02 03:04 PM"
	appointmentFullTime, _ := time.Parse(layout, req.AppointmentDate+" "+req.TimeSlot)
	if appointmentFullTime.Before(now) {
		helper.SendJSON(w, http.StatusBadRequest, "failed", "Time out! This slot has already passed.", nil)
		return
	}
	// --- NEW: PAST DATE VALIDATION ---
	now = time.Now()
	currentDate := now.Format("2006-01-02")
	if date < currentDate {
		helper.SendJSON(w, http.StatusBadRequest, "failed", "Past dates are not accepted. Please select a valid date.", nil)
		return
	}

	// 2. Define Operating Hours
	var allSlots []string
	layout = "03:04 PM"
	startTime, _ := time.Parse(layout, "09:00 AM")
	endTime, _ := time.Parse(layout, "09:00 PM")

	// --- NEW: CURRENT TIME FOR TODAY'S FILTER ---
	currentTimeStr := time.Now().Format("03:04 PM")
	currentTime, _ := time.Parse(layout, currentTimeStr)

	for t := startTime; t.Before(endTime) || t.Equal(endTime); t = t.Add(30 * time.Minute) {
		slotStr := t.Format(layout)

		// Skip Lunch Break
		if slotStr == "01:00 PM" || slotStr == "01:30 PM" {
			continue
		}

		// --- NEW: TIME FILTER LOGIC ---
		if date == currentDate {

			if t.After(currentTime) {
				allSlots = append(allSlots, slotStr)
			}
		} else {
			allSlots = append(allSlots, slotStr)
		}
	}

	// 3. Fetch Booked Slots (Nee Old DB Logic as it is)
	var bookedBookings []appointment.Booking
	err := database.DB.Where("doctor_id = ? AND hospital_id = ? AND appointment_date = ? AND status IN (?)",
		doctorID, hospitalID, date, []string{"confirmed", "pending"}).Find(&bookedBookings).Error

	if err != nil {
		helper.SendJSON(w, http.StatusInternalServerError, "error", "Database error fetching booked slots", nil)
		return
	}

	// 4. Map Booked Slots (Nee Old Logic)
	bookedMap := make(map[string]bool)
	for _, b := range bookedBookings {
		cleanSlot := strings.ToUpper(strings.TrimSpace(b.TimeSlot))
		bookedMap[cleanSlot] = true
	}

	// 5. Filter Booked vs Available
	var availableSlots []string
	for _, s := range allSlots {
		if !bookedMap[strings.ToUpper(s)] {
			availableSlots = append(availableSlots, s)
		}
	}

	// 6. Return Response
	if len(availableSlots) == 0 {
		helper.SendJSON(w, http.StatusOK, "success", "No future slots available for today. Try another date.", map[string]interface{}{
			"available_slots": []string{},
			"total_count":     0,
		})
		return
	}

	helper.SendJSON(w, http.StatusOK, "success", "Available slots fetched successfully", map[string]interface{}{
		"date":            date,
		"doctor_id":       doctorID,
		"hospital_id":     hospitalID,
		"available_slots": availableSlots,
		"total_count":     len(availableSlots),
	})
}

// 3. GET MY BOOKINGS (With Filters)
func GetMyBookings(w http.ResponseWriter, r *http.Request) {
	mobile := r.URL.Query().Get("mobile")
	filter := r.URL.Query().Get("filter")

	if mobile == "" {
		helper.SendJSON(w, http.StatusBadRequest, "failed", "Mobile number required", nil)
		return
	}

	today := time.Now().Format("2006-01-02")
	query := database.DB.Where("account_holder_mobile = ?", mobile)

	switch strings.ToLower(filter) {
	case "today":
		query = query.Where("appointment_date = ?", today)
	case "upcoming":
		query = query.Where("appointment_date > ?", today)
	case "past":
		query = query.Where("appointment_date < ?", today)
	}

	var history []appointment.Booking
	query.Order("appointment_date desc").Find(&history)
	helper.SendJSON(w, http.StatusOK, "success", fmt.Sprintf("Found %d records", len(history)), history)
}
